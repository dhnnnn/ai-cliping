package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
)

// JobResponse adalah response API tanpa field transcript (terlalu besar)
type JobResponse struct {
	ID          string                 `json:"id"`
	URL         string                 `json:"url"`
	Status      models.JobStatus       `json:"status"`
	Result      string                 `json:"result"`
	VideoType   models.VideoType       `json:"videoType"`
	Preferences models.ClipPreferences `json:"preferences"`
	Highlights  []models.Highlight     `json:"highlights,omitempty"`
	ClipPaths   []string               `json:"clipPaths,omitempty"`
}

func StatusHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/status/")
		if id == "" {
			http.Error(w, "job id is required", http.StatusBadRequest)
			return
		}

		// Ambil job dari Redis store
		job, ok := q.Get(r.Context(), id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		// Buat response tanpa transcript (transcript tersedia di file)
		response := JobResponse{
			ID:          job.ID,
			URL:         job.URL,
			Status:      job.Status,
			Result:      job.Result,
			VideoType:   job.VideoType,
			Preferences: job.Preferences,
			Highlights:  job.Highlights,
			ClipPaths:   job.ClipPaths,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
