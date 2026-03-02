# ─────────────────────────────────────────────
# Stage 1: Build Go binary
# ─────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ai-clipping-backend .

# ─────────────────────────────────────────────
# Stage 2: Runtime image
# Pakai Ubuntu 22.04 — repo lebih stabil di cloud
# ─────────────────────────────────────────────
FROM ubuntu:22.04

# Nonaktifkan interactive prompt saat install
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Step 1: Install system packages (ffmpeg, python3, curl)
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Step 2: Install Python packages
# faster-whisper = lebih ringan dari openai-whisper
RUN pip3 install --no-cache-dir \
    yt-dlp \
    faster-whisper

# Copy Go binary dari stage builder
COPY --from=builder /app/ai-clipping-backend .

# Copy Python scripts
COPY scripts/ ./scripts/

# Buat folder storage
RUN mkdir -p storage/videos storage/audio storage/transcripts storage/highlights storage/clips

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

CMD ["./ai-clipping-backend"]
