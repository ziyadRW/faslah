package routes

import (
	"github.com/labstack/echo/v4"
	podcast "github.com/ziyadrw/faslah/internal/modules/podcast/routes"
	user "github.com/ziyadrw/faslah/internal/modules/user/routes"
	"gorm.io/gorm"
)

func RegisterAllRoutes(e *echo.Echo, db *gorm.DB) {
	user.RegisterRoutes(e, db)
	podcast.RegisterRoutes(e, db)
	RegisterSwaggerRoutes(e)
}
