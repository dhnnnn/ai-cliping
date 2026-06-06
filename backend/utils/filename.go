package utils

import (
	"regexp"
	"strings"
)

// invalidFilenameChars cocok dengan karakter yang tidak aman untuk nama file
// di Windows maupun Linux (termasuk \ / : * ? " < > |) dan karakter kontrol.
var invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

// SanitizeFilename mengubah judul dari AI menjadi nama file yang aman.
// Contoh: `Cara Bikin Kopi: Tips & Trik!` -> `Cara Bikin Kopi - Tips & Trik`
func SanitizeFilename(name string) string {
	// Ganti karakter ilegal dengan spasi
	cleaned := invalidFilenameChars.ReplaceAllString(name, " ")

	// Rapikan spasi berlebih jadi satu spasi
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	// Hilangkan titik/spasi di awal & akhir (Windows tidak suka trailing dot/space)
	cleaned = strings.Trim(cleaned, " .")

	// Batasi panjang agar tidak melebihi limit filesystem
	const maxLen = 100
	if len(cleaned) > maxLen {
		cleaned = strings.TrimSpace(cleaned[:maxLen])
	}

	// Fallback kalau hasilnya kosong
	if cleaned == "" {
		cleaned = "clip"
	}

	return cleaned
}
