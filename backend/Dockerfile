# ─────────────────────────────────────────────
# Stage 1: Build Go binary (Pakai Alpine)
# ─────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Alpine pakai apk
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

# Ubuntu pakai apt-get
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Python packages
RUN pip3 install --no-cache-dir \
    yt-dlp \
    faster-whisper

COPY --from=builder /app/ai-clipping-backend .
COPY scripts/ ./scripts/

RUN mkdir -p storage/videos storage/audio storage/transcripts storage/highlights storage/clips

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/health || exit 1

CMD ["./ai-clipping-backend"]
