package podcast

import (
	"github.com/labstack/echo/v4"
	podcastHandlers "github.com/ziyadrw/faslah/internal/modules/podcast/handlers"
	podcastRepositories "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/podcast/services"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	podcastRepo := podcastRepositories.NewPodcastRepository(db)
	podcastService := podcastServices.NewPodcastService(podcastRepo)
	podcastHandler := podcastHandlers.NewPodcastHandler(podcastService)

	crmGroup := e.Group("/crm")
	crmGroup.POST("create-content", podcastHandler.CreateContent)
}
