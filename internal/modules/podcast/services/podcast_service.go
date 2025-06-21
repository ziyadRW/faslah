package podcast

import (
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/podcast/dtos"
	podcast "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
)

type PodcastService struct {
	PodcastRepository *podcast.PodcastRepository
}

func NewPodcastService(podcastRepository *podcast.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

func (ps *PodcastService) CreateContent(dto podcastDTOs.CreateContentRequest) base.Response {
	return base.Response{}
}
