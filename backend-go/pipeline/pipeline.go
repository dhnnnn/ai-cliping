package pipeline

import (
	"fmt"
	"log"
	"path/filepath"

	"ai-clipping-backend/models"
	"ai-clipping-backend/utils"
)

func StartWorker(jobChan <-chan *models.Job) {
	for job := range jobChan {
		process(job)
	}
}

func process(job *models.Job) {
	log.Println("Processing job:", job.ID)
	job.Status = models.StatusDownloading

	videoPath := filepath.Join("storage", "videos", job.ID+".mp4")

	err := utils.RunCommand(
		"./yt-dlp.exe",
		"-t", "mp4",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"--no-check-certificates",
		"--retries", "5",
		"--fragment-retries", "5",
		"-o", videoPath,
		job.URL,
	)

	if err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("download failed: %v", err)
		log.Printf("Job %s failed at download: %v", job.ID, err)
		return
	}

	job.Status = models.StatusDownloaded
	log.Printf("Job %s: video downloaded successfully", job.ID)

	// 2. Extract audio
	job.Status = models.StatusAudioExtract
	audioPath := filepath.Join("storage", "audio", job.ID+".wav")

	err = utils.ExtractAudio(videoPath, audioPath)
	if err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("audio extraction failed: %v", err)
		log.Printf("Job %s failed at audio extraction: %v", job.ID, err)
		return
	}
	job.Status = models.StatusAudioReady
	log.Printf("Job %s: audio extracted successfully", job.ID)

	// 3. Transcribe audio using Whisper
	job.Status = models.StatusTranscribing
	log.Printf("Job %s: starting transcription", job.ID)

	transcript, err := utils.TranscribeAudio(audioPath, "base")
	if err != nil {
		job.Status = models.StatusFailed
		job.Result = fmt.Sprintf("transcription failed: %v", err)
		log.Printf("Job %s failed at transcription: %v", job.ID, err)
		return
	}

	job.Transcript = transcript
	job.Status = models.StatusTranscribed
	log.Printf("Job %s: transcription completed, %d segments", job.ID, len(transcript))

	// Save transcript to JSON file
	transcriptPath := filepath.Join("storage", "transcripts", job.ID+".json")
	log.Printf("Job %s: saving transcript to %s", job.ID, transcriptPath)

	if err := utils.SaveTranscript(transcript, transcriptPath); err != nil {
		log.Printf("Job %s: WARNING - Failed to save transcript file: %v", job.ID, err)
		// Don't fail the job, transcript is already in job.Transcript field
	} else {
		log.Printf("Job %s: transcript file saved successfully", job.ID)
	}

	job.Result = fmt.Sprintf("video saved at %s, audio at %s, transcript at %s", videoPath, audioPath, transcriptPath)
}
