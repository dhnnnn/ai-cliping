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
	job.Result = fmt.Sprintf("video saved at %s, audio at %s", videoPath, audioPath)
}
