# Testing Phase 1: Transcription

## 📝 Setup Postman

1. **Import Collection:**
   - Buka Postman
   - Click "Import" → "Upload Files"
   - Pilih file `AI_Clipping_API.postman_collection.json`
   - Collection akan muncul di sidebar

2. **Set Variable:**
   - Click pada collection "AI Clipping Backend API"
   - Tab "Variables"
   - Set `job_id` value setelah dapat response dari request pertama

## 🧪 Testing Flow

### Test 1: Submit Video untuk Transcription

**Request:** "1. Process Video (Simple)"
```json
POST http://localhost:8080/api/process
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}
```

**Expected Response:**
```json
{
  "id": "abc-123-xyz",
  "url": "...",
  "status": "queued",
  "videoType": "general",
  "preferences": {
    "minDuration": 15,
    "maxDuration": 60,
    "maxClips": 5,
    "aspectRatio": "9:16"
  }
}
```

**Action:** Copy `id` value untuk check status

### Test 2: Check Status (Polling)

**Request:** "5. Check Job Status"
- Ganti `{{job_id}}` di URL dengan ID dari Test 1
- Atau set variable `job_id` di collection

**Expected Status Flow:**
```
queued (awal)
  ↓
downloading (yt-dlp running)
  ↓
downloaded (video saved)
  ↓
audio_extracting (ffmpeg running)
  ↓
audio_ready (audio extracted)
  ↓
transcribing (Whisper running) ← BARU!
  ↓
transcribed (transcript ready) ← BARU!
  ↓
(akan lanjut ke analyzing nanti)
```

**Final Response Should Include:**
```json
{
  "id": "abc-123-xyz",
  "status": "transcribed",
  "result": "video saved at..., audio at..., transcript at...",
  "transcript": [
    {
      "start": 0.0,
      "end": 3.5,
      "text": "Hello everyone..."
    },
    ...
  ]
}
```

### Test 3: Test dengan Video Types Berbeda

Coba request:
- "2. Process Video (Tutorial)"
- "3. Process Video (Gaming)" 
- "4. Process Video (Podcast)"

Setiap video type akan punya optimasi berbeda nanti di Phase 2 (highlight detection).

## 📊 What to Check

### ✅ Success Indicators

1. **Server Logs:**
   ```
   Processing job: abc-123-xyz
   Job abc-123-xyz: video downloaded successfully
   Job abc-123-xyz: audio extracted successfully
   Job abc-123-xyz: starting transcription
   Job abc-123-xyz: transcription completed, 25 segments
   ```

2. **Files Created:**
   - `storage/videos/{job_id}.mp4` ✓
   - `storage/audio/{job_id}.wav` ✓
   - `storage/transcripts/{job_id}.json` ✓ NEW!

3. **Transcript JSON Content:**
   ```json
   [
     {
       "start": 0.0,
       "end": 3.5,
       "text": "Actual transcribed text from video"
     }
   ]
   ```

### ❌ Common Issues

**Issue: "whisper execution failed"**
- Cek Python installed: `python --version`
- Cek Whisper installed: `pip show openai-whisper`
- Cek logs untuk detail error

**Issue: Transcription lambat**
- Normal untuk model `base` (~1-2 menit per video)
- Gunakan video pendek untuk testing (<5 menit)
- Model akan di-download saat first run (~1GB)

**Issue: "No module named 'whisper'"**
```bash
pip install --upgrade openai-whisper
```

## 🎯 Next After Testing

Setelah transcription berhasil, lanjut ke:
- **Phase 2:** Implement highlight detection (keyword, volume analysis)
- **Phase 3:** Implement video clipping dengan FFmpeg
- **Phase 4:** Add API endpoints untuk transcript & clips

## 💡 Tips

1. **Gunakan video pendek** untuk testing awal (<2-3 menit)
2. **Monitor server logs** untuk lihat progress real-time
3. **Check transcript JSON** untuk quality hasil Whisper
4. Set environment variable jika Python tidak ditemukan:
   ```powershell
   $env:PATH += ";C:\Users\YourUser\AppData\Local\Programs\Python\Python312"
   ```
