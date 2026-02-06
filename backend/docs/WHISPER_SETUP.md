# Whisper Setup Requirements

## Install Python Dependencies

```bash
# Install Python 3.8+ if not installed
# Download from: https://www.python.org/downloads/

# Install Whisper
pip install openai-whisper

# Install ffmpeg (required by Whisper for audio processing)
# Already have ffmpeg.exe in project, but make sure Python can find it
```

## Test Whisper Installation

```bash
# Test the transcribe script
python scripts/transcribe.py storage/audio/test.wav base
```

## Model Sizes

- **tiny**: ~1GB, fastest, least accurate
- **base**: ~1GB, good balance **(RECOMMENDED)**
- **small**: ~2GB, better accuracy
- **medium**: ~5GB, great accuracy
- **large**: ~10GB, best accuracy

First run will download the model (~1GB for base).

## Language Support

Currently set to Indonesian (`id`). To auto-detect or use other language:
- Edit `scripts/transcribe.py`
- Change `language="id"` to `language=None` for auto-detect
- Or set specific language code (e.g., `"en"` for English)
