package routes

import (
	"github.com/labstack/echo/v4"
	cms "github.com/ziyadrw/faslah/internal/modules/cms/routes"
	discovery "github.com/ziyadrw/faslah/internal/modules/discovery/routes"
	user "github.com/ziyadrw/faslah/internal/modules/user/routes"
	"gorm.io/gorm"
)

func RegisterAllRoutes(e *echo.Echo, db *gorm.DB) {
	user.RegisterRoutes(e, db)
	cms.RegisterRoutes(e, db)
	discovery.RegisterRoutes(e, db)
	RegisterSwaggerRoutes(e)
}
