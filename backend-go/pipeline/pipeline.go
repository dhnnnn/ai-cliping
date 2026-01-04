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
	job.Status = models.StatusProcessing

	videoPath := filepath.Join("storage", "videos", job.ID+".%(ext)s")

	err := utils.RunCommand(
		"./yt-dlp.exe",
		"-f", "bestvideo+bestaudio/best",
		"--merge-output-format", "mp4",
		"-o", videoPath,
		job.URL,
	)


	if err != nil {
		job.Status = models.StatusFailed
		job.Result = "download failed"
		return
	}

	job.Status = models.StatusDone
	job.Result = fmt.Sprintf("video saved at %s", videoPath)
}
