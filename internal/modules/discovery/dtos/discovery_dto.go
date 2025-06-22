package discovery

import (
	podcastModels "github.com/ziyadrw/faslah/internal/modules/cms/models"
	"time"

	"github.com/ziyadrw/faslah/internal/base"
)

// PodcastListRequest represents the request parameters for listing podcasts
type PodcastListRequest struct {
	base.PaginationRequest
	Sort          string     `query:"sort" form:"sort" validate:"omitempty,oneof=newest oldest popular"`
	PublishedFrom *time.Time `query:"published_from" form:"published_from"`
	PublishedTo   *time.Time `query:"published_to" form:"published_to"`
	Tag           string     `query:"tag" form:"tag"`
}

// PodcastSearchRequest represents the request parameters for searching podcasts
type PodcastSearchRequest struct {
	base.PaginationRequest
	Query string `query:"q" form:"q" validate:"required" message:"مصطلح البحث مطلوب"`
}

// PopularPodcastResponse represents a podcast in the popular podcasts list
type PopularPodcastResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaURL    string `json:"media_url"`
	PlayCount   int    `json:"play_count"`
}

func MapPodcastToDTO(podcast *podcastModels.Podcast) map[string]interface{} {
	return map[string]interface{}{
		"id":            podcast.ID.String(),
		"user_id":       podcast.UserID,
		"title":         podcast.Title,
		"description":   podcast.Description,
		"tags":          podcast.Tags,
		"media_url":     podcast.MediaURL,
		"source_url":    podcast.SourceURL,
		"duration_secs": podcast.DurationSecs,
		"published_at":  podcast.PublishedAt,
		"created_at":    podcast.CreatedAt,
		"updated_at":    podcast.UpdatedAt,
	}
}
