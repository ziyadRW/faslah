package user

import (
	"github.com/ziyadrw/faslah/internal/base"
	"time"
)

type WatchHistory struct {
	base.Model

	UserID         string     `gorm:"type:uuid;not null;index" json:"user_id"`
	PodcastID      string     `gorm:"type:uuid;not null;index" json:"podcast_id"`
	PlaybackSecond int        `gorm:"default:0" json:"playback_second"`
	LastPlayedAt   *time.Time `json:"last_played_at"`
}
