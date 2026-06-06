package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai-clipping-backend/models"
	"ai-clipping-backend/queue"
	"ai-clipping-backend/youtube"
)

// YouTubeAuthHandler memulai alur OAuth2 dengan redirect ke halaman consent Google.
func YouTubeAuthHandler(up *youtube.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if up == nil {
			http.Error(w, "YouTube uploader belum dikonfigurasi (cek YOUTUBE_CLIENT_ID/SECRET)", http.StatusServiceUnavailable)
			return
		}
		url := up.AuthURL("state-token")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// YouTubeCallbackHandler menerima authorization code dari Google dan menyimpan token.
func YouTubeCallbackHandler(up *youtube.Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if up == nil {
			http.Error(w, "YouTube uploader belum dikonfigurasi", http.StatusServiceUnavailable)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "authorization code tidak ditemukan", http.StatusBadRequest)
			return
		}

		if err := up.ExchangeAndSave(r.Context(), code); err != nil {
			http.Error(w, "gagal menyimpan token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","message":"YouTube berhasil terhubung. Sekarang kamu bisa upload clip."}`))
	}
}

// UploadRequest adalah body untuk endpoint upload.
type UploadRequest struct {
	JobID         string `json:"jobId"`                   // wajib
	ClipIndex     *int   `json:"clipIndex,omitempty"`     // opsional: index 1-based clip tertentu; jika kosong, upload semua clip
	PrivacyStatus string `json:"privacyStatus,omitempty"` // "private" | "unlisted" | "public" (default)
	Description   string `json:"description,omitempty"`
}

// YouTubeUploadHandler memulai proses upload clip dari sebuah job ke YouTube.
// Upload berjalan di background; status tiap clip bisa dipantau via
// GET /api/youtube/status/<jobId>.
func YouTubeUploadHandler(up *youtube.Uploader, q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if up == nil {
			http.Error(w, "YouTube uploader belum dikonfigurasi (cek YOUTUBE_CLIENT_ID/SECRET)", http.StatusServiceUnavailable)
			return
		}
		if !up.IsAuthorized() {
			http.Error(w, "belum terhubung ke YouTube. Buka /api/youtube/auth dulu di browser.", http.StatusUnauthorized)
			return
		}

		var req UploadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.JobID == "" {
			http.Error(w, "jobId is required", http.StatusBadRequest)
			return
		}

		job, ok := q.Get(r.Context(), req.JobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		if len(job.ClipPaths) == 0 {
			http.Error(w, "job belum punya clip (mungkin belum selesai diproses)", http.StatusBadRequest)
			return
		}

		// Tentukan clip mana yang akan diupload.
		indices := make([]int, 0, len(job.ClipPaths))
		if req.ClipIndex != nil {
			idx := *req.ClipIndex - 1 // request 1-based → slice 0-based
			if idx < 0 || idx >= len(job.ClipPaths) {
				http.Error(w, "clipIndex di luar jangkauan", http.StatusBadRequest)
				return
			}
			indices = append(indices, idx)
		} else {
			for i := range job.ClipPaths {
				indices = append(indices, i)
			}
		}

		// Tandai clip yang akan diproses sebagai "uploading" (skip yang sudah sukses).
		queued := make([]int, 0, len(indices))
		for _, i := range indices {
			clipPath := job.ClipPaths[i]
			if existing := findUpload(job.Uploads, clipPath); existing != nil && existing.Status == "uploaded" {
				continue // sudah pernah sukses, jangan upload ulang
			}
			setUploadRecord(job, models.UploadRecord{
				ClipPath: clipPath,
				Status:   "uploading",
				Title:    titleForClip(job, i),
			})
			queued = append(queued, i)
		}

		if err := q.Update(r.Context(), job); err != nil {
			http.Error(w, "gagal menyimpan status awal upload: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Jalankan upload di background agar request tidak menggantung lama.
		go runUploads(up, q, req, job.ID, queued)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "uploading",
			"jobId":   req.JobID,
			"message": fmt.Sprintf("%d clip sedang diupload. Pantau via GET /api/youtube/status/%s", len(queued), req.JobID),
			"queued":  len(queued),
		})
	}
}

// runUploads memproses upload tiap clip secara berurutan dan meng-update
// status di Redis setelah tiap clip selesai (sukses maupun gagal).
func runUploads(up *youtube.Uploader, q *queue.JobQueue, req UploadRequest, jobID string, indices []int) {
	ctx := context.Background()

	for _, i := range indices {
		job, ok := q.Get(ctx, jobID)
		if !ok {
			fmt.Printf("[YouTube] job %s hilang saat upload, berhenti\n", jobID)
			return
		}
		if i >= len(job.ClipPaths) {
			continue
		}
		clipPath := job.ClipPaths[i]

		// Pastikan file masih ada.
		if _, err := os.Stat(clipPath); err != nil {
			markUpload(q, ctx, jobID, clipPath, func(rec *models.UploadRecord) {
				rec.Status = "failed"
				rec.Error = "file clip tidak ditemukan: " + clipPath
			})
			continue
		}

		res, err := up.Upload(ctx, youtube.UploadOptions{
			FilePath:      clipPath,
			Title:         titleForClip(job, i),
			Description:   req.Description,
			Tags:          tagsForClip(job, i),
			PrivacyStatus: req.PrivacyStatus,
		})
		if err != nil {
			markUpload(q, ctx, jobID, clipPath, func(rec *models.UploadRecord) {
				rec.Status = "failed"
				rec.Error = err.Error()
			})
			fmt.Printf("[YouTube] upload gagal job %s clip %s: %v\n", jobID, clipPath, err)
			continue
		}

		markUpload(q, ctx, jobID, clipPath, func(rec *models.UploadRecord) {
			rec.Status = "uploaded"
			rec.VideoID = res.VideoID
			rec.URL = res.URL
			rec.Title = res.Title
			rec.Error = ""
			rec.UploadedAt = time.Now().UTC().Format(time.RFC3339)
		})
		fmt.Printf("[YouTube] job %s clip %s -> %s\n", jobID, clipPath, res.URL)
	}
}

// UploadStatusResponse adalah ringkasan status upload sebuah job.
type UploadStatusResponse struct {
	JobID     string                `json:"jobId"`
	Total     int                   `json:"totalClips"`
	Uploaded  int                   `json:"uploaded"`
	Uploading int                   `json:"uploading"`
	Failed    int                   `json:"failed"`
	Pending   int                   `json:"pending"` // clip yang belum pernah diupload sama sekali
	Done      bool                  `json:"done"`    // true jika tidak ada lagi yang "uploading"
	Uploads   []models.UploadRecord `json:"uploads"`
}

// YouTubeStatusHandler mengembalikan status upload semua clip dari sebuah job.
func YouTubeStatusHandler(q *queue.JobQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/youtube/status/")
		if id == "" {
			http.Error(w, "job id is required", http.StatusBadRequest)
			return
		}

		job, ok := q.Get(r.Context(), id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		resp := UploadStatusResponse{
			JobID:   id,
			Total:   len(job.ClipPaths),
			Uploads: job.Uploads,
		}

		for _, rec := range job.Uploads {
			switch rec.Status {
			case "uploaded":
				resp.Uploaded++
			case "uploading":
				resp.Uploading++
			case "failed":
				resp.Failed++
			}
		}
		resp.Pending = resp.Total - len(job.Uploads)
		if resp.Pending < 0 {
			resp.Pending = 0
		}
		resp.Done = resp.Uploading == 0

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

// titleForClip mengembalikan judul untuk clip ke-i (dari AI, fallback ke nama file).
func titleForClip(job *models.Job, i int) string {
	if i < len(job.Highlights) {
		if t := strings.TrimSpace(job.Highlights[i].Title); t != "" {
			return t
		}
	}
	base := filepath.Base(job.ClipPaths[i])
	return base[:len(base)-len(filepath.Ext(base))]
}

// tagsForClip mengembalikan keywords sebagai tags untuk clip ke-i.
func tagsForClip(job *models.Job, i int) []string {
	if i < len(job.Highlights) {
		return job.Highlights[i].Keywords
	}
	return nil
}

// findUpload mengembalikan record upload untuk clipPath tertentu, atau nil jika belum ada.
func findUpload(uploads []models.UploadRecord, clipPath string) *models.UploadRecord {
	for i := range uploads {
		if uploads[i].ClipPath == clipPath {
			return &uploads[i]
		}
	}
	return nil
}

// setUploadRecord menambah atau menimpa record untuk clipPath di dalam job (in-memory).
func setUploadRecord(job *models.Job, rec models.UploadRecord) {
	for i := range job.Uploads {
		if job.Uploads[i].ClipPath == rec.ClipPath {
			job.Uploads[i] = rec
			return
		}
	}
	job.Uploads = append(job.Uploads, rec)
}

// markUpload mengambil job terbaru dari Redis, menerapkan perubahan pada record
// clipPath lewat fn, lalu menyimpannya kembali. Re-fetch mencegah race/overwrite.
func markUpload(q *queue.JobQueue, ctx context.Context, jobID, clipPath string, fn func(*models.UploadRecord)) {
	job, ok := q.Get(ctx, jobID)
	if !ok {
		return
	}

	rec := findUpload(job.Uploads, clipPath)
	if rec == nil {
		job.Uploads = append(job.Uploads, models.UploadRecord{ClipPath: clipPath})
		rec = &job.Uploads[len(job.Uploads)-1]
	}
	fn(rec)

	if err := q.Update(ctx, job); err != nil {
		fmt.Printf("[YouTube] WARNING: gagal update status upload job %s: %v\n", jobID, err)
	}
}
