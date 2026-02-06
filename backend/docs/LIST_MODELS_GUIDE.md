# List Available Gemini Models

Gunakan script ini untuk melihat model Gemini mana yang available untuk API key kamu.

## 🐍 Option 1: Python Script (Recommended)

### Install Dependencies
```powershell
pip install google-generativeai
```

### Run Script
```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Option A: Use API key dari environment (.env)
python scripts\list_gemini_models.py

# Option B: Pass API key directly
python scripts\list_gemini_models.py AIzaSyYourKeyHere
```

---

## 🔷 Option 2: Go Script

### Run Script
```powershell
cd d:\.dev\golang\ai-clipping\backend-go\scripts

# Option A: Use API key dari environment
go run test_gemini_models.go

# Option B: Pass API key directly
go run test_gemini_models.go AIzaSyYourKeyHere
```

---

## 📊 Expected Output

```
============================================================
GEMINI AVAILABLE MODELS
============================================================

Models that support generateContent:
------------------------------------------------------------
✓ models/gemini-pro
  Display name: Gemini Pro
  Description: The best model for scaling across a wide range of tasks

✓ models/gemini-1.5-flash-001
  Display name: Gemini 1.5 Flash
  Description: Fast and versatile multimodal model

✓ models/gemini-1.5-pro-001
  Display name: Gemini 1.5 Pro
  Description: Mid-size multimodal model optimized for performance

============================================================
Total: 3 models support generateContent
============================================================
```

---

## 🔧 Next Steps

1. **Run script** untuk lihat available models
2. **Copy model name** yang mau dipakai (contoh: `models/gemini-pro`)
3. **Update `analysis/gemini.go` line 50:**
   ```go
   model := client.GenerativeModel("gemini-pro")  // Ganti dengan model yang available
   ```
4. **Restart server** dan test lagi

---

## 💡 Model Recommendations

| Model | Speed | Quality | Use Case |
|-------|-------|---------|----------|
| **gemini-pro** | ⭐⭐⭐ | ⭐⭐⭐⭐ | Balanced, stable (RECOMMENDED) |
| **gemini-1.5-flash-001** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | Fastest, good for highlighting |
| **gemini-1.5-pro-001** | ⭐⭐ | ⭐⭐⭐⭐⭐ | Best quality, slower |

---

## ❌ Troubleshooting

### Error: "Invalid API key"
- Double check API key di `.env` file
- Make sure no extra spaces or quotes
- Verify key starts with `AIza`

### Error: "Module not found"
```powershell
# Python
pip install google-generativeai

# Go - dependencies already installed via go mod
```

### No models shown
- Check internet connection
- Verify API key is active at https://aistudio.google.com/app/apikey
- Check API quota not exceeded
