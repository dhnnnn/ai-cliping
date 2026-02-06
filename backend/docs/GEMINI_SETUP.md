# Gemini API Setup Guide

## Step 1: Get Your Free API Key

1. Go to: **https://makersuite.google.com/app/apikey**
2. Sign in with your Google account
3. Click **"Get API Key"** or **"Create API Key"**
4. Click **"Create API key in new project"** (or select existing project)
5. **Copy** the API key (starts with `AIza...`)

**Important:** Keep this key secret! Don't commit to git.

---

## Step 2: Configure API Key

### Option A: Using .env File (Recommended)

1. Create `.env` file di project root:
   ```powershell
   cd d:\.dev\golang\ai-clipping\backend-go
   copy .env.example .env
   ```

2. Edit `.env` file:
   ```
   GEMINI_API_KEY=AIzaSy...your-actual-key-here
   ```

3. Save file

### Option B: Set Environment Variable

**PowerShell (temporary - current session only):**
```powershell
$env:GEMINI_API_KEY="AIzaSy...your-actual-key-here"
```

**PowerShell (permanent):**
```powershell
[System.Environment]::SetEnvironmentVariable('GEMINI_API_KEY', 'AIzaSy...your-actual-key-here', 'User')
```

**CMD:**
```cmd
set GEMINI_API_KEY=AIzaSy...your-actual-key-here
```

---

## Step 3: Install Dependencies

```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Install Gemini SDK
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option

# Install dotenv for .env file support
go get github.com/joho/godotenv

# Update all dependencies
go mod tidy
```

---

## Step 4: Verify Setup

```powershell
# Start server
go run main.go
```

Check logs for:
```
✅ Server running on :8080
```

If no API key error, you're good to go!

If you see:
```
⚠️  WARNING: GEMINI_API_KEY not set, highlight detection will be skipped
```

Then API key not loaded correctly. Check .env file or environment variable.

---

## Step 5: Test with Postman

1. Submit a video in Postman (any of the existing requests)
2. Wait for transcription to complete
3. Check status - should see new field:
   ```json
   {
     "status": "completed",
     "highlights": [
       {
         "start": 10.5,
         "end": 25.3,
         "score": 85,
         "reason": "Key concept explained clearly",
         "keywords": ["important", "tutorial", "step"]
       }
     ]
   }
   ```

---

## 🎯 Expected Behavior

### With API Key:
```
Job xxx: starting transcription
Job xxx: transcription completed, 25 segments
Job xxx: starting highlight analysis with Gemini
Job xxx: found 3 highlights
Job xxx: completed successfully
```

### Without API Key:
```
Job xxx: starting transcription
Job xxx: transcription completed, 25 segments
Job xxx: Skipping highlight analysis (no Gemini API key)
Job xxx: completed successfully
```

Both cases work - without API key just skips highlight detection.

---

## 💰 Rate Limits (Free Tier)

- **15 requests per minute**
- **1500 requests per day**

For your use case: **1500 videos per day = plenty!**

If exceeded:
- Wait 1 minute (RPM limit)
- Or wait until next day (daily limit)

---

## 🐛 Troubleshooting

### Error: "gemini API key not configured"
- Check `.env` file exists
- Check `GEMINI_API_KEY=` has actual key
- Restart server after adding key

### Error: "failed to create gemini client"
- Check internet connection
- Verify API key is correct (copy-paste again)
- Check: https://aistudio.google.com/app/apikey to verify key is active

### Error: "API quota exceeded"
- Wait 1 minute for RPM reset
- Or wait for daily reset (00:00 UTC)
- Check usage: https://console.cloud.google.com/apis/api/generativelanguage.googleapis.com/quotas

### Highlights not appearing
- Check server logs for "Gemini analysis failed"
- Verify transcript has content
- Check Gemini response in logs

---

## ✅ You're Done!

Gemini is now integrated. Every video will:
1. Download ✓
2. Extract audio ✓
3. Transcribe ✓
4. **Analyze highlights with AI** 🆕✓
5. Ready for clipping (Phase 3)

Next: Implement video clipping to cut the highlights! 🎬
