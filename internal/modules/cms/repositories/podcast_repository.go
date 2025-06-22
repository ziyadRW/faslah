package cms

import (
	"errors"
	"github.com/google/uuid"
	podcastModels "github.com/ziyadrw/faslah/internal/modules/cms/models"
	"gorm.io/gorm"
)

type PodcastRepository struct {
	DB *gorm.DB
}

func NewPodcastRepository(DB *gorm.DB) *PodcastRepository {
	return &PodcastRepository{DB: DB}
}

// CreatePodcast creates a new podcast in the database
func (pr *PodcastRepository) CreatePodcast(podcast *podcastModels.Podcast) error {
	return pr.DB.Create(podcast).Error
}

// GetPodcastByID retrieves a podcast by ID
func (pr *PodcastRepository) GetPodcastByID(id uuid.UUID) (*podcastModels.Podcast, error) {
	var podcast podcastModels.Podcast

	result := pr.DB.Where("id = ?", id).First(&podcast)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &podcast, nil
}

// UpdatePodcast updates a podcast in the database
func (pr *PodcastRepository) UpdatePodcast(podcast *podcastModels.Podcast) error {
	return pr.DB.Save(podcast).Error
}

// DeletePodcast soft-deletes a podcast from the database
func (pr *PodcastRepository) DeletePodcast(id uuid.UUID) error {
	return pr.DB.Delete(&podcastModels.Podcast{}, id).Error
}

// GetPodcastsByUserID retrieves all podcasts created by a specific user
func (pr *PodcastRepository) GetPodcastsByUserID(userID string) ([]podcastModels.Podcast, error) {
	var podcasts []podcastModels.Podcast

	result := pr.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&podcasts)
	if result.Error != nil {
		return nil, result.Error
	}

	return podcasts, nil
}
