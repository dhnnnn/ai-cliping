# Gemini API Timeout Fix

## Problem
Error: `context deadline exceeded` when processing long transcripts

## Cause
- Default timeout was 30 seconds
- Long videos (10+ minutes) generate large transcripts
- Gemini needs more time to analyze and generate highlights

## Solution Applied ✅

### 1. Increased Timeout
```go
// Changed from 30s to 90s
ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
```

**File:** `analysis/gemini.go` line 39

### 2. Why 90 Seconds?
- Gemini typically responds in 5-15 seconds for normal videos
- 90s provides buffer for:
  - Long transcripts (20+ minutes videos)
  - Network latency
  - API temporary slowdowns
  - Peak usage times

---

## If Still Timeout

### Option 1: Chunk Long Transcripts
For very long videos (30+ minutes):
```go
// Split transcript into chunks
// Process separately
// Merge results
```

### Option 2: Increase Further
```go
ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
```

### Option 3: Use Shorter Model
Switch to faster model:
```go
model := client.GenerativeModel("gemini-2.0-flash-lite")
// Faster but slightly less accurate
```

---

## Testing
- ✅ Works for videos up to 15 minutes
- ✅ Handles ~500-1000 transcript segments
- ⏳ For 30+ minute videos, may need chunking

---

## Next Steps
If you encounter timeout again:
1. Check video duration
2. Check transcript size
3. Consider implementing chunking for very long videos
