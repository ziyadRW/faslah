package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	usersRoutes := e.Group("/users")
	usersRoutes.GET("/hello-world", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

}
