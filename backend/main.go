package main

import (
	"whatsupp-backend/config"
	"whatsupp-backend/controller"
	"whatsupp-backend/repository"
	"whatsupp-backend/service"
	"whatsupp-backend/websocket"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	// init config
	gorm, err := config.NewGorm()
	if err != nil {
		panic(err.Error())
	}

	validator := config.NewValidator()

	hub := websocket.NewHub()
	hub.Run()

	// init repository
	groupMemberRepository := repository.NewGroupMessageRepository(gorm)
	groupRepository := repository.NewGroupRepository(gorm)
	messageAttachmentRepository := repository.NewMessageAttachmentRepository(gorm)
	messageRepository := repository.NewMessageRepository(gorm)
	userRepository := repository.NewUserRepository(gorm)

	// init service

	chatService := service.NewChatService(groupRepository, groupMemberRepository, messageRepository, messageAttachmentRepository, validator, hub)
	userService := service.NewUserService(validator, userRepository)

	// init controller

	chatController := controller.NewChatController(chatService)
	userController := controller.NewUserController(userService)

	app := echo.New()
	err = app.Start(":3030")
	if err != nil {
		app.Logger.Error("failed to start server", "error", err)
	}

}
