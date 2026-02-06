package server

import (
	"net/http"

	"ai-clipping-backend/handlers"
	"ai-clipping-backend/queue"
)

func NewServer(q *queue.JobQueue) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/process", handlers.ProcessHandler(q))
	mux.HandleFunc("/api/status/", handlers.StatusHandler(q))

	return mux
}
