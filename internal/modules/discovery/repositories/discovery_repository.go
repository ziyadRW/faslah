package repositories

import (
	podcastModels "github.com/ziyadrw/faslah/internal/modules/podcast/models"
	"gorm.io/gorm"
	"time"
)

type DiscoveryRepository struct {
	DB *gorm.DB
}

func NewDiscoveryRepository(DB *gorm.DB) *DiscoveryRepository {
	return &DiscoveryRepository{DB: DB}
}

// ListPodcasts retrieves a paginated list of podcasts with optional filters
func (dr *DiscoveryRepository) ListPodcasts(page, perPage int, sort string, publishedFrom, publishedTo *time.Time, tag string) ([]podcastModels.Podcast, int64, error) {
	var podcasts []podcastModels.Podcast
	var totalCount int64

	query := dr.DB.Model(&podcastModels.Podcast{})
	if publishedFrom != nil {
		query = query.Where("published_at >= ?", publishedFrom)
	}

	if publishedTo != nil {
		query = query.Where("published_at <= ?", publishedTo)
	}

	if tag != "" {
		query = query.Where("? = ANY(tags)", tag)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "newest":
		query = query.Order("published_at DESC")
	case "oldest":
		query = query.Order("published_at ASC")
	case "popular":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("published_at DESC")
	}

	offset := (page - 1) * perPage
	query = query.Offset(offset).Limit(perPage)

	if err := query.Find(&podcasts).Error; err != nil {
		return nil, 0, err
	}

	return podcasts, totalCount, nil
}

// SearchPodcasts searches for podcasts by title, description, or tags
func (dr *DiscoveryRepository) SearchPodcasts(query string, page, perPage int) ([]podcastModels.Podcast, int64, error) {
	var podcasts []podcastModels.Podcast
	var totalCount int64

	searchQuery := dr.DB.Model(&podcastModels.Podcast{}).
		Where("title ILIKE ? OR description ILIKE ? OR ? = ANY(tags)",
			"%"+query+"%", "%"+query+"%", query)

	if err := searchQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := searchQuery.Offset(offset).Limit(perPage).Find(&podcasts).Error; err != nil {
		return nil, 0, err
	}

	return podcasts, totalCount, nil
}
