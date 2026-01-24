package controller

import (
	"whatsupp-backend/service"
	"whatsupp-backend/util"

	"github.com/labstack/echo/v5"
)

type ChatController struct {
	chatService *service.ChatService
}

func NewChatController(chatService *service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

func (cc *ChatController) UpgradeWs(c *echo.Context) error {
	ctx := c.Request().Context()

	err := cc.chatService.HandleUpgradeWs(ctx, c.Response(), c.Request())
	if err != nil {
		return err
	}

	return util.ResponseOk(c, "upgrade connection success", nil)
}
