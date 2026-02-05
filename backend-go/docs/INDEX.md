# Documentation Index

Welcome to AI Clipping Backend documentation! 📚

## 📖 Quick Start Guides

- **[README](../README.md)** - Main project overview and quick setup
- **[Quick Start: Gemini](QUICK_START_GEMINI.md)** - 5-minute setup for AI highlight detection
- **[Installation Guide](INSTALL_DEPS.md)** - Install all dependencies

## 🔧 Setup Guides

- **[Whisper Setup](WHISPER_SETUP.md)** - Install and configure Whisper for transcription
- **[Python Setup](PYTHON_SETUP.md)** - Fix Python PATH issues on Windows
- **[Gemini Setup](GEMINI_SETUP.md)** - Detailed Gemini API configuration
- **[Debug .env File](DEBUG_ENV.md)** - Troubleshoot environment variable loading

## 📊 User Guides

- **[Testing Guide](TESTING_GUIDE.md)** - How to test with Postman
- **[Saved Files Structure](SAVED_FILES.md)** - Where results are stored

## 🏗️ Project Structure

```
backend-go/
├── README.md              # Main documentation (GitHub landing page)
├── main.go               # Application entry point
├── .env                  # Environment variables (not in git)
├── .env.example          # Template for .env
│
├── docs/                 # 📁 All documentation files
│   ├── INDEX.md          # This file
│   ├── QUICK_START_GEMINI.md
│   ├── GEMINI_SETUP.md
│   ├── WHISPER_SETUP.md
│   ├── PYTHON_SETUP.md
│   ├── TESTING_GUIDE.md
│   ├── SAVED_FILES.md
│   ├── DEBUG_ENV.md
│   └── INSTALL_DEPS.md
│
├── analysis/             # AI analysis logic
├── handlers/             # API endpoints
├── models/               # Data models
├── pipeline/             # Job processing
├── queue/                # Job queue
├── utils/                # Utilities
├── scripts/              # Python scripts
└── storage/              # File storage
```

## 🚀 Workflow

1. **Start Here:** [README](../README.md) - Project overview
2. **Setup:** [INSTALL_DEPS.md](INSTALL_DEPS.md) - Install dependencies
3. **Configure:** [GEMINI_SETUP.md](GEMINI_SETUP.md) - Get API key
4. **Test:** [TESTING_GUIDE.md](TESTING_GUIDE.md) - Use Postman
5. **Understand:** [SAVED_FILES.md](SAVED_FILES.md) - Check results

## 🐛 Troubleshooting

- Python not found? → [PYTHON_SETUP.md](PYTHON_SETUP.md)
- .env not loading? → [DEBUG_ENV.md](DEBUG_ENV.md)
- Whisper issues? → [WHISPER_SETUP.md](WHISPER_SETUP.md)

## 📝 Notes

All documentation uses GitHub Flavored Markdown for best compatibility.
