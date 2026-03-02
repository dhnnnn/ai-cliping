# ─────────────────────────────────────────────
# Stage 1: Build Go binary
# ─────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (dibutuhkan oleh beberapa go module)
RUN apk add --no-cache git

# Copy dependency files dulu (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary untuk Linux (bukan .exe)
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-clipping-backend .

# ─────────────────────────────────────────────
# Stage 2: Runtime image (lebih kecil)
# ─────────────────────────────────────────────
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies:
# - ffmpeg: untuk extract audio & create clips
# - yt-dlp: untuk download video (via pip, versi Linux)
# - curl: untuk healthcheck
RUN apt-get update && apt-get install -y \
    ffmpeg \
    curl \
    && pip3 install --no-cache-dir yt-dlp openai-whisper \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy Go binary dari stage builder
COPY --from=builder /app/ai-clipping-backend .

# Copy Python scripts (untuk Whisper transcription)
COPY scripts/ ./scripts/

# Buat folder storage yang dibutuhkan
RUN mkdir -p storage/videos storage/audio storage/transcripts storage/highlights storage/clips

# Port yang diekspos (Railway akan inject PORT env otomatis)
EXPOSE 8080

# Healthcheck agar Railway tahu container sehat
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

# Jalankan binary
CMD ["./ai-clipping-backend"]
