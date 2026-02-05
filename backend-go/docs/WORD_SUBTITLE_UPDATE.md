# Word-by-Word Subtitle Update 🎬

## Changes Made

### Before (Full Sentences)
```
1
00:00:00,000 --> 00:00:05,000
Alasan pertama kenapa lu masih miskin adalah lu nyaman untuk jadi seorang loser
```

### After (Word-by-Word)
```
1
00:00:00,000 --> 00:00:00,500
Alasan pertama

2
00:00:00,500 --> 00:00:01,000
kenapa lu

3  
00:00:01,000 --> 00:00:01,500
masih miskin
```

## Implementation

**Key features:**
- ✅ Shows **1-3 words at a time**
- ✅ **Auto-calculates** timing per word
- ✅ **Synced with speech** based on segment duration
- ✅ **Larger font** (32px) for better readability
- ✅ **Bolder text** (Arial Bold) for impact
- ✅ **Thicker outline** (3px) for visibility

## Algorithm

```go
// For each transcript segment:
1. Split text into words
2. Calculate: timePerWord = segmentDuration / wordCount
3. Group into 1-3 words per subtitle
4. Assign timing: wordStart = segStart + (wordIndex * timePerWord)
5. Create SRT entry for each group
```

**Example:**
```
Segment: "kenapa lu masih miskin" (2.0 seconds, 4 words)
Time per word: 2.0s / 4 = 0.5s

Result:
0.0-1.0s: "kenapa lu"     (2 words)
1.0-2.0s: "masih miskin"  (2 words)
```

## Benefits

**Engagement:**
- ✅ More dynamic and eye-catching
- ✅ Easier to follow along
- ✅ Better for fast-paced content
- ✅ TikTok/Reels style

**Accessibility:**
- ✅ Easier to read (less text at once)
- ✅ Better sync with audio
- ✅ More natural flow

## Testing

Restart server dan process video baru:
```powershell
go run main.go
```

**Verify:**
1. Subtitles show 1-3 words at a time ✓
2. Words change quickly (synced with speech) ✓
3. Text is large and bold ✓
4. Easy to read on vertical video ✓

## Customization

Adjust in `subtitle.go`:

```go
// Show fewer/more words per subtitle
wordsPerSubtitle := 2  // Change to 1 or 3

// Font size
FontSize: 32  // Make bigger/smaller

// Outline thickness
OutlineWidth: 3  // Adjust for visibility
```
