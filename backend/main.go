package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"ai-clipping-backend/pipeline"
	"ai-clipping-backend/queue"
	"ai-clipping-backend/server"
	"ai-clipping-backend/store"
	"ai-clipping-backend/tasks"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("[Main] No .env file found, using environment variables")
	}

	// ── Konfigurasi umum ──────────────────────────────────────────────────
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Println("[Main] WARNING: GEMINI_API_KEY not set, highlight detection will be skipped")
	}

	workerCount := getenvInt("WORKER_COUNT", 3)

	// Railway inject PORT tanpa titik dua (misal: "8080")
	// http.ListenAndServe butuh format ":8080"
	rawPort := getenv("PORT", "8080")
	httpPort := fmt.Sprintf(":%s", strings.TrimPrefix(rawPort, ":"))

	// ── Setup Redis ───────────────────────────────────────────────────────
	// Railway menyediakan REDIS_URL secara otomatis.
	// Untuk lokal, gunakan REDIS_ADDR + REDIS_PASSWORD di .env
	redisOpt, redisAddr, redisPassword, redisDB := buildRedisOpt()
	log.Printf("[Main] Config: Redis=%s | Workers=%d | Port=%s", redisAddr, workerCount, httpPort)

	// ── Setup Job Store (Redis) ───────────────────────────────────────────
	jobStore := store.NewJobStore(redisAddr, redisPassword, redisDB)
	if err := jobStore.Ping(context.Background()); err != nil {
		log.Fatalf("[Main] ❌ Cannot connect to Redis at %s: %v", redisAddr, err)
	}
	log.Printf("[Main] ✅ Connected to Redis at %s", redisAddr)

	// ── Setup Job Queue (Asynq Producer) ─────────────────────────────────
	jobQueue := queue.NewJobQueue(redisOpt, jobStore)
	defer jobQueue.Close()

	// ── Setup Asynq Worker Server (Consumer) ─────────────────────────────
	processor := pipeline.NewVideoProcessor(jobStore, geminiAPIKey)

	workerServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: workerCount,
		Queues: map[string]int{
			"default": 1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("[Worker] ❌ Task failed: type=%s, err=%v", task.Type(), err)
		}),
	})

	// Daftarkan handler untuk task TypeVideoProcess
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeVideoProcess, processor.ProcessTask)

	// Jalankan worker di goroutine terpisah
	go func() {
		log.Printf("[Worker] 🚀 Starting %d workers...", workerCount)
		if err := workerServer.Run(mux); err != nil {
			log.Fatalf("[Worker] ❌ Worker server failed: %v", err)
		}
	}()

	// ── Jalankan HTTP Server ──────────────────────────────────────────────
	handler := server.NewServer(jobQueue)
	log.Printf("[Main] 🌐 HTTP server running on %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, handler))
}

// buildRedisOpt membangun konfigurasi Redis dari environment variables.
// Prioritas: REDIS_URL (Railway) → REDIS_ADDR + REDIS_PASSWORD (lokal)
func buildRedisOpt() (asynq.RedisClientOpt, string, string, int) {
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		// Parse REDIS_URL format: redis://:password@host:port
		u, err := url.Parse(redisURL)
		if err != nil {
			log.Fatalf("[Main] ❌ Invalid REDIS_URL: %v", err)
		}
		password, _ := u.User.Password()
		log.Println("[Main] Using REDIS_URL (Railway mode)")
		opt := asynq.RedisClientOpt{
			Addr:     u.Host,
			Password: password,
			DB:       0,
		}
		return opt, u.Host, password, 0
	}

	// Fallback ke manual config (development lokal)
	addr := getenv("REDIS_ADDR", "localhost:6379")
	password := getenv("REDIS_PASSWORD", "")
	db := getenvInt("REDIS_DB", 0)
	opt := asynq.RedisClientOpt{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	return opt, addr, password, db
}

// getenv mengambil nilai env atau mengembalikan fallback jika kosong
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getenvInt mengambil nilai env sebagai integer
func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
