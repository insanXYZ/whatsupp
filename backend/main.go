package main

import (
	"fmt"
	"os"
	"whatsupp-backend/config"
	"whatsupp-backend/controller"
	"whatsupp-backend/repository"
	"whatsupp-backend/service"
	"whatsupp-backend/websocket"

	"github.com/joho/godotenv"
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
	app := config.NewEcho()

	hub := websocket.NewHub()
	go hub.Run()

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

	SetRoute(&RouteConfig{
		userController: userController,
		chatController: chatController,
		app:            app,
	})

	err = app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		app.Logger.Error("failed to start server", "error", err)
	}

	app.Logger.Info(fmt.Sprintf("app start on %s", os.Getenv("APP_PORT")))

}
