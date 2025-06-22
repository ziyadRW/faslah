package podcast

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/podcast/dtos"
	podcastModels "github.com/ziyadrw/faslah/internal/modules/podcast/models"
	podcast "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
)

type PodcastService struct {
	PodcastRepository *podcast.PodcastRepository
}

func NewPodcastService(podcastRepository *podcast.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

// CreateContent creates a new podcast from either a YouTube URL or a direct MP4 upload
func (ps *PodcastService) CreateContent(dto podcastDTOs.CreatePodcastRequest, userID string) base.Response {
	if (dto.SourceURL != "" && dto.File != nil) || (dto.SourceURL == "" && dto.File == nil) {
		return base.SetErrorMessage(
			"خطأ في البيانات المدخلة",
			"يجب توفير إما رابط يوتيوب أو ملف MP4، وليس كلاهما أو لا شيء",
		)
	}

	if dto.File != nil && dto.Title == "" {
		return base.SetErrorMessage(
			"خطأ في البيانات المدخلة",
			"يجب توفير العنوان عند رفع ملف MP4",
		)
	}

	var mediaURL string
	var title string
	var description string
	var tags []string
	var durationSecs int
	var sourceURL string

	if dto.SourceURL != "" {
		sourceURL = dto.SourceURL

		youtubeResponse := ps.FetchYouTube(dto.SourceURL)
		if youtubeResponse.HTTPStatus != http.StatusOK {
			return youtubeResponse
		}

		youtubeData, ok := youtubeResponse.Data.(podcastDTOs.YouTubeMetadataResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في استخراج بيانات الفيديو")
		}

		uniqueFilename := fmt.Sprintf("%s.mp4", uuid.New().String())
		uploadResponse := ps.UploadMedia(youtubeData.VideoFile, uniqueFilename)
		if uploadResponse.HTTPStatus != http.StatusOK {
			return uploadResponse
		}

		uploadData, ok := uploadResponse.Data.(podcastDTOs.MediaUploadResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في رفع الفيديو")
		}

		mediaURL = uploadData.MediaURL

		if dto.Title == "" {
			title = youtubeData.Title
		} else {
			title = dto.Title
		}

		if dto.Description == "" {
			description = youtubeData.Description
		} else {
			description = dto.Description
		}

		if dto.Tags == nil || len(dto.Tags) == 0 {
		} else {
			tags = dto.Tags
		}

		durationSecs = youtubeData.DurationSecs
	} else {
		src, err := dto.File.Open()
		if err != nil {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في فتح الملف")
		}
		defer src.Close()

		fileContent := make([]byte, dto.File.Size)
		if _, err = src.Read(fileContent); err != nil {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في قراءة محتوى الملف")
		}

		uploadResponse := ps.UploadMedia(fileContent, dto.File.Filename)
		if uploadResponse.HTTPStatus != http.StatusOK {
			return uploadResponse
		}

		uploadData, ok := uploadResponse.Data.(podcastDTOs.MediaUploadResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في رفع الملف")
		}

		mediaURL = uploadData.MediaURL
		title = dto.Title
		description = dto.Description
		tags = dto.Tags

		duration, err := ps.GetVideoDuration(fileContent)
		if err != nil {
			log.Printf("Failed to extract video duration: %v", err)
			durationSecs = 0
		} else {
			durationSecs = duration
		}
	}

	podcast := &podcastModels.Podcast{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Tags:         tags,
		MediaURL:     mediaURL,
		SourceURL:    sourceURL,
		DurationSecs: durationSecs,
		PublishedAt:  dto.PublishedAt,
	}

	if err := ps.PodcastRepository.CreatePodcast(podcast); err != nil {
		return base.SetErrorMessage("فشل في إنشاء البودكاست", err.Error())
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)
	return base.SetData(response, "تم إنشاء البودكاست بنجاح")
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
