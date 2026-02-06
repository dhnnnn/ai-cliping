package utils

func ExtractAudio(videoPath, audioPath string) error {
	return RunCommand(
		"./ffmpeg.exe",
		"-y",
		"-i", videoPath,
		"-vn",
		"-ac", "1",
		"-ar", "16000",
		audioPath,
	)
}
