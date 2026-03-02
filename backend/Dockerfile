# ─────────────────────────────────────────────
# Stage 1: Build Go binary
# ─────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy dependency files dulu (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary untuk Linux (bukan .exe)
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-clipping-backend .

# ─────────────────────────────────────────────
# Stage 2: Runtime image
# ─────────────────────────────────────────────
# Pakai bookworm (Debian stable) bukan trixie (testing/unstable)
FROM python:3.11-slim-bookworm

WORKDIR /app

# Step 1: Install system packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Step 2: Install Python packages
# faster-whisper jauh lebih ringan dari openai-whisper (tidak butuh PyTorch penuh)
RUN pip3 install --no-cache-dir \
    yt-dlp \
    faster-whisper

# Copy Go binary dari stage builder
COPY --from=builder /app/ai-clipping-backend .

# Copy Python scripts (untuk Whisper transcription)
COPY scripts/ ./scripts/

# Buat folder storage yang dibutuhkan
RUN mkdir -p storage/videos storage/audio storage/transcripts storage/highlights storage/clips

# Railway inject PORT env secara otomatis
EXPOSE 8080

# Healthcheck agar Railway tahu container sehat
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

CMD ["./ai-clipping-backend"]
