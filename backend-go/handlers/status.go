package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
)

// JobResponse is the API response without the full transcript (too large)
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
		job, ok := q.Get(id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		// Create response without transcript (transcript available in file)
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

		json.NewEncoder(w).Encode(response)
	}
}
