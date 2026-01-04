package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-clipping-backend/queue"
)

func StatusHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/status/")
		job, ok := q.Get(id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(job)
	}
}
