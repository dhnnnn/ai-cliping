package utils

func ExtractAudio(videoPath, audioPath string) error {
	return RunCommand(
		FFmpegBinary(),
		"-y",
		"-i", videoPath,
		"-vn",
		"-ac", "1",
		"-ar", "16000",
		audioPath,
	)
}
