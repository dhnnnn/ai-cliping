package utils

import (
	"fmt"
	"log"
)

// ClipConfig holds configuration for video clipping
type ClipConfig struct {
	PaddingBefore   float64 // Seconds to add before highlight start
	PaddingAfter    float64 // Seconds to add after highlight end
	FadeInDuration  float64 // Audio fade in duration (seconds)
	FadeOutDuration float64 // Audio fade out duration (seconds)
	NormalizeAudio  bool    // Apply audio normalization
	AspectRatio     string  // Target aspect ratio (e.g., "9:16")
}

// DefaultClipConfig returns sensible defaults for smooth clipping
func DefaultClipConfig() ClipConfig {
	return ClipConfig{
		PaddingBefore:   0.5,  // Half second before
		PaddingAfter:    0.5,  // Half second after
		FadeInDuration:  0.3,  // 300ms fade in
		FadeOutDuration: 0.5,  // 500ms fade out
		NormalizeAudio:  true, // Normalize volume
		AspectRatio:     "9:16",
	}
}

// CreateClip creates a video clip with smart audio handling
func CreateClip(inputVideo, outputPath string, start, end float64, config ClipConfig) error {
	// Calculate actual timestamps with padding
	actualStart := start - config.PaddingBefore
	if actualStart < 0 {
		actualStart = 0
	}
	actualEnd := end + config.PaddingAfter
	clipDuration := actualEnd - actualStart

	// Build video filter for aspect ratio conversion
	videoFilter := buildVideoFilter(config.AspectRatio)

	// Build audio filter for smooth transitions
	audioFilter := buildAudioFilter(config, clipDuration)

	// Build FFmpeg command
	args := []string{
		"-ss", fmt.Sprintf("%.3f", actualStart), // Start time
		"-to", fmt.Sprintf("%.3f", actualEnd), // End time
		"-i", inputVideo, // Input file
	}

	// Add video filter
	if videoFilter != "" {
		args = append(args, "-vf", videoFilter)
	}

	// Add audio filter
	if audioFilter != "" {
		args = append(args, "-af", audioFilter)
	}

	// Add encoding options
	args = append(args,
		"-c:v", "libx264", // Video codec
		"-preset", "fast", // Encoding speed
		"-crf", "23", // Quality (18-28, lower = better)
		"-c:a", "aac", // Audio codec
		"-b:a", "128k", // Audio bitrate
		"-y",       // Overwrite output
		outputPath, // Output file
	)

	// Execute FFmpeg - use same path format as ffmpeg.go
	return RunCommand(FFmpegBinary(), args...)
}

// buildVideoFilter creates video filter string for aspect ratio conversion
func buildVideoFilter(aspectRatio string) string {
	if aspectRatio == "" || aspectRatio == "original" {
		return ""
	}

	switch aspectRatio {
	case "9:16":
		// Vertical format (1080x1920) for TikTok/Instagram
		// Scale to fit height, then crop to width
		return "scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920"

	case "16:9":
		// Horizontal format (1920x1080)
		return "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2"

	case "1:1":
		// Square format (1080x1080)
		return "scale=1080:1080:force_original_aspect_ratio=increase,crop=1080:1080"

	default:
		log.Printf("Warning: unknown aspect ratio %s, using original", aspectRatio)
		return ""
	}
}

// buildAudioFilter creates audio filter string for smooth transitions
func buildAudioFilter(config ClipConfig, clipDuration float64) string {
	filters := []string{}

	// Add fade in at the start
	if config.FadeInDuration > 0 {
		fadeIn := fmt.Sprintf("afade=t=in:st=0:d=%.2f", config.FadeInDuration)
		filters = append(filters, fadeIn)
	}

	// Add fade out at the end
	if config.FadeOutDuration > 0 {
		fadeOutStart := clipDuration - config.FadeOutDuration
		if fadeOutStart < 0 {
			fadeOutStart = 0
		}
		fadeOut := fmt.Sprintf("afade=t=out:st=%.2f:d=%.2f", fadeOutStart, config.FadeOutDuration)
		filters = append(filters, fadeOut)
	}

	// Add audio normalization for consistent volume
	if config.NormalizeAudio {
		filters = append(filters, "loudnorm")
	}

	if len(filters) == 0 {
		return ""
	}

	// Join filters with comma
	result := ""
	for i, filter := range filters {
		if i > 0 {
			result += ","
		}
		result += filter
	}
	return result
}
