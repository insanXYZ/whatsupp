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
	clientStorage, err := config.NewSupabaseStorageClient()
	if err != nil {
		panic(err.Error())
	}

	// int hub ws

	hub := websocket.NewHub()
	go hub.Run()

	// init repository
	memberRepository := repository.NewMemberRepository(gorm)
	conversationRepository := repository.NewConversationRepository(gorm)
	messageAttachmentRepository := repository.NewMessageAttachmentRepository(gorm)
	messageRepository := repository.NewMessageRepository(gorm)
	userRepository := repository.NewUserRepository(gorm)

	// init service

	messageService := service.NewMessageService(validator, conversationRepository, memberRepository, messageRepository, messageAttachmentRepository, userRepository, hub, clientStorage)
	userService := service.NewUserService(validator, userRepository)
	conversationService := service.NewConversationService(validator, memberRepository, conversationRepository)

	// init controller

	messageController := controller.NewChatController(messageService)
	userController := controller.NewUserController(userService)
	conversationController := controller.NewConversationController(conversationService)

	SetRoute(&RouteConfig{
		userController:         userController,
		messageController:      messageController,
		conversationController: conversationController,
		app:                    app,
	})

	err = app.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		app.Logger.Error("failed to start server", "error", err)
	}

	app.Logger.Info("app stopped")

}
