package server

import (
	"net/http"

	"ai-clipping-backend/handlers"
	"ai-clipping-backend/queue"
)

func NewServer(q *queue.JobQueue) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint (digunakan Railway untuk cek status container)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/process", handlers.ProcessHandler(q))
	mux.HandleFunc("/api/status/", handlers.StatusHandler(q))

	return mux
}
