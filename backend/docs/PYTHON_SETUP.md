# Python PATH Setup for Windows

## Issue
Error: `exec: "python": executable file not found in %PATH%`

## Quick Fix Options

### Option 1: Verify Python Installation (Recommended)

```powershell
# Test which Python command works
py --version
python --version
python3 --version
```

Salah satu command di atas harus work. Jika tidak, install Python terlebih dahulu.

### Option 2: Add Python to PATH

**Cara Manual:**

1. Cari lokasi Python install:
   ```powershell
   # Biasanya di salah satu lokasi ini:
   C:\Python312\
   C:\Users\<YourUsername>\AppData\Local\Programs\Python\Python312\
   ```

2. Add ke PATH:
   - Open "Environment Variables" (ketik di Start menu)
   - Edit "Path" di System variables
   - Click "New" dan tambahkan path Python, contoh:
     ```
     C:\Users\YourUsername\AppData\Local\Programs\Python\Python312
     C:\Users\YourUsername\AppData\Local\Programs\Python\Python312\Scripts
     ```
   - Click OK semua dialog
   - **Restart terminal/CMD** atau restart computer

3. Test ulang:
   ```powershell
   python --version
   ```

**Cara Temporary (untuk sesi terminal saat ini):**

```powershell
# Ganti dengan path Python yang sebenarnya
$env:PATH += ";C:\Users\YourUsername\AppData\Local\Programs\Python\Python312"

# Test
python --version

# Jalankan server
go run main.go
```

### Option 3: Install Python dengan Checkbox "Add to PATH"

Jika Python belum terinstall atau mau re-install:

1. Download dari: https://www.python.org/downloads/
2. **PENTING:** Saat install, **centang "Add Python to PATH"**
3. Complete installation
4. Restart terminal

### Option 4: Use Python Launcher (Windows)

Windows Python biasanya install dengan `py` launcher:

```powershell
# Gunakan py instead of python
py --version
pip install openai-whisper
```

Code sudah diupdate untuk try `py`, `python`, dan `python3` secara otomatis!

## Verify Whisper Installation

Setelah Python terdeteksi:

```powershell
# Check Whisper installed
pip show openai-whisper

# Or dengan py
py -m pip show openai-whisper

# Install jika belum ada
pip install openai-whisper
# or
py -m pip install openai-whisper
```

## Test Python Script Directly

```powershell
# Test transcribe script manual
cd d:\.dev\golang\ai-clipping\backend-go

# Try dengan py
py scripts/transcribe.py

# Or dengan python
python scripts/transcribe.py
```

Expected error (normal karena no audio file):
```json
{"success": false, "error": "No audio file provided"}
```

## After Fix

1. Restart terminal yang menjalankan Go server
2. Atau restart `go run main.go`
3. Test ulang di Postman

Server seharusnya sekarang bisa call Python Whisper script!
