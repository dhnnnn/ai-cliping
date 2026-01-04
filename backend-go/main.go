package main

import (
	"log"
	"net/http"

	"ai-clipping-backend/pipeline"
	"ai-clipping-backend/queue"
	"ai-clipping-backend/server"
)

func main() {
	jobQueue := queue.NewJobQueue(10)

	go pipeline.StartWorker(jobQueue.Jobs)

	handler := server.NewServer(jobQueue)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
