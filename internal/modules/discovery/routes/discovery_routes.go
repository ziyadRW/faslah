package discovery

import (
	"github.com/labstack/echo/v4"
	discoveryHandlers "github.com/ziyadrw/faslah/internal/modules/discovery/handlers"
	discoveryServices "github.com/ziyadrw/faslah/internal/modules/discovery/services"
	podcastRepositories "github.com/ziyadrw/faslah/internal/modules/podcast/repositories"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	podcastRepo := podcastRepositories.NewPodcastRepository(db)
	discoveryService := discoveryServices.NewDiscoveryService(podcastRepo)
	discoveryHandler := discoveryHandlers.NewDiscoveryHandler(discoveryService)

	discoveryGroup := e.Group("/discovery")
	discoveryGroup.GET("", discoveryHandler.ListPodcasts)
	discoveryGroup.GET("/search", discoveryHandler.SearchPodcasts)
}
