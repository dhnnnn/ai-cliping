# Complete Setup: Install All Dependencies

Run these commands in order:

```powershell
cd d:\.dev\golang\ai-clipping\backend-go

# Install ALL required dependencies
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option
go get github.com/joho/godotenv
go get github.com/google/uuid

# Clean up and update go.mod
go mod tidy

# Verify dependencies installed
go list -m all | findstr godotenv
```

Should output:
```
github.com/joho/godotenv v1.5.1
```

Then run server:
```powershell
go run main.go
```
