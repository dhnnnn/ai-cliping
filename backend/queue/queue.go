package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-clipping-backend/models"
	"ai-clipping-backend/store"
	"ai-clipping-backend/tasks"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// JobQueue mengelola pengiriman job ke Asynq (Redis)
type JobQueue struct {
	client   *asynq.Client
	jobStore *store.JobStore
}

// NewJobQueue membuat instance JobQueue baru
func NewJobQueue(redisOpt asynq.RedisClientOpt, jobStore *store.JobStore) *JobQueue {
	client := asynq.NewClient(redisOpt)
	return &JobQueue{
		client:   client,
		jobStore: jobStore,
	}
}

// Add membuat job baru, menyimpannya ke Redis, dan mengirimnya ke antrian Asynq
func (q *JobQueue) Add(ctx context.Context, url string, videoType models.VideoType, prefs models.ClipPreferences) (*models.Job, error) {
	// Buat job baru
	job := &models.Job{
		ID:          uuid.NewString(),
		URL:         url,
		Status:      models.StatusQueued,
		VideoType:   videoType,
		Preferences: prefs,
	}

	// Simpan job ke Redis Job Store terlebih dahulu
	if err := q.jobStore.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to save job to store: %w", err)
	}

	// Buat payload task Asynq
	payload, err := tasks.NewVideoProcessPayload(job)
	if err != nil {
		return nil, fmt.Errorf("failed to create task payload: %w", err)
	}

	// Enqueue task ke Asynq dengan konfigurasi retry & timeout
	task := asynq.NewTask(tasks.TypeVideoProcess, payload)
	info, err := q.client.EnqueueContext(ctx, task,
		asynq.MaxRetry(3),             // Retry maksimal 3 kali jika gagal
		asynq.Timeout(30*time.Minute), // Timeout 30 menit per job
		asynq.Queue("default"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("[Queue] Job %s enqueued (asynq task ID: %s)", job.ID, info.ID)
	return job, nil
}

// Get mengambil status job dari Redis Job Store
func (q *JobQueue) Get(ctx context.Context, id string) (*models.Job, bool) {
	return q.jobStore.Get(ctx, id)
}

// Close menutup koneksi Asynq client
func (q *JobQueue) Close() error {
	return q.client.Close()
}
