package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
)

type ProcessRequest struct {
	URL string `json:"url"`
}

func ProcessHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProcessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		job := &models.Job{
			ID:     uuid.NewString(),
			URL:    req.URL,
			Status: models.StatusQueued,
		}

		q.Add(job)

		json.NewEncoder(w).Encode(job)
	}
}
