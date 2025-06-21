package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/config"
	"github.com/ziyadrw/faslah/internal/base"
	"github.com/ziyadrw/faslah/internal/middlewares"
	"github.com/ziyadrw/faslah/internal/migrations"
	"github.com/ziyadrw/faslah/internal/routes"
)

func main() {
	e := echo.New()

	base.RegisterValidator(e)
	middlewares.RegisterAllGlobalMiddlewares(e)

	config.Connect()
	db := config.GetDB()
	migrations.Migrate()

	routes.RegisterAllRoutes(e, db)

	startServer(e)
}
