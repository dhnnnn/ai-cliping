package pipeline

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai-clipping-backend/analysis"
	"ai-clipping-backend/models"
	"ai-clipping-backend/utils"
)

func StartWorker(jobChan <-chan *models.Job, geminiAPIKey string) {
	for job := range jobChan {
		process(job, geminiAPIKey)
	}
}

func process(job *models.Job, geminiAPIKey string) {
	log.Printf("Processing job: %s", job.ID)
	job.Status = models.StatusDownloading

	videoPath := filepath.Join("storage", "videos", job.ID+".mp4")

	err := utils.RunCommand(
		"./yt-dlp.exe",
		"-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"--merge-output-format", "mp4", // Force MP4 output
		"-o", videoPath,
		job.URL,
	)

	if err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("Download failed: %v", err)
		log.Printf("Job %s failed: %v", job.ID, err)
		return
	}

	job.Status = models.StatusDownloaded

	// 2. Extract audio
	job.Status = models.StatusAudioExtract
	audioPath := filepath.Join("storage", "audio", job.ID+".wav")

	if err := utils.ExtractAudio(videoPath, audioPath); err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("Audio extraction failed: %v", err)
		log.Printf("Job %s failed: %v", job.ID, err)
		return
	}

	job.Status = models.StatusAudioReady

	// 3. Transcribe audio
	job.Status = models.StatusTranscribing
	log.Printf("Job %s: transcribing audio", job.ID)

	transcript, err := utils.TranscribeAudio(audioPath, "base")
	if err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("Transcription failed: %v", err)
		log.Printf("Job %s failed: %v", job.ID, err)
		return
	}

	job.Transcript = transcript
	job.Status = models.StatusTranscribed
	log.Printf("Job %s: transcription completed (%d segments)", job.ID, len(transcript))

	// Save transcript to JSON
	transcriptPath := filepath.Join("storage", "transcripts", job.ID+".json")
	if err := utils.SaveTranscript(transcript, transcriptPath); err != nil {
		log.Printf("Job %s: WARNING - Failed to save transcript: %v", job.ID, err)
	}

	// 4. Analyze highlights using Gemini
	if geminiAPIKey != "" {
		job.Status = models.StatusAnalyzing
		log.Printf("Job %s: analyzing highlights", job.ID)

		analyzer := analysis.NewGeminiAnalyzer(geminiAPIKey, job.VideoType, job.Preferences)
		highlights, err := analyzer.FindHighlights(transcript)

		if err != nil {
			log.Printf("Job %s: WARNING - Highlight analysis failed: %v", job.ID, err)
		} else {
			job.Highlights = highlights
			log.Printf("Job %s: found %d highlights", job.ID, len(highlights))

			// Save highlights to JSON
			highlightsPath := filepath.Join("storage", "highlights", job.ID+".json")
			if err := utils.SaveHighlights(highlights, highlightsPath); err != nil {
				log.Printf("Job %s: WARNING - Failed to save highlights: %v", job.ID, err)
			}
		}
	} else {
		log.Printf("Job %s: Skipping highlight analysis (no Gemini API key)", job.ID)
	}

	// 5. Create video clips from highlights
	if len(job.Highlights) > 0 {
		job.Status = models.StatusClipping
		log.Printf("Job %s: creating %d clips", job.ID, len(job.Highlights))

		clipPaths, err := createClips(job, videoPath)
		if err != nil {
			log.Printf("Job %s: WARNING - Clipping failed: %v", job.ID, err)
		} else {
			job.ClipPaths = clipPaths
			log.Printf("Job %s: created %d clips successfully", job.ID, len(clipPaths))
		}
	}

	job.Status = models.StatusCompleted
	job.Result = fmt.Sprintf("Processing complete: %d highlights, %d clips created",
		len(job.Highlights), len(job.ClipPaths))
	log.Printf("Job %s: completed successfully", job.ID)
}

// createClips generates video clips from highlights with optional subtitles
func createClips(job *models.Job, videoPath string) ([]string, error) {
	// Create clips in job-specific subfolder for better organization
	clipDir := filepath.Join("storage", "clips", job.ID)

	// Create directory if not exists
	if err := os.MkdirAll(clipDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create clip directory: %v", err)
	}

	// Use default smart clipping config
	clipConfig := utils.DefaultClipConfig()

	// Override aspect ratio from job preferences
	if job.Preferences.AspectRatio != "" {
		clipConfig.AspectRatio = job.Preferences.AspectRatio
	}

	// Enable subtitles by default
	subConfig := utils.DefaultSubtitleConfig()

	var clipPaths []string
	for i, highlight := range job.Highlights {
		outputPath := filepath.Join(clipDir, fmt.Sprintf("clip_%d.mp4", i+1))

		// Create clip with burned-in subtitles
		err := utils.CreateClipWithSubtitles(
			videoPath,
			outputPath,
			highlight.Start,
			highlight.End,
			job.Transcript,
			clipConfig,
			subConfig,
		)

		if err != nil {
			return nil, fmt.Errorf("clip %d failed: %v", i+1, err)
		}

		clipPaths = append(clipPaths, outputPath)
	}

	return clipPaths, nil
}
