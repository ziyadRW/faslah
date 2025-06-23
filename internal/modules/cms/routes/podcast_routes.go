package cms

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/middlewares"
	podcastHandlers "github.com/ziyadrw/faslah/internal/modules/cms/handlers"
	podcastRepositories "github.com/ziyadrw/faslah/internal/modules/cms/repositories"
	podcastServices "github.com/ziyadrw/faslah/internal/modules/cms/services"
	user "github.com/ziyadrw/faslah/internal/modules/user/enums"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	podcastRepo := podcastRepositories.NewPodcastRepository(db)
	podcastService := podcastServices.NewPodcastService(podcastRepo)
	podcastHandler := podcastHandlers.NewPodcastHandler(podcastService)

	cmsGroup := e.Group("/cms", middlewares.RoleMiddleware(db, user.TypeCreator, user.TypeAdmin))
	cmsGroup.POST("/create-content", podcastHandler.CreateContent)
	cmsGroup.POST("/my-content", podcastHandler.MyContent)
	cmsGroup.GET("/retreive-content/:id", podcastHandler.GetContent)
	cmsGroup.PUT("/update-content/:id", podcastHandler.UpdateContent)
	cmsGroup.DELETE("/delete-content/:id", podcastHandler.DeleteContent)

	cmsGroupNoAuth := e.Group("/cms")
	cmsGroupNoAuth.POST("/fetch-youtube-content", podcastHandler.FetchFromYouTube)
}
