package main

import (
	"whatsupp-backend/controller"
	"whatsupp-backend/middleware"

	"github.com/labstack/echo/v5"
)

type RouteConfig struct {
	userController    *controller.UserController
	messageController *controller.MessageController
	groupController   *controller.GroupController
	app               *echo.Echo
}

func SetRoute(cfg *RouteConfig) {
	middleware.InitMiddleware()

	app := cfg.app

	api := app.Group("/api")
	api.POST("/login", cfg.userController.Login)
	api.POST("/register", cfg.userController.Register)

	hasJwtRoute := api.Group("", middleware.HasJWT)

	hasJwtRoute.GET("/ws", cfg.messageController.UpgradeWs)

	user := hasJwtRoute.Group("/me", middleware.HasJWT)
	user.PUT("", cfg.userController.UpdateMe)
	user.GET("", cfg.userController.Me)
	user.GET("/logout", cfg.userController.Logout)

	group := hasJwtRoute.Group("/groups")
	group.GET("/recent", cfg.groupController.LoadRecentGroups)
	group.GET("", cfg.groupController.Lists)

	chat := hasJwtRoute.Group("/messages")
	chat.POST("/attachments", cfg.messageController.UploadFileAttachments)
	chat.GET("/:groupId", cfg.messageController.GetMessages)
}
