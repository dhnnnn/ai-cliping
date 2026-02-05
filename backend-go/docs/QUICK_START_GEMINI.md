# Quick Setup: Gemini Highlight Detection

## 🚀 5-Minute Setup

### Step 1: Get Gemini API Key (2 minutes)

1. Open: https://makersuite.google.com/app/apikey
2. Click **"Get API Key"** → **"Create API key in new project"**
3. **Copy** the key (starts with `AIza...`)

### Step 2: Configure API Key (1 minute)

```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Create .env file
echo "GEMINI_API_KEY=AIzaSy...paste-your-key-here" > .env
```

### Step 3: Install Dependencies (2 minutes)

```powershell
# Install required packages
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option
go get github.com/joho/godotenv

# Tidy dependencies  
go mod tidy
```

### Step 4: Start Server

```powershell
go run main.go
```

Expected output:
```
Server running on :8080
```

---

## ✅ Test It!

### In Postman:

1. Submit video (request "1. Process Video (Simple)")
2. Check status until `"status": "completed"`
3. Response akan include:

```json
{
  "status": "completed",
  "highlights": [
    {
      "start": 15.5,
      "end": 32.8,
      "score": 88,
      "reason": "Clear explanation of key concept",
      "keywords": ["important", "tutorial", "step"]
    },
    {
      "start": 45.2,
      "end": 63.1,
      "score": 82,
      "reason": "Practical example demonstration",
      "keywords": ["example", "like this"]
    }
  ],
  "result": "Processing complete: ..., 3 highlights found"
}
```

---

## 📊 What Happened?

**Pipeline Flow:**
```
1. Download video ✓
2. Extract audio ✓
3. Transcribe with Whisper ✓
4. Analyze with Gemini AI 🆕✓
   - Reads transcript
   - Understands context
   - Finds 3-5 best moments
5. Ready for clipping!
```

**Gemini Does:**
- Reads entire transcript
- Understands video type (tutorial/gaming/podcast)
- Finds most interesting/engaging moments
- Scores each highlight (0-100)
- Returns timestamps + explanations

---

## 🎯 Next Step: Video Clipping

Setelah highlights detected, lanjut ke **Phase 3**:
- Use FFmpeg to cut video clips
- Convert to 9:16 vertical format
- Export multiple clips from highlights

Mau lanjut implement clipping sekarang? 🎬
