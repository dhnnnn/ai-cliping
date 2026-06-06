package pipeline

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"ai-clipping-backend/analysis"
	"ai-clipping-backend/models"
	"ai-clipping-backend/store"
	"ai-clipping-backend/tasks"
	"ai-clipping-backend/utils"

	"github.com/hibiken/asynq"
)

// ytDLPBinary mengembalikan nama binary yt-dlp sesuai OS
func ytDLPBinary() string {
	if runtime.GOOS == "windows" {
		return "./yt-dlp.exe"
	}
	return "yt-dlp" // di Linux sudah ada di PATH (diinstall via pip di Docker)
}

// VideoProcessor adalah Asynq task handler untuk memproses video
type VideoProcessor struct {
	jobStore     *store.JobStore
	geminiAPIKey string
}

// NewVideoProcessor membuat instance VideoProcessor baru
func NewVideoProcessor(jobStore *store.JobStore, geminiAPIKey string) *VideoProcessor {
	return &VideoProcessor{
		jobStore:     jobStore,
		geminiAPIKey: geminiAPIKey,
	}
}

// ProcessTask adalah entry point yang dipanggil oleh Asynq worker
// Signature ini wajib sesuai dengan asynq.HandlerFunc
func (p *VideoProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// Parse payload dari task
	payload, err := tasks.ParseVideoProcessPayload(t.Payload())
	if err != nil {
		return fmt.Errorf("failed to parse task payload: %w", err)
	}

	// Ambil job dari Redis store
	job, ok := p.jobStore.Get(ctx, payload.JobID)
	if !ok {
		return fmt.Errorf("job %s not found in store", payload.JobID)
	}

	log.Printf("[Worker] Starting job: %s (type: %s)", job.ID, job.VideoType)
	return p.process(ctx, job)
}

