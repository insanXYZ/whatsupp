package main

import (
	"whatsupp-backend/controller"
	"whatsupp-backend/middleware"

	"github.com/labstack/echo/v5"
)

type RouteConfig struct {
	userController *controller.UserController
	chatController *controller.ChatController
	app            *echo.Echo
}

func SetRoute(cfg *RouteConfig) {
	middleware.InitMiddleware()

	app := cfg.app

	api := app.Group("/api")
	api.POST("/login", cfg.userController.Login)
	api.POST("/register", cfg.userController.Register)

	user := api.Group("/me", middleware.HasJWT)
	user.PUT("", cfg.userController.UpdateMe)
	user.GET("", cfg.userController.Me)

	api.GET("/ws", cfg.chatController.UpgradeWs)
}
