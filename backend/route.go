package main

import (
	"whatsupp-backend/controller"
	"whatsupp-backend/middleware"

	"github.com/labstack/echo/v5"
)

type RouteConfig struct {
	userController  *controller.UserController
	chatController  *controller.ChatController
	groupController *controller.GroupController
	app             *echo.Echo
}

func SetRoute(cfg *RouteConfig) {
	middleware.InitMiddleware()

	app := cfg.app

	api := app.Group("/api")
	api.POST("/login", cfg.userController.Login)
	api.POST("/register", cfg.userController.Register)

	hasJwtRoute := api.Group("", middleware.HasJWT)

	user := hasJwtRoute.Group("/me", middleware.HasJWT)
	user.PUT("", cfg.userController.UpdateMe)
	user.GET("", cfg.userController.Me)

	group := hasJwtRoute.Group("/groups")
	group.GET("", cfg.groupController.Lists)

	chat := hasJwtRoute.Group("/chats")
	chat.GET("/ws", cfg.chatController.UpgradeWs)
	chat.POST("/attachments", cfg.chatController.UploadFileAttachments)
}
