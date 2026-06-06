package utils

import (
	"log"
	"os/exec"
	"runtime"
)

// FFmpegBinary mengembalikan nama binary ffmpeg sesuai OS.
// Windows pakai ./ffmpeg.exe (binary lokal), Linux/Docker pakai ffmpeg dari PATH.
func FFmpegBinary() string {
	if runtime.GOOS == "windows" {
		return "./ffmpeg.exe"
	}
	return "ffmpeg"
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	// Always log output for debugging
	if len(output) > 0 {
		log.Println(string(output))
	}

	if err != nil {
		log.Printf("Command failed: %v", err)
		return err
	}

	return nil
}
