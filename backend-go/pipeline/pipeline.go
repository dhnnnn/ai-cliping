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
		"-f", "mp4",
		"--merge-output-format", "mp4",
		"-o", videoPath,
		job.URL,
	)


	if err != nil {
		job.Status = models.StatusFailed
		job.Result = "download failed"
		return
	}

	job.Status = models.StatusDownloaded

	// 2. Extract audio
	job.Status = models.StatusAudioExtract
	audioPath := filepath.Join("storage", "audio", job.ID+".wav")

	err = utils.ExtractAudio(videoPath, audioPath)
	if err != nil {
		job.Status = models.StatusFailed
		return
	}
	job.Status = models.StatusAudioReady
	job.Result = fmt.Sprintf("video saved at %s", videoPath)
}
