# Saved Files Structure

## 📁 Where Are the Results?

After processing a video, hasil analisis disimpan di berbagai folder:

### 1. Video & Audio
```
storage/videos/{job_id}.mp4    # Original downloaded video
storage/audio/{job_id}.wav     # Extracted audio (mono 16kHz)
```

### 2. Transcript (Whisper)
```
storage/transcripts/{job_id}.json
```

**Content:**
```json
[
  {
    "start": 0.0,
    "end": 3.5,
    "text": "Hello everyone, welcome to this tutorial"
  },
  {
    "start": 3.5,
    "end": 8.2,
    "text": "Today we're going to learn about..."
  }
]
```

### 3. Highlights (Gemini AI Analysis)
```
storage/highlights/{job_id}.json
```

**Content:**
```json
[
  {
    "start": 12.5,
    "end": 28.3,
    "score": 90,
    "reason": "Clear step-by-step instruction on key concept",
    "keywords": ["important", "step", "tutorial", "remember"]
  },
  {
    "start": 45.0,
    "end": 62.5,
    "score": 85,
    "reason": "Practical example demonstration",
    "keywords": ["example", "like this", "demonstration"]
  },
  {
    "start": 89.2,
    "end": 105.8,
    "score": 80,
    "reason": "Summary of key points",
    "keywords": ["conclusion", "summary", "important"]
  }
]
```

### 4. Clips (After Phase 3)
```
storage/clips/{job_id}_clip_1.mp4
storage/clips/{job_id}_clip_2.mp4
storage/clips/{job_id}_clip_3.mp4
```

---

## 🔍 How to Find Your Results

### Option 1: Check Folder Directly

```powershell
# List all transcripts
dir storage\transcripts\

# List all highlights
dir storage\highlights\

# View specific transcript
cat storage\transcripts\{job_id}.json

# View specific highlights
cat storage\highlights\{job_id}.json
```

### Option 2: From API Response

The `result` field shows all file paths:
```json
{
  "id": "abc-123-xyz",
  "status": "completed",
  "result": "Processing complete: video at storage\\videos\\abc-123-xyz.mp4, audio at storage\\audio\\abc-123-xyz.wav, transcript at storage\\transcripts\\abc-123-xyz.json, 3 highlights found"
}
```

### Option 3: Server Logs

```
Job abc-123-xyz: transcript file saved successfully
Job abc-123-xyz: highlights file saved successfully
```

---

## 📊 File Sizes (Approximate)

| File Type | Size | Example |
|-----------|------|---------|
| **Video** (5 min) | 10-50 MB | video.mp4 |
| **Audio** (5 min) | 5 MB | audio.wav |
| **Transcript** | 5-20 KB | transcript.json |
| **Highlights** | 1-5 KB | highlights.json |
| **Clip** (30s) | 3-8 MB | clip.mp4 |

---

## 🗂️ Complete Example

After processing job `abc-123-xyz`:

```
storage/
├── videos/
│   └── abc-123-xyz.mp4          (15 MB)
├── audio/
│   └── abc-123-xyz.wav          (4.5 MB)
├── transcripts/
│   └── abc-123-xyz.json         (8 KB)
├── highlights/
│   └── abc-123-xyz.json         (2 KB)
└── clips/
    ├── abc-123-xyz_clip_1.mp4   (5 MB)
    ├── abc-123-xyz_clip_2.mp4   (6 MB)
    └── abc-123-xyz_clip_3.mp4   (4 MB)
```

**Total:** ~45 MB per video (includes clips)

---

## 🔄 Workflow

```
User submits URL
    ↓
Download video → storage/videos/{id}.mp4
    ↓
Extract audio → storage/audio/{id}.wav
    ↓
Transcribe → storage/transcripts/{id}.json ✅
    ↓
Analyze → storage/highlights/{id}.json ✅
    ↓
Clip (Phase 3) → storage/clips/{id}_clip_*.mp4
```

---

## 🧹 Cleanup

Jika butuh hapus hasil lama:

```powershell
# Hapus semua videos
del storage\videos\*.mp4

# Hapus semua audio
del storage\audio\*.wav

# Hapus semua transcripts
del storage\transcripts\*.json

# Hapus semua highlights
del storage\highlights\*.json

# Hapus semua clips
del storage\clips\*.mp4
```

**Note:** `.gitkeep` files akan tetap ada untuk maintain folder structure.

---

## ✅ Verify Files Created

Setelah job completed:

```powershell
# Check if all files exist for a job
$jobId = "abc-123-xyz"

Test-Path "storage\videos\$jobId.mp4"
Test-Path "storage\audio\$jobId.wav"
Test-Path "storage\transcripts\$jobId.json"
Test-Path "storage\highlights\$jobId.json"
```

Should all return `True` ✅
