package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
)

type ProcessRequest struct {
	URL         string                  `json:"url"`
	VideoType   models.VideoType        `json:"videoType"`   // optional, default to "general"
	Preferences *models.ClipPreferences `json:"preferences"` // optional
}

func ProcessHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProcessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		// Set defaults
		if req.VideoType == "" {
			req.VideoType = models.VideoTypeGeneral
		}

		var prefs models.ClipPreferences
		if req.Preferences != nil {
			prefs = *req.Preferences
		} else {
			// Default preferences for 9:16 vertical clips
			prefs = models.ClipPreferences{
				MinDuration: 15,
				MaxDuration: 60,
				MaxClips:    5,
				AspectRatio: "9:16",
			}
		}

		job := &models.Job{
			ID:          uuid.NewString(),
			URL:         req.URL,
			Status:      models.StatusQueued,
			VideoType:   req.VideoType,
			Preferences: prefs,
		}

		q.Add(job)

		json.NewEncoder(w).Encode(job)
	}
}
