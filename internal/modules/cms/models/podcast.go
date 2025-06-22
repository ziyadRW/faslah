package cms

import (
	"github.com/ziyadrw/faslah/internal/base"
	user "github.com/ziyadrw/faslah/internal/modules/user/models"
	"time"
)

type Podcast struct {
	base.Model

	UserID       string     `gorm:"type:uuid;not null;index" json:"user_id"`
	Title        string     `gorm:"type:text;not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Tags         []string   `gorm:"type:text[]" json:"tags"`
	MediaURL     string     `gorm:"type:text;not null" json:"media_url"`
	SourceURL    string     `gorm:"type:text" json:"source_url"`
	DurationSecs int        `json:"duration_secs"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`

	User *user.User `gorm:"foreignKey:UserID" json:"-"`
}
