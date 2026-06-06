# ─────────────────────────────────────────────
# Stage 1: Build Go binary (Pakai Alpine)
# ─────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ai-clipping-backend .

# ─────────────────────────────────────────────
# Stage 2: Runtime image (Biar hemat storage & stabil)
# ─────────────────────────────────────────────
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Install system packages (ffmpeg + python untuk transcribe & yt-dlp)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ffmpeg \
    python3 \
    python3-pip \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Python packages (yt-dlp untuk download, faster-whisper untuk transcribe)
RUN pip3 install --no-cache-dir --upgrade pip && \
    pip3 install --no-cache-dir \
    yt-dlp \
    faster-whisper

COPY --from=builder /app/ai-clipping-backend .
COPY scripts/ ./scripts/

# Pastikan folder storage ada
RUN mkdir -p storage/videos storage/audio storage/transcripts storage/highlights storage/clips

# Port yang diekspos
EXPOSE 8080

# Jalankan binary
CMD ["./ai-clipping-backend"]
