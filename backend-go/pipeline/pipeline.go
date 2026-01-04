package pipeline

import (
	"log"
	"time"

	"ai-clipping-backend/models"
)

func StartWorker(jobChan <-chan *models.Job) {
	for job := range jobChan {
		process(job)
	}
}

func process(job *models.Job) {
	log.Println("Processing job:", job.ID)
	job.Status = models.StatusProcessing

	// STEP 1: Download video (yt-dlp)
	time.Sleep(2 * time.Second)

	// STEP 2: Extract audio (ffmpeg)
	time.Sleep(1 * time.Second)

	// STEP 3: Call AI worker (Python)
	time.Sleep(2 * time.Second)

	// STEP 4: Clip video
	time.Sleep(1 * time.Second)

	job.Status = models.StatusDone
	job.Result = "clips/generated"
}
