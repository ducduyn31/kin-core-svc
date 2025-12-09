package messaging

import "github.com/google/uuid"

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImage    ContentType = "image"
	ContentTypeVideo    ContentType = "video"
	ContentTypeAudio    ContentType = "audio"
	ContentTypeFile     ContentType = "file"
	ContentTypeLocation ContentType = "location"
	ContentTypeSticker  ContentType = "sticker"
)

type Content struct {
	Type     ContentType `json:"type"`
	Text     *string     `json:"text,omitempty"`
	MediaID  *uuid.UUID  `json:"media_id,omitempty"`
	MediaURL *string     `json:"media_url,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}

type Metadata struct {
	FileName  *string `json:"file_name,omitempty"`
	FileSize  *int64  `json:"file_size,omitempty"`
	MimeType  *string `json:"mime_type,omitempty"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Duration  *int    `json:"duration,omitempty"` // Seconds for audio/video
	Thumbnail *string `json:"thumbnail,omitempty"`

	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	PlaceName *string  `json:"place_name,omitempty"`
	Address   *string  `json:"address,omitempty"`

	StickerID   *string `json:"sticker_id,omitempty"`
	StickerPack *string `json:"sticker_pack,omitempty"`
}

func NewTextContent(text string) Content {
	return Content{
		Type: ContentTypeText,
		Text: &text,
	}
}

func NewImageContent(mediaID uuid.UUID, url string, width, height int) Content {
	return Content{
		Type:     ContentTypeImage,
		MediaID:  &mediaID,
		MediaURL: &url,
		Metadata: &Metadata{
			Width:  &width,
			Height: &height,
		},
	}
}

func NewVideoContent(mediaID uuid.UUID, url string, duration int, thumbnail *string) Content {
	return Content{
		Type:     ContentTypeVideo,
		MediaID:  &mediaID,
		MediaURL: &url,
		Metadata: &Metadata{
			Duration:  &duration,
			Thumbnail: thumbnail,
		},
	}
}

func NewAudioContent(mediaID uuid.UUID, url string, duration int) Content {
	return Content{
		Type:     ContentTypeAudio,
		MediaID:  &mediaID,
		MediaURL: &url,
		Metadata: &Metadata{
			Duration: &duration,
		},
	}
}

func NewFileContent(mediaID uuid.UUID, url, fileName string, fileSize int64, mimeType string) Content {
	return Content{
		Type:     ContentTypeFile,
		MediaID:  &mediaID,
		MediaURL: &url,
		Metadata: &Metadata{
			FileName: &fileName,
			FileSize: &fileSize,
			MimeType: &mimeType,
		},
	}
}

func NewLocationContent(lat, lng float64, placeName, address *string) Content {
	return Content{
		Type: ContentTypeLocation,
		Metadata: &Metadata{
			Latitude:  &lat,
			Longitude: &lng,
			PlaceName: placeName,
			Address:   address,
		},
	}
}

func IsValidContentType(ct ContentType) bool {
	switch ct {
	case ContentTypeText, ContentTypeImage, ContentTypeVideo, ContentTypeAudio,
		ContentTypeFile, ContentTypeLocation, ContentTypeSticker:
		return true
	default:
		return false
	}
}
