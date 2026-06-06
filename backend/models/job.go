package models

type JobStatus string

const (
	StatusQueued       = "queued"
	StatusDownloading  = "downloading"
	StatusDownloaded   = "downloaded"
	StatusAudioExtract = "audio_extracting"
	StatusAudioReady   = "audio_ready"
	StatusTranscribing = "transcribing"
	StatusTranscribed  = "transcribed"
	StatusAnalyzing    = "analyzing"
	StatusClipping     = "clipping"
	StatusCompleted    = "completed"
	StatusFailed       = "failed"
)

type VideoType string

const (
	VideoTypeTutorial VideoType = "tutorial"
	VideoTypeGaming   VideoType = "gaming"
	VideoTypePodcast  VideoType = "podcast"
	VideoTypeGeneral  VideoType = "general"
)

type ClipPreferences struct {
	MinDuration int    `json:"minDuration"` // seconds
	MaxDuration int    `json:"maxDuration"` // seconds
	MaxClips    int    `json:"maxClips"`
	AspectRatio string `json:"aspectRatio"` // "9:16", "16:9", "1:1"
}

type TranscriptSegment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type Highlight struct {
	Start    float64  `json:"start"`
	End      float64  `json:"end"`
	Score    float64  `json:"score"`
	Title    string   `json:"title"`    // AI-generated, catchy title siap diupload ke YouTube/TikTok
	Reason   string   `json:"reason"`   // e.g., "keyword_match", "volume_spike"
	Keywords []string `json:"keywords"` // matched keywords if any
}

type Job struct {
	ID          string              `json:"id"`
	URL         string              `json:"url"`
	Status      JobStatus           `json:"status"`
	Result      string              `json:"result"`
	VideoType   VideoType           `json:"videoType"`
	Preferences ClipPreferences     `json:"preferences"`
	Transcript  []TranscriptSegment `json:"transcript,omitempty"`
	Highlights  []Highlight         `json:"highlights,omitempty"`
	ClipPaths   []string            `json:"clipPaths,omitempty"`
}
