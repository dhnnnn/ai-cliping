package utils

import (
	"fmt"
	"os"
	"strings"

	"ai-clipping-backend/models"
)

// SubtitleConfig holds subtitle styling options
type SubtitleConfig struct {
	Enabled      bool
	FontSize     int
	FontName     string
	FontColor    string
	OutlineColor string
	OutlineWidth int
	Position     string // "bottom", "top", "middle"
	MaxWidth     int    // Max characters per line
}

// DefaultSubtitleConfig returns default subtitle styling
func DefaultSubtitleConfig() SubtitleConfig {
	return SubtitleConfig{
		Enabled:      true,
		FontSize:     32, // Larger for word-by-word (was 24)
		FontName:     "Arial Bold",
		FontColor:    "white",
		OutlineColor: "black",
		OutlineWidth: 3, // Thicker outline for better visibility (was 2)
		Position:     "bottom",
		MaxWidth:     40,
	}
}

// CreateSubtitleFile generates SRT subtitle file for a clip with word-by-word timing
func CreateSubtitleFile(transcript []models.TranscriptSegment, start, end float64, outputPath string) error {
	// Filter transcript segments within clip timerange
	var clipSegments []models.TranscriptSegment
	for _, seg := range transcript {
		if seg.End >= start && seg.Start <= end {
			adjustedSeg := models.TranscriptSegment{
				Start: seg.Start - start,
				End:   seg.End - start,
				Text:  seg.Text,
			}

			if adjustedSeg.Start < 0 {
				adjustedSeg.Start = 0
			}
			if adjustedSeg.End > (end - start) {
				adjustedSeg.End = end - start
			}

			clipSegments = append(clipSegments, adjustedSeg)
		}
	}

	if len(clipSegments) == 0 {
		return fmt.Errorf("no transcript segments found for clip range %.1f-%.1f", start, end)
	}

	// Create SRT content with word-by-word timing
	var srt strings.Builder
	subtitleIndex := 1

	for _, seg := range clipSegments {
		// Split text into words
		words := strings.Fields(seg.Text)
		if len(words) == 0 {
			continue
		}

		// Calculate time per word
		segmentDuration := seg.End - seg.Start
		timePerWord := segmentDuration / float64(len(words))

		// Group words (show 1-3 words at a time for better readability)
		wordsPerSubtitle := 2
		if len(words) <= 3 {
			wordsPerSubtitle = len(words)
		}

		for i := 0; i < len(words); i += wordsPerSubtitle {
			// Get word group
			endIdx := i + wordsPerSubtitle
			if endIdx > len(words) {
				endIdx = len(words)
			}
			wordGroup := strings.Join(words[i:endIdx], " ")

			// Calculate timing for this group
			wordStart := seg.Start + (float64(i) * timePerWord)
			wordEnd := seg.Start + (float64(endIdx) * timePerWord)

			// Create SRT entry
			fmt.Fprintf(&srt, "%d\n", subtitleIndex)
			fmt.Fprintf(&srt, "%s --> %s\n", formatSRTTime(wordStart), formatSRTTime(wordEnd))
			fmt.Fprintf(&srt, "%s\n\n", wordGroup)

			subtitleIndex++
		}
	}

	// Write to file
	return os.WriteFile(outputPath, []byte(srt.String()), 0644)
}

// formatSRTTime converts seconds to SRT time format (HH:MM:SS,mmm)
func formatSRTTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, millis)
}

// CreateClipWithSubtitles creates a video clip with burned-in subtitles
func CreateClipWithSubtitles(
	inputVideo, outputPath string,
	start, end float64,
	transcript []models.TranscriptSegment,
	clipConfig ClipConfig,
	subConfig SubtitleConfig,
) error {
	// Calculate actual timestamps with padding
	actualStart := start - clipConfig.PaddingBefore
	if actualStart < 0 {
		actualStart = 0
	}
	actualEnd := end + clipConfig.PaddingAfter
	clipDuration := actualEnd - actualStart

	// Create temporary subtitle file
	subtitlePath := outputPath + ".srt"
	defer os.Remove(subtitlePath) // Clean up after

	if subConfig.Enabled {
		if err := CreateSubtitleFile(transcript, actualStart, actualEnd, subtitlePath); err != nil {
			// If subtitle creation fails, create clip without subtitles
			return CreateClip(inputVideo, outputPath, start, end, clipConfig)
		}
	}

	// Build video filter
	videoFilter := buildVideoFilter(clipConfig.AspectRatio)

	// Add subtitle filter
	if subConfig.Enabled {
		subtitleFilter := buildSubtitleFilter(subtitlePath, subConfig)
		if videoFilter != "" {
			videoFilter = videoFilter + "," + subtitleFilter
		} else {
			videoFilter = subtitleFilter
		}
	}

	// Build audio filter
	audioFilter := buildAudioFilter(clipConfig, clipDuration)

	// Build FFmpeg command
	args := []string{
		"-ss", fmt.Sprintf("%.3f", actualStart),
		"-to", fmt.Sprintf("%.3f", actualEnd),
		"-i", inputVideo,
	}

	if videoFilter != "" {
		args = append(args, "-vf", videoFilter)
	}

	if audioFilter != "" {
		args = append(args, "-af", audioFilter)
	}

	args = append(args,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-y",
		outputPath,
	)

	// Execute FFmpeg
	return RunCommand("./ffmpeg.exe", args...)
}

// buildSubtitleFilter creates FFmpeg subtitle filter string
func buildSubtitleFilter(subtitlePath string, config SubtitleConfig) string {
	// Convert Windows path to FFmpeg-compatible path
	escapedPath := strings.ReplaceAll(subtitlePath, "\\", "/")
	escapedPath = strings.ReplaceAll(escapedPath, ":", "\\\\:")

	// Build subtitles filter with styling
	// Alignment: 2=bottom center, 5=middle center, 8=top center
	alignment := 2 // Default bottom
	switch config.Position {
	case "top":
		alignment = 8
	case "middle":
		alignment = 5
	}

	filter := fmt.Sprintf("subtitles=%s:force_style='FontName=%s,FontSize=%d,PrimaryColour=&H%s,OutlineColour=&H%s,Outline=%d,Alignment=%d'",
		escapedPath,
		config.FontName,
		config.FontSize,
		colorToFFmpeg(config.FontColor),
		colorToFFmpeg(config.OutlineColor),
		config.OutlineWidth,
		alignment,
	)

	return filter
}

// colorToFFmpeg converts color name to FFmpeg format (AABBGGRR)
func colorToFFmpeg(color string) string {
	colors := map[string]string{
		"white":  "00FFFFFF",
		"black":  "00000000",
		"red":    "000000FF",
		"yellow": "0000FFFF",
		"green":  "0000FF00",
		"blue":   "00FF0000",
	}

	if hex, ok := colors[strings.ToLower(color)]; ok {
		return hex
	}
	return "00FFFFFF" // Default white
}
