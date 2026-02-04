# AI Clipping Backend

Backend service untuk mengunduh video dari URL dan mengekstrak audio menggunakan Go, yt-dlp, dan FFmpeg.

_Backend service for downloading videos from URLs and extracting audio using Go, yt-dlp, and FFmpeg._

## 📋 Deskripsi / Description

Aplikasi backend ini menyediakan API untuk:
- Mengunduh video dari berbagai platform (YouTube, dll) menggunakan yt-dlp
- Mengekstrak audio dari video yang diunduh dalam format WAV
- Mengelola antrian job processing secara asynchronous
- Melacak status setiap job (queued, downloading, downloaded, audio_extracting, audio_ready, failed)

_This backend application provides APIs for:_
- _Downloading videos from various platforms (YouTube, etc.) using yt-dlp_
- _Extracting audio from downloaded videos in WAV format_
- _Managing job processing queue asynchronously_
- _Tracking status of each job (queued, downloading, downloaded, audio_extracting, audio_ready, failed)_

## 🛠️ Teknologi / Technology Stack

- **Go** 1.25.5
- **yt-dlp** - Video downloader
- **FFmpeg** - Audio extraction
- **UUID** - Job ID generation

## 📦 Prerequisites / Prasyarat

Pastikan sudah terinstall:
- Go 1.25.5 atau lebih baru
- Git (opsional, untuk clone repository)

_Make sure you have installed:_
- _Go 1.25.5 or newer_
- _Git (optional, for cloning repository)_

## 🚀 Instalasi / Installation

### 1. Clone atau Download Project

```bash
# Jika menggunakan Git
git clone <repository-url>
cd backend-go

# Atau download dan extract zip file
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Verifikasi ffmpeg.exe dan yt-dlp.exe

Pastikan file `ffmpeg.exe` dan `yt-dlp.exe` sudah ada di root directory project. Jika belum ada:

_Make sure `ffmpeg.exe` and `yt-dlp.exe` files exist in the project root directory. If not:_

- **ffmpeg.exe**: Download dari [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html)
- **yt-dlp.exe**: Download dari [https://github.com/yt-dlp/yt-dlp/releases](https://github.com/yt-dlp/yt-dlp/releases)

### 4. Buat Folder Storage

```bash
# Buat folder untuk menyimpan hasil download
mkdir -p storage\videos
mkdir -p storage\audio
```

## ▶️ Cara Menjalankan / How to Run

### Development Mode

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

_Server will run at `http://localhost:8080`_

### Build dan Run

```bash
# Build executable
go build -o ai-clipping-backend.exe

# Jalankan executable
.\ai-clipping-backend.exe
```

## 📡 API Endpoints

### 1. Process Video (Submit Job)

**Endpoint:** `POST /api/process`

**Request Body:**
```json
{
  "url": "https://www.youtube.com/watch?v=example"
}
```

**Response:**
```json
{
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "URL": "https://www.youtube.com/watch?v=example",
  "Status": "queued",
  "Result": ""
}
```

**Contoh menggunakan curl:**
```bash
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://www.youtube.com/watch?v=dQw4w9WgXcQ\"}"
```

### 2. Check Status (Get Job Status)

**Endpoint:** `GET /api/status/{job_id}`

**Response:**
```json
{
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "URL": "https://www.youtube.com/watch?v=example",
  "Status": "audio_ready",
  "Result": "video saved at storage\\videos\\550e8400-e29b-41d4-a716-446655440000.mp4, audio at storage\\audio\\550e8400-e29b-41d4-a716-446655440000.wav"
}
```

**Contoh menggunakan curl:**
```bash
curl http://localhost:8080/api/status/550e8400-e29b-41d4-a716-446655440000
```

## 📊 Job Status Flow

1. **queued** - Job masuk ke antrian
2. **downloading** - Sedang mengunduh video
3. **downloaded** - Video berhasil diunduh
4. **audio_extracting** - Sedang mengekstrak audio
5. **audio_ready** - Audio berhasil diekstrak ✅
6. **failed** - Proses gagal ❌

## 📁 Struktur Project / Project Structure

