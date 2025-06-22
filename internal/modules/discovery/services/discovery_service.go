package discovery

import (
	"github.com/ziyadrw/faslah/internal/base"
	discoveryDTOs "github.com/ziyadrw/faslah/internal/modules/discovery/dtos"
	podcast "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
)

type DiscoveryService struct {
	PodcastRepository *podcast.PodcastRepository
}

func NewDiscoveryService(podcastRepository *podcast.PodcastRepository) *DiscoveryService {
	return &DiscoveryService{PodcastRepository: podcastRepository}
}

// ListPodcasts retrieves a paginated list of podcasts with optional filters
func (ds *DiscoveryService) ListPodcasts(dto discoveryDTOs.PodcastListRequest) base.Response {
	podcasts, totalCount, err := ds.PodcastRepository.ListPodcasts(
		dto.Page,
		dto.PerPage,
		dto.Sort,
		dto.PublishedFrom,
		dto.PublishedTo,
		dto.Tag,
	)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}

	var responseItems []interface{}
	for _, podcast := range podcasts {
		responseItems = append(responseItems, discoveryDTOs.MapPodcastToDTO(&podcast))
	}

	return base.SetPaginatedResponse(responseItems, dto.Page, dto.PerPage, int(totalCount))
}

// SearchPodcasts searches for podcasts by title, description, or tags
func (ds *DiscoveryService) SearchPodcasts(dto discoveryDTOs.PodcastSearchRequest) base.Response {
	podcasts, totalCount, err := ds.PodcastRepository.SearchPodcasts(
		dto.Query,
		dto.Page,
		dto.PerPage,
	)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}

	var responseItems []interface{}
	for _, podcast := range podcasts {
		responseItems = append(responseItems, discoveryDTOs.MapPodcastToDTO(&podcast))
	}

	return base.SetPaginatedResponse(responseItems, dto.Page, dto.PerPage, int(totalCount))
}
