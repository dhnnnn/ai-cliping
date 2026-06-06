package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// tokenFile adalah lokasi penyimpanan OAuth token (di dalam folder storage
// agar ikut ter-persist lewat volume Docker).
const tokenFile = "storage/youtube_token.json"

// mu melindungi akses baca/tulis file token.
var mu sync.Mutex

// Uploader membungkus konfigurasi OAuth2 untuk YouTube Data API v3.
type Uploader struct {
	config *oauth2.Config
}

// NewUploader membuat Uploader dari environment variables.
// Mengembalikan error jika kredensial OAuth belum diset.
func NewUploader() (*Uploader, error) {
	clientID := os.Getenv("YOUTUBE_CLIENT_ID")
	clientSecret := os.Getenv("YOUTUBE_CLIENT_SECRET")
	redirectURL := os.Getenv("YOUTUBE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("YOUTUBE_CLIENT_ID / YOUTUBE_CLIENT_SECRET belum diset")
	}
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/youtube/callback"
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{youtube.YoutubeUploadScope},
		Endpoint:     google.Endpoint,
	}

	return &Uploader{config: config}, nil
}

// AuthURL menghasilkan URL consent Google untuk memulai otorisasi.
// AccessTypeOffline + prompt consent memastikan kita dapat refresh token.
func (u *Uploader) AuthURL(state string) string {
	return u.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

// ExchangeAndSave menukar authorization code dengan token lalu menyimpannya.
func (u *Uploader) ExchangeAndSave(ctx context.Context, code string) error {
	token, err := u.config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("gagal menukar code dengan token: %w", err)
	}
	return saveToken(token)
}

// IsAuthorized mengecek apakah token sudah tersimpan.
func (u *Uploader) IsAuthorized() bool {
	_, err := loadToken()
	return err == nil
}

// UploadResult berisi info video yang berhasil diupload.
type UploadResult struct {
	VideoID string `json:"videoId"`
	URL     string `json:"url"`
	Title   string `json:"title"`
}

// UploadOptions adalah parameter untuk satu kali upload.
type UploadOptions struct {
	FilePath      string
	Title         string
	Description   string
	Tags          []string
	CategoryID    string // default "22" (People & Blogs)
	PrivacyStatus string // "private" | "unlisted" | "public"
}

// Upload mengunggah satu file video ke channel YouTube yang terotorisasi.
func (u *Uploader) Upload(ctx context.Context, opts UploadOptions) (*UploadResult, error) {
	token, err := loadToken()
	if err != nil {
		return nil, fmt.Errorf("belum terotorisasi, buka /api/youtube/auth dulu: %w", err)
	}

	// client otomatis me-refresh access token saat kedaluwarsa.
	client := u.config.Client(ctx, token)

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat youtube service: %w", err)
	}

	// Default value
	if opts.CategoryID == "" {
		opts.CategoryID = "22"
	}
	if opts.PrivacyStatus == "" {
		opts.PrivacyStatus = "public"
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       sanitizeYouTubeTitle(opts.Title),
			Description: opts.Description,
			Tags:        opts.Tags,
			CategoryId:  opts.CategoryID,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: opts.PrivacyStatus,
		},
	}

	file, err := os.Open(opts.FilePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file clip %s: %w", opts.FilePath, err)
	}
	defer file.Close()

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)
	resp, err := call.Media(file).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("upload ke youtube gagal: %w", err)
	}

	return &UploadResult{
		VideoID: resp.Id,
		URL:     "https://youtu.be/" + resp.Id,
		Title:   upload.Snippet.Title,
	}, nil
}

// sanitizeYouTubeTitle membuang karakter yang ditolak YouTube (< dan >)
// serta membatasi panjang judul ke 100 karakter.
func sanitizeYouTubeTitle(title string) string {
	cleaned := strings.NewReplacer("<", "", ">", "").Replace(title)
	cleaned = strings.TrimSpace(cleaned)
	if len([]rune(cleaned)) > 100 {
		cleaned = string([]rune(cleaned)[:100])
	}
	if cleaned == "" {
		cleaned = "Untitled Clip"
	}
	return cleaned
}

// saveToken menyimpan OAuth token ke file JSON.
func saveToken(token *oauth2.Token) error {
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(tokenFile), 0755); err != nil {
		return fmt.Errorf("gagal membuat folder token: %w", err)
	}

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("gagal marshal token: %w", err)
	}

	if err := os.WriteFile(tokenFile, data, 0600); err != nil {
		return fmt.Errorf("gagal menyimpan token: %w", err)
	}
	return nil
}

// loadToken membaca OAuth token dari file JSON.
func loadToken() (*oauth2.Token, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