// process menjalankan seluruh pipeline: download → extract → transcribe → analyze → clip
func (p *VideoProcessor) process(ctx context.Context, job *models.Job) error {
	// Helper untuk update status ke Redis
	updateStatus := func(status models.JobStatus, result string) {
		job.Status = status
		if result != "" {
			job.Result = result
		}
		if err := p.jobStore.UpdateJob(ctx, job); err != nil {
			log.Printf("[Worker] Warning: failed to update status for job %s: %v", job.ID, err)
		}
	}

	// ── Step 1: Download video ─────────────────────────────────────────────
	updateStatus(models.StatusDownloading, "")
	videoPath := filepath.Join("storage", "videos", job.ID+".mp4")

	err := func() error {
		args := []string{
			"-f", "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
			"--merge-output-format", "mp4",
			"-o", videoPath,
		}

		// Jika ada cookies.txt di root, gunakan untuk bypass bot detection
		if _, err := os.Stat("cookies.txt"); err == nil {
			args = append(args, "--cookies", "cookies.txt")
			log.Println("[Worker] 🍪 Using cookies.txt for download")
		}

		args = append(args, job.URL)
		return utils.RunCommand(ytDLPBinary(), args...)
	}()
	if err != nil {
		updateStatus(models.StatusFailed, fmt.Sprintf("Download failed: %v", err))
		log.Printf("[Worker] Job %s failed at download: %v", job.ID, err)
		return fmt.Errorf("download failed: %w", err)
	}
	updateStatus(models.StatusDownloaded, "")

	// ── Step 2: Extract audio ──────────────────────────────────────────────
	updateStatus(models.StatusAudioExtract, "")
	audioPath := filepath.Join("storage", "audio", job.ID+".wav")

	if err := utils.ExtractAudio(videoPath, audioPath); err != nil {
		updateStatus(models.StatusFailed, fmt.Sprintf("Audio extraction failed: %v", err))
		log.Printf("[Worker] Job %s failed at audio extraction: %v", job.ID, err)
		return fmt.Errorf("audio extraction failed: %w", err)
	}
	updateStatus(models.StatusAudioReady, "")

	// ── Step 3: Transcribe audio ───────────────────────────────────────────
	updateStatus(models.StatusTranscribing, "")
	log.Printf("[Worker] Job %s: transcribing audio...", job.ID)

	transcript, err := utils.TranscribeAudio(audioPath, "base")
	if err != nil {
		updateStatus(models.StatusFailed, fmt.Sprintf("Transcription failed: %v", err))
		log.Printf("[Worker] Job %s failed at transcription: %v", job.ID, err)
		return fmt.Errorf("transcription failed: %w", err)
	}

	job.Transcript = transcript
	updateStatus(models.StatusTranscribed, "")
	log.Printf("[Worker] Job %s: transcribed %d segments", job.ID, len(transcript))

	// Simpan transcript ke file JSON
	transcriptPath := filepath.Join("storage", "transcripts", job.ID+".json")
	if err := utils.SaveTranscript(transcript, transcriptPath); err != nil {
		log.Printf("[Worker] Job %s: WARNING - failed to save transcript: %v", job.ID, err)
	}

	// ── Step 4: Analyze highlights dengan Gemini ──────────────────────────
	if p.geminiAPIKey != "" {
		updateStatus(models.StatusAnalyzing, "")
		log.Printf("[Worker] Job %s: analyzing highlights...", job.ID)

		analyzer := analysis.NewGeminiAnalyzer(p.geminiAPIKey, job.VideoType, job.Preferences)
		highlights, err := analyzer.FindHighlights(transcript)
		if err != nil {
			log.Printf("[Worker] Job %s: WARNING - highlight analysis failed: %v", job.ID, err)
		} else {
			job.Highlights = highlights
			log.Printf("[Worker] Job %s: found %d highlights", job.ID, len(highlights))

			// Simpan highlights ke file JSON
			highlightsPath := filepath.Join("storage", "highlights", job.ID+".json")
			if err := utils.SaveHighlights(highlights, highlightsPath); err != nil {
				log.Printf("[Worker] Job %s: WARNING - failed to save highlights: %v", job.ID, err)
			}
		}
	} else {
		log.Printf("[Worker] Job %s: skipping highlight analysis (no Gemini API key)", job.ID)
	}

	// ── Step 5: Create video clips ────────────────────────────────────────
	if len(job.Highlights) > 0 {
		updateStatus(models.StatusClipping, "")
		log.Printf("[Worker] Job %s: creating %d clips...", job.ID, len(job.Highlights))

		clipPaths, err := createClips(job, videoPath)
		if err != nil {
			log.Printf("[Worker] Job %s: WARNING - clipping failed: %v", job.ID, err)
		} else {
			job.ClipPaths = clipPaths
			log.Printf("[Worker] Job %s: created %d clips", job.ID, len(clipPaths))
		}
	}

	// ── Selesai ───────────────────────────────────────────────────────────
	finalResult := fmt.Sprintf("Processing complete: %d highlights, %d clips created",
		len(job.Highlights), len(job.ClipPaths))
	updateStatus(models.StatusCompleted, finalResult)
	log.Printf("[Worker] Job %s: DONE - %s", job.ID, finalResult)

	return nil
}

// createClips menghasilkan video clip dari setiap highlight
func createClips(job *models.Job, videoPath string) ([]string, error) {
	clipDir := filepath.Join("storage", "clips", job.ID)

	if err := os.MkdirAll(clipDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create clip directory: %v", err)
	}

	clipConfig := utils.DefaultClipConfig()
	if job.Preferences.AspectRatio != "" {
		clipConfig.AspectRatio = job.Preferences.AspectRatio
	}

	var clipPaths []string
	usedNames := make(map[string]bool)
	for i, highlight := range job.Highlights {
		// Pakai judul dari AI sebagai nama file (sudah di-sanitize agar aman)
		baseName := utils.SanitizeFilename(highlight.Title)

		// Prefix nomor urut biar clip terurut & nama tetap unik
		fileName := fmt.Sprintf("%02d - %s", i+1, baseName)

		// Jaga-jaga kalau ada judul yang sama persis
		candidate := fileName
		for n := 2; usedNames[candidate]; n++ {
			candidate = fmt.Sprintf("%s (%d)", fileName, n)
		}
		usedNames[candidate] = true

		outputPath := filepath.Join(clipDir, candidate+".mp4")

		err := utils.CreateClip(videoPath, outputPath, highlight.Start, highlight.End, clipConfig)
		if err != nil {
			return nil, fmt.Errorf("clip %d failed: %v", i+1, err)
		}

		clipPaths = append(clipPaths, outputPath)
	}

	return clipPaths, nil
}
