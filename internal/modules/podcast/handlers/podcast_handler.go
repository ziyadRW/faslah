package podcast

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/podcast/dtos"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/podcast/services"
)

type PodcastHandler struct {
	PodcastService *podcastServices.PodcastService
}

func NewPodcastHandler(podcastService *podcastServices.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		PodcastService: podcastService,
	}
}

func (ph *PodcastHandler) CreateContent(c echo.Context) error {
	var dto podcastDTOs.CreateContentRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := ph.PodcastService.CreateContent(dto)
	return c.JSON(response.HTTPStatus, response)
}
