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

func (ph *PodcastHandler) GetContent(c echo.Context) error {
	id := c.Param("id")
	response := ph.PodcastService.GetContent(id)
	return c.JSON(response.HTTPStatus, response)
}

func (ph *PodcastHandler) UpdateContent(c echo.Context) error {
	id := c.Param("id")

	var dto podcastDTOs.UpdateContentRequest
	if res, ok := base.BindAndValidate(c, &dto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := ph.PodcastService.UpdateContent(id, dto)
	return c.JSON(response.HTTPStatus, response)
}

func (ph *PodcastHandler) DeleteContent(c echo.Context) error {
	id := c.Param("id")
	response := ph.PodcastService.DeleteContent(id)
	return c.JSON(response.HTTPStatus, response)
}
