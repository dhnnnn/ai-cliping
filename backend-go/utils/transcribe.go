package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ai-clipping-backend/models"
)

type WhisperResponse struct {
	Success  bool                       `json:"success"`
	Segments []models.TranscriptSegment `json:"segments"`
	Language string                     `json:"language"`
	Error    string                     `json:"error"`
}

// TranscribeAudio transcribes audio file using Whisper Python script
// modelSize: "tiny", "base", "small", "medium", "large"
// "base" is recommended for balance between speed and accuracy
func TranscribeAudio(audioPath string, modelSize string) ([]models.TranscriptSegment, error) {
	if modelSize == "" {
		modelSize = "base"
	}

	// Try different Python commands (Windows compatibility)
	pythonCommands := []string{"py", "python", "python3"}

	var output []byte
	var err error
	var lastErr error

	for _, pythonCmd := range pythonCommands {
		cmd := exec.Command(pythonCmd, "./scripts/transcribe.py", audioPath, modelSize)
		output, err = cmd.CombinedOutput()

		if err == nil {
			// Command executed successfully, break loop
			break
		}

		// Store error for reporting if all commands fail
		lastErr = fmt.Errorf("%s command failed: %v", pythonCmd, err)
	}

	if err != nil {
		return nil, fmt.Errorf("whisper execution failed (tried py, python, python3): %v, output: %s", lastErr, string(output))
	}

	// Parse JSON response
	// Note: Whisper outputs progress bars and warnings to stdout
	// The actual JSON result is always on the last line
	outputStr := string(output)
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	// Get the last non-empty line (should be JSON)
	var jsonLine string
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" && (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")) {
			jsonLine = trimmed
			break
		}
	}

	if jsonLine == "" {
		return nil, fmt.Errorf("no JSON output found in whisper response, full output: %s", outputStr)
	}

	var response WhisperResponse
	if err := json.Unmarshal([]byte(jsonLine), &response); err != nil {
		return nil, fmt.Errorf("failed to parse whisper response: %v, json line: %s", err, jsonLine)
	}

	if !response.Success {
		return nil, fmt.Errorf("transcription failed: %s", response.Error)
	}

	return response.Segments, nil
}

// SaveTranscript saves transcript segments to a JSON file
func SaveTranscript(segments []models.TranscriptSegment, filePath string) error {
	data, err := json.MarshalIndent(segments, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transcript: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
