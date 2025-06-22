package podcast

import (
	"github.com/google/uuid"
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

// GetContent returns a podcast
func (ps *PodcastService) GetContent(id string) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)

	return base.SetData(response)
}

// UpdateContent updates a podcast
func (ps *PodcastService) UpdateContent(id string, dto podcastDTOs.UpdateContentRequest) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	if dto.Title != "" {
		podcast.Title = dto.Title
	}
	if dto.Description != "" {
		podcast.Description = dto.Description
	}
	if dto.Tags != nil {
		podcast.Tags = dto.Tags
	}
	if dto.MediaURL != "" {
		podcast.MediaURL = dto.MediaURL
	}
	if dto.SourceURL != "" {
		podcast.SourceURL = dto.SourceURL
	}
	if dto.DurationSecs != 0 {
		podcast.DurationSecs = dto.DurationSecs
	}
	if dto.PublishedAt != nil {
		podcast.PublishedAt = dto.PublishedAt
	}

	if err := ps.PodcastRepository.UpdatePodcast(podcast); err != nil {
		return base.SetErrorMessage("فشل في تحديث البودكاست", err.Error())
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)

	return base.SetData(response, "تم تحديث البودكاست بنجاح")
}

func (ps *PodcastService) DeleteContent(id string) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	if err := ps.PodcastRepository.DeletePodcast(podcastID); err != nil {
		return base.SetErrorMessage("فشل في حذف البودكاست", err.Error())
	}

	return base.SetSuccessMessage("تم حذف البودكاست بنجاح")
}
