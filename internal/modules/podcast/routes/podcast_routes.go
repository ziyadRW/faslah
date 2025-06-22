package podcast

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/middlewares"
	podcastHandlers "github.com/ziyadrw/faslah/internal/modules/podcast/handlers"
	podcastRepositories "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/podcast/services"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	podcastRepo := podcastRepositories.NewPodcastRepository(db)
	podcastService := podcastServices.NewPodcastService(podcastRepo)
	podcastHandler := podcastHandlers.NewPodcastHandler(podcastService)

	cmsGroup := e.Group("/cms", middlewares.RoleMiddleware(db))
	cmsGroup.POST("/create-content", podcastHandler.CreateContent)

	cmsGroup.GET("/retreive-content/:id", podcastHandler.GetContent)
	cmsGroup.PUT("/update-content/:id", podcastHandler.UpdateContent)
	cmsGroup.DELETE("/delete-content/:id", podcastHandler.DeleteContent)
}
