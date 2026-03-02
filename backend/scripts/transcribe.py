import sys
import json
import warnings
warnings.filterwarnings("ignore")

def transcribe_audio(audio_path, model_name="base"):
    """
    Transcribe audio file using faster-whisper (lebih ringan dari openai-whisper)
    
    Args:
        audio_path: Path to audio file (WAV recommended)
        model_name: Whisper model size (tiny, base, small, medium, large)
    """
    try:
        from faster_whisper import WhisperModel

        # Gunakan CPU dengan int8 quantization agar lebih hemat memori
        model = WhisperModel(model_name, device="cpu", compute_type="int8")

        # Transcribe
        segments_gen, info = model.transcribe(
            audio_path,
            language="id",       # Bahasa Indonesia, ubah ke None untuk auto-detect
            word_timestamps=True
        )

        # Format output
        segments = []
        for segment in segments_gen:
            segments.append({
                "start": segment.start,
                "end":   segment.end,
                "text":  segment.text.strip()
            })

        output = {
            "success":  True,
            "segments": segments,
            "language": info.language
        }

        return json.dumps(output, ensure_ascii=False)

    except Exception as e:
        error_output = {
            "success": False,
            "error":   str(e)
        }
        return json.dumps(error_output)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"success": False, "error": "No audio file provided"}))
        sys.exit(1)

    audio_path = sys.argv[1]
    model_name = sys.argv[2] if len(sys.argv) > 2 else "base"

    result = transcribe_audio(audio_path, model_name)
    print(result)
