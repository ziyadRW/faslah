package podcast

import "gorm.io/gorm"

type PodcastRepository struct {
	DB *gorm.DB
}

func NewPodcastRepository(DB *gorm.DB) *PodcastRepository {
	return &PodcastRepository{DB: DB}
}
