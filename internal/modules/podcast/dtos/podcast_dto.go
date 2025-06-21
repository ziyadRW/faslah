package podcast

import (
	"time"
)

type CreateContentRequest struct {
	MediaURL  string `json:"media_url,omitempty"`  // For uploaded files (e.g. https://media.faslah.com/...)
	SourceURL string `json:"source_url,omitempty"` // For YouTube or external links to be fetched

	Title        string     `json:"title,omitempty"`         // Optional: user can override extracted title
	Description  string     `json:"description,omitempty"`   // Optional: user can override description
	Tags         []string   `json:"tags,omitempty"`          // Optional: user-supplied tags
	PublishedAt  *time.Time `json:"published_at,omitempty"`  // Optional: for scheduling or back-dated content
	DurationSecs int        `json:"duration_secs,omitempty"` // Optional: can be auto-filled from file/YT
}
