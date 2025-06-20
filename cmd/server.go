package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ziyadrw/faslah/config"
	"log"
)

func startServer(e *echo.Echo, setPort ...string) {
	port := ":8080"
	if len(setPort) > 0 {
		port = setPort[0]
	}
	log.Println("Server running on http://" + config.GetEnv("PUBLIC_IP", "127.0.0.1") + port)
	e.Logger.Fatal(e.Start("0.0.0.0" + port))
}
