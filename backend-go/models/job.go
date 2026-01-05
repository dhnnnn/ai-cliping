package models

type JobStatus string

const (
	StatusQueued         = "queued"
	StatusDownloading    = "downloading"
	StatusDownloaded     = "downloaded"
	StatusAudioExtract   = "audio_extracting"
	StatusAudioReady     = "audio_ready"
	StatusFailed         = "failed"
)

type Job struct {
	ID     string
	URL    string
	Status JobStatus
	Result string
}
