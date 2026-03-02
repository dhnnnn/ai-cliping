package handlers

import (
	"encoding/json"
	"net/http"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
)

type ProcessRequest struct {
	URL         string                  `json:"url"`
	VideoType   models.VideoType        `json:"videoType"`   // optional, default "general"
	Preferences *models.ClipPreferences `json:"preferences"` // optional
}

func ProcessHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProcessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		// Set default video type
		if req.VideoType == "" {
			req.VideoType = models.VideoTypeGeneral
		}

		// Set default preferences jika tidak diberikan
		var prefs models.ClipPreferences
		if req.Preferences != nil {
			prefs = *req.Preferences
		} else {
			prefs = models.ClipPreferences{
				MinDuration: 15,
				MaxDuration: 60,
				MaxClips:    5,
				AspectRatio: "9:16",
			}
		}

		// Tambahkan job ke antrian (enqueue ke Redis via Asynq)
		job, err := q.Add(r.Context(), req.URL, req.VideoType, prefs)
		if err != nil {
			http.Error(w, "failed to enqueue job: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		json.NewEncoder(w).Encode(job)
	}
}
