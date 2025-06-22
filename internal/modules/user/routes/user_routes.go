package user

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/internal/middlewares"
	userHandlers "github.com/ziyadrw/faslah/internal/modules/user/handlers"
	userRepositories "github.com/ziyadrw/faslah/internal/modules/user/repositories"
	userServices "github.com/ziyadrw/faslah/internal/modules/user/services"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepo := userRepositories.NewUserRepository(db)
	userService := userServices.NewUserService(userRepo)
	userHandler := userHandlers.NewUserHandler(userService)

	usersRoutes := e.Group("/users")
	usersRoutes.POST("/signup", userHandler.Signup)
	usersRoutes.POST("/login", userHandler.Login)

	authRoutes := usersRoutes.Group("/me", middlewares.RoleMiddleware(db))
	authRoutes.GET("/profile", userHandler.GetProfile)
	authRoutes.GET("/history", userHandler.GetWatchHistory)

	podcastRoutes := e.Group("/podcasts", middlewares.RoleMiddleware(db))
	podcastRoutes.POST("/:id/track-play", userHandler.TrackPlay)
}
