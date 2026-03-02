package tasks

import (
	"encoding/json"

	"ai-clipping-backend/models"
)

// TypeVideoProcess adalah nama task yang dikenali oleh Asynq worker
const TypeVideoProcess = "video:process"

// VideoProcessPayload adalah data yang disimpan di Redis untuk setiap job
type VideoProcessPayload struct {
	JobID       string                 `json:"job_id"`
	URL         string                 `json:"url"`
	VideoType   models.VideoType       `json:"video_type"`
	Preferences models.ClipPreferences `json:"preferences"`
}

// NewVideoProcessPayload membuat payload JSON dari sebuah Job
func NewVideoProcessPayload(job *models.Job) ([]byte, error) {
	payload := VideoProcessPayload{
		JobID:       job.ID,
		URL:         job.URL,
		VideoType:   job.VideoType,
		Preferences: job.Preferences,
	}
	return json.Marshal(payload)
}

// ParseVideoProcessPayload mem-parse payload JSON menjadi struct
func ParseVideoProcessPayload(data []byte) (*VideoProcessPayload, error) {
	var p VideoProcessPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
