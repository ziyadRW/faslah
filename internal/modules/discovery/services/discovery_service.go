package discovery

import (
	"github.com/ziyadrw/faslah/internal/base"
	discoveryDTOs "github.com/ziyadrw/faslah/internal/modules/discovery/dtos"
	discoveryRepositories "github.com/ziyadrw/faslah/internal/modules/discovery/repositories"
)

type DiscoveryService struct {
	DiscoveryRepository *discoveryRepositories.DiscoveryRepository
}

func NewDiscoveryService(discoveryRepository *discoveryRepositories.DiscoveryRepository) *DiscoveryService {
	return &DiscoveryService{DiscoveryRepository: discoveryRepository}
}

// ListPodcasts retrieves a paginated list of podcasts with optional filters
func (ds *DiscoveryService) ListPodcasts(dto discoveryDTOs.PodcastListRequest) base.Response {
	podcasts, totalCount, err := ds.DiscoveryRepository.ListPodcasts(
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
	podcasts, totalCount, err := ds.DiscoveryRepository.SearchPodcasts(
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

// GetPopularPodcasts retrieves the top 10 podcasts by play count in the last 24 hours
func (ds *DiscoveryService) GetPopularPodcasts() base.Response {
	popularPodcasts, err := ds.DiscoveryRepository.GetPopularPodcasts()
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}

	var response []discoveryDTOs.PopularPodcastResponse
	for _, podcast := range popularPodcasts {
		response = append(response, discoveryDTOs.PopularPodcastResponse{
			ID:          podcast["id"].(string),
			Title:       podcast["title"].(string),
			Description: podcast["description"].(string),
			MediaURL:    podcast["media_url"].(string),
			PlayCount:   int(podcast["play_count"].(int64)),
		})
	}

	return base.SetData(response)
}