```
backend-go/
├── handlers/          # HTTP request handlers
│   ├── process.go    # Handler untuk submit job
│   └── status.go     # Handler untuk cek status
├── models/           # Data models
│   └── job.go        # Job model dan status constants
├── pipeline/         # Processing pipeline
│   └── pipeline.go   # Worker untuk proses video
├── queue/            # Job queue management
│   └── queue.go      # Queue implementation
├── server/           # HTTP server setup
│   └── server.go     # Route configuration
├── storage/          # Output storage
│   ├── videos/       # Downloaded videos
│   └── audio/        # Extracted audio files
├── utils/            # Utility functions
│   ├── command.go    # Command execution helper
│   └── ffmpeg.go     # FFmpeg audio extraction
├── ffmpeg.exe        # FFmpeg binary
├── yt-dlp.exe        # yt-dlp binary
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
└── main.go           # Application entry point
```

## ⚙️ Konfigurasi / Configuration

### Port Server
Default port adalah `8080`. Untuk mengubahnya, edit file `main.go`:

```go
log.Fatal(http.ListenAndServe(":8080", handler)) // Ubah 8080 ke port lain
```

### Queue Size
Default queue size adalah `10` concurrent jobs. Untuk mengubahnya, edit file `main.go`:

```go
jobQueue := queue.NewJobQueue(10) // Ubah 10 ke jumlah yang diinginkan
```

## 🔧 Troubleshooting

### Error: "ffmpeg.exe not found" atau "yt-dlp.exe not found"
**Solusi:** Pastikan file `ffmpeg.exe` dan `yt-dlp.exe` ada di root directory project.

### Error: "storage/videos: The system cannot find the path specified"
**Solusi:** Buat folder storage secara manual:
```bash
mkdir storage\videos
mkdir storage\audio
```

### Video download gagal / HTTP Error 403: Forbidden
**Masalah:** YouTube memblokir download request (HTTP 403 Forbidden)

**Solusi:** 
1. **Update yt-dlp ke versi terbaru** (sangat penting!):
   ```bash
   # Download versi terbaru dari https://github.com/yt-dlp/yt-dlp/releases
   # Replace yt-dlp.exe dengan versi baru
   ```

2. **Gunakan cookies dari browser:**
   - Install browser extension untuk export cookies (seperti "Get cookies.txt LOCALLY")
   - Export cookies dari YouTube
   - Simpan sebagai `cookies.txt` di root project
   - Tambahkan flag `--cookies cookies.txt` di `pipeline.go`

3. **Cek URL video:**
   - Pastikan URL video valid dan dapat diakses
   - Pastikan video tidak private atau restricted
   - Coba buka URL di browser terlebih dahulu

4. **Gunakan proxy jika perlu:**
   - Tambahkan flag `--proxy socks5://127.0.0.1:1080` jika menggunakan proxy

5. **Install JavaScript runtime (opsional, untuk performa lebih baik):**
   ```bash
   # Install deno
   # Windows: https://deno.land/manual/getting_started/installation
   # Atau gunakan scoop: scoop install deno
   ```

### Error: "No supported JavaScript runtime"
**Solusi:** Install deno untuk ekstraksi YouTube yang lebih baik:
- Download dari: https://deno.land/
- Atau gunakan package manager: `scoop install deno` atau `choco install deno`

### Audio extraction gagal
**Solusi:**
- Pastikan ffmpeg.exe berfungsi dengan baik
- Periksa apakah video sudah berhasil didownload

## 📝 Contoh Penggunaan Lengkap / Complete Usage Example

### 1. Start Server
```bash
go run main.go
```

### 2. Submit Job untuk Download Video
```bash
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://www.youtube.com/watch?v=dQw4w9WgXcQ\"}"
```

Response:
```json
{
  "ID": "abc-123-def",
  "URL": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "Status": "queued",
  "Result": ""
}
```

### 3. Cek Status Job (gunakan ID dari response di atas)
```bash
curl http://localhost:8080/api/status/abc-123-def
```

Response:
```json
{
  "ID": "abc-123-def",
  "URL": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "Status": "audio_ready",
  "Result": "video saved at storage\\videos\\abc-123-def.mp4, audio at storage\\audio\\abc-123-def.wav"
}
```

### 4. Ambil File Hasil
- **Video:** `storage/videos/abc-123-def.mp4`
- **Audio:** `storage/audio/abc-123-def.wav` (mono, 16kHz sample rate)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

[Specify your license here]

## 👤 Author

[Your name/organization]

---

**Happy Coding! 🚀**
