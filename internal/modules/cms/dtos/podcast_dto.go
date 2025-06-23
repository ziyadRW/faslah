package cms

import (
	podcastModels "github.com/ziyadrw/faslah/internal/modules/cms/models"
	"mime/multipart"
	"time"
)

// CreateContentRequest represents the request body for creating a podcast
type CreateContentRequest struct {
	MediaURL  string `json:"media_url,omitempty" validate:"required_without=SourceURL" message:"يجب توفير رابط الوسائط أو رابط المصدر"`
	SourceURL string `json:"source_url,omitempty" validate:"required_without=MediaURL" message:"يجب توفير رابط الوسائط أو رابط المصدر"`

	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
	DurationSecs int        `json:"duration_secs,omitempty"`
}

// UpdateContentRequest represents the request body for updating a podcast
type UpdateContentRequest struct {
	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	MediaURL     string     `json:"media_url,omitempty"`
	SourceURL    string     `json:"source_url,omitempty"`
	DurationSecs int        `json:"duration_secs,omitempty"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
}

// PodcastResponse represents the response body for a podcast
type PodcastResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Tags         []string   `json:"tags"`
	MediaURL     string     `json:"media_url"`
	SourceURL    string     `json:"source_url"`
	DurationSecs int        `json:"duration_secs"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// MediaUploadResponse represents the response body for media upload
type MediaUploadResponse struct {
	MediaURL string `json:"media_url"`
}

// YouTubeMetadataResponse represents the response body for YouTube metadata extraction
type YouTubeMetadataResponse struct {
	VideoFile    []byte `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	DurationSecs int    `json:"duration_secs"`
}

// FetchYouTubeRequest represents the request body for fetching YouTube metadata
type FetchYouTubeRequest struct {
	YoutubeURL string `json:"youtube_url" form:"youtube_url" validate:"required" message:"يجب توفير رابط يوتيوب"`
}

// CreatePodcastRequest represents the request body for creating a podcast
type CreatePodcastRequest struct {
	SourceURL   string                `json:"source_url,omitempty" form:"source_url" validate:"omitempty,required_without=File"`
	File        *multipart.FileHeader `json:"-" form:"file" validate:"omitempty,required_without=SourceURL"`
	Title       string                `json:"title,omitempty" form:"title" validate:"required_if=File.Filename .,omitempty"`
	Description string                `json:"description,omitempty" form:"description" validate:"required_if=File.Filename .,omitempty"`
	Tags        []string              `json:"tags,omitempty" form:"tags" validate:"required_if=File.Filename .,omitempty"`
	PublishedAt *time.Time            `json:"published_at,omitempty" form:"published_at"`
}

func MapPodcastToDTO(podcast *podcastModels.Podcast) PodcastResponse {
	return PodcastResponse{
		ID:           podcast.ID.String(),
		UserID:       podcast.UserID,
		Title:        podcast.Title,
		Description:  podcast.Description,
		Tags:         podcast.Tags,
		MediaURL:     podcast.MediaURL,
		SourceURL:    podcast.SourceURL,
		DurationSecs: podcast.DurationSecs,
		PublishedAt:  podcast.PublishedAt,
		CreatedAt:    podcast.CreatedAt,
		UpdatedAt:    podcast.UpdatedAt,
	}
}
