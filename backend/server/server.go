package server

import (
	"net/http"

	"ai-clipping-backend/handlers"
	"ai-clipping-backend/queue"
	"ai-clipping-backend/youtube"
)

func NewServer(q *queue.JobQueue, yt *youtube.Uploader) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint (digunakan Railway untuk cek status container)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/process", handlers.ProcessHandler(q))
	mux.HandleFunc("/api/status/", handlers.StatusHandler(q))

	// YouTube Data API v3 (OAuth2 + upload)
	mux.HandleFunc("/api/youtube/auth", handlers.YouTubeAuthHandler(yt))
	mux.HandleFunc("/api/youtube/callback", handlers.YouTubeCallbackHandler(yt))
	mux.HandleFunc("/api/youtube/upload", handlers.YouTubeUploadHandler(yt, q))
	mux.HandleFunc("/api/youtube/status/", handlers.YouTubeStatusHandler(q))

	// Serve file hasil (clips, dll) dari folder storage.
	// Contoh akses: GET /storage/clips/<job-id>/01 - Judul.mp4
	fs := http.FileServer(http.Dir("storage"))
	mux.Handle("/storage/", http.StripPrefix("/storage/", fs))

	return mux
}
