package discovery

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	discoveryHandlers "github.com/ziyadrw/faslah/internal/modules/discovery/handlers"
	discoveryRepositories "github.com/ziyadrw/faslah/internal/modules/discovery/repositories"
	discoveryServices "github.com/ziyadrw/faslah/internal/modules/discovery/services"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	discoveryRepo := discoveryRepositories.NewDiscoveryRepository(db)
	discoveryService := discoveryServices.NewDiscoveryService(discoveryRepo)
	discoveryHandler := discoveryHandlers.NewDiscoveryHandler(discoveryService)

	discoveryGroup := e.Group("/discovery")
	discoveryGroup.GET("", discoveryHandler.ListPodcasts)
	discoveryGroup.GET("/search", discoveryHandler.SearchPodcasts)
}
