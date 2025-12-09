package media

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeFile  MediaType = "file"
)

type Media struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Type         MediaType `json:"type"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	StorageKey   string    `json:"storage_key"`
	URL          string    `json:"url"`
	ThumbnailKey *string   `json:"thumbnail_key,omitempty"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	Width        *int      `json:"width,omitempty"`
	Height       *int      `json:"height,omitempty"`
	Duration     *int      `json:"duration,omitempty"` // Seconds for audio/video
	CreatedAt    time.Time `json:"created_at"`
}

func NewMedia(userID uuid.UUID, fileName string, fileSize int64, mimeType string) *Media {
	id := uid.New()
	ext := filepath.Ext(fileName)
	storageKey := generateStorageKey(userID, id, ext)

	return &Media{
		ID:         id,
		UserID:     userID,
		Type:       MediaTypeFromMime(mimeType),
		FileName:   fileName,
		FileSize:   fileSize,
		MimeType:   mimeType,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
	}
}

func (m *Media) SetURL(url string) {
	m.URL = url
}

func (m *Media) SetThumbnail(key, url string) {
	m.ThumbnailKey = &key
	m.ThumbnailURL = &url
}

func (m *Media) SetDimensions(width, height int) {
	m.Width = &width
	m.Height = &height
}

func (m *Media) SetDuration(duration int) {
	m.Duration = &duration
}

func (m *Media) IsImage() bool {
	return m.Type == MediaTypeImage
}

func (m *Media) IsVideo() bool {
	return m.Type == MediaTypeVideo
}

func (m *Media) IsAudio() bool {
	return m.Type == MediaTypeAudio
}

func generateStorageKey(userID, mediaID uuid.UUID, ext string) string {
	now := time.Now()
	return strings.Join([]string{
		userID.String(),
		now.Format("2006/01/02"),
		mediaID.String() + ext,
	}, "/")
}

func MediaTypeFromMime(mimeType string) MediaType {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return MediaTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return MediaTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return MediaTypeAudio
	default:
		return MediaTypeFile
	}
}

var AllowedMimeTypes = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"image/webp":         true,
	"image/heic":         true,
	"image/heif":         true,
	"video/mp4":          true,
	"video/quicktime":    true,
	"video/webm":         true,
	"audio/mpeg":         true,
	"audio/mp4":          true,
	"audio/wav":          true,
	"audio/ogg":          true,
	"audio/webm":         true,
	"audio/aac":          true,
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"text/plain": true,
}

func IsAllowedMimeType(mimeType string) bool {
	return AllowedMimeTypes[mimeType]
}

var MaxFileSizes = map[MediaType]int64{
	MediaTypeImage: 10 * 1024 * 1024,  // 10 MB
	MediaTypeVideo: 100 * 1024 * 1024, // 100 MB
	MediaTypeAudio: 50 * 1024 * 1024,  // 50 MB
	MediaTypeFile:  25 * 1024 * 1024,  // 25 MB
}

func IsWithinSizeLimit(mediaType MediaType, size int64) bool {
	maxSize, ok := MaxFileSizes[mediaType]
	if !ok {
		return false
	}
	return size <= maxSize
}
