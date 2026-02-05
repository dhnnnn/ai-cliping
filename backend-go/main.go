package main

import (
	"log"
	"net/http"
	"os"

	"ai-clipping-backend/pipeline"
	"ai-clipping-backend/queue"
	"ai-clipping-backend/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get Gemini API key from environment
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Println("WARNING: GEMINI_API_KEY not set, highlight detection will be skipped")
	}

	jobQueue := queue.NewJobQueue(10)

	go pipeline.StartWorker(jobQueue.Jobs, geminiAPIKey)

	handler := server.NewServer(jobQueue)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
