package main

import (
	"Go-StandingbookServer/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoServer *echo.Echo

func main() {
	echoServer = echo.New()
	echoServer.Use(middleware.Logger())
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	staticFile()

	echoServer.Logger.Fatal(echoServer.Start(":9527"))
}

func staticFile() {
	echoServer.GET("/ping", routes.PingHandler)
	echoServer.POST("/register", routes.Register)
	echoServer.POST("/login", routes.Login)
	echoServer.GET("/session", routes.GetSession)
	echoServer.POST("/content", routes.Upload)
}
