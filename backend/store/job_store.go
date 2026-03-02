package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ai-clipping-backend/models"

	"github.com/redis/go-redis/v9"
)

const jobTTL = 24 * time.Hour // Job disimpan selama 24 jam

// JobStore menyimpan dan mengambil data Job dari Redis
type JobStore struct {
	client *redis.Client
}

// NewJobStore membuat instance JobStore baru dengan koneksi Redis
func NewJobStore(addr, password string, db int) *JobStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &JobStore{client: client}
}

// Ping mengecek koneksi ke Redis
func (s *JobStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// key menghasilkan Redis key untuk sebuah job
func (s *JobStore) key(id string) string {
	return fmt.Sprintf("job:%s", id)
}

// Save menyimpan job ke Redis
func (s *JobStore) Save(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}
	return s.client.Set(ctx, s.key(job.ID), data, jobTTL).Err()
}

// Get mengambil job dari Redis berdasarkan ID
func (s *JobStore) Get(ctx context.Context, id string) (*models.Job, bool) {
	data, err := s.client.Get(ctx, s.key(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false // Job tidak ditemukan
		}
		return nil, false
	}
	var job models.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return nil, false
	}
	return &job, true
}

// UpdateStatus mengupdate status dan result sebuah job secara atomik
func (s *JobStore) UpdateStatus(ctx context.Context, id string, status models.JobStatus, result string) error {
	job, ok := s.Get(ctx, id)
	if !ok {
		return fmt.Errorf("job %s not found", id)
	}
	job.Status = status
	if result != "" {
		job.Result = result
	}
	return s.Save(ctx, job)
}

// UpdateJob menyimpan ulang seluruh objek job (untuk update highlights, clips, dll)
func (s *JobStore) UpdateJob(ctx context.Context, job *models.Job) error {
	return s.Save(ctx, job)
}
