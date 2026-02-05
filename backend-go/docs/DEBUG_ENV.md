# Debug: .env File Issues

## Problem
Error: "No .env file found" padahal file ada

### Common Causes:

1. **File Encoding Issue** (Most Common)
   - Windows create file dengan UTF-16LE encoding
   - `godotenv` library hanya support UTF-8
   - Solution: Re-create file dengan UTF-8

2. **File Location**
   - File harus di root project (same folder as main.go)
   - Check: `d:\.dev\golang\ai-clipping\backend-go\.env`

3. **File Content**
   - Typo di key name
   - Extra spaces
   - Missing `=`

---

## Quick Fix

### Option 1: Re-create .env File (Recommended)

```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Delete old .env
del .env

# Create new .env with UTF-8 encoding
Set-Content -Path ".env" -Value "GEMINI_API_KEY=your-api-key-here" -Encoding UTF8
```

Replace `your-api-key-here` dengan actual API key!

### Option 2: Use Notepad++ or VS Code

1. Open `.env` file
2. Check encoding di bottom-right corner
3. If UTF-16: Convert to UTF-8
   - Notepad++: Encoding → Convert to UTF-8
   - VS Code: Click encoding → Save with Encoding → UTF-8

### Option 3: Set Environment Variable Directly

```powershell
# Temporary (current session)
$env:GEMINI_API_KEY="your-api-key-here"
go run main.go

# Permanent
[System.Environment]::SetEnvironmentVariable('GEMINI_API_KEY', 'your-api-key-here', 'User')
```

---

## Verify .env Working

After fix, run:

```powershell
go run main.go
```

**Should see:**
```
2026/02/04 13:50:02 Server running on :8080
```

**Should NOT see:**
```
❌ No .env file found
❌ WARNING: GEMINI_API_KEY not set
```

If still seeing warnings → API key not loaded correctly.

---

## Test API Key Loaded

Quick test in PowerShell:

```powershell
# Show current environment variable
echo $env:GEMINI_API_KEY

# Should output your API key
# If empty → not loaded
```

---

## Alternative: Edit main.go for Debug

Temporarily add debug logging:

```go
func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
        log.Printf("DEBUG: Error details: %v", err) // ADD THIS
    }

    geminiAPIKey := os.Getenv("GEMINI_API_KEY")
    log.Printf("DEBUG: API Key loaded: %v", geminiAPIKey != "") // ADD THIS
    
    if geminiAPIKey == "" {
        log.Println("WARNING: GEMINI_API_KEY not set")
    }
    ...
}
```

This will show exact error!

---

## My Recommendation

**Easiest solution:**

```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Remove old file
Remove-Item .env -Force

# Create new with PowerShell (guaranteed UTF-8)
@"
GEMINI_API_KEY=AIzaSyYourActualKeyHere
"@ | Out-File -FilePath .env -Encoding utf8 -NoNewline

# Verify file created
cat .env

# Run server
go run main.go
```

Ganti `AIzaSyYourActualKeyHere` dengan API key asli kamu!
