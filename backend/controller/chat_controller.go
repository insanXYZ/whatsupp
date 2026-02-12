package controller

import (
	"fmt"
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/message"
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
	claims := util.GetClaims(c)

	err := cc.chatService.HandleUpgradeWs(ctx, claims, c.Response(), c.Request())
	if err != nil {
		fmt.Println("error ws:", err.Error())
		return err
	}

	return nil
}

func (cc *ChatController) UploadFileAttachments(c *echo.Context) error {
	ctx := c.Request().Context()
	messageId := c.FormValue("message_id")
	form, err := c.MultipartForm()
	if err != nil {
		return util.ResponseErr(c, message.ERR_RETRIEVE_FILES, err)
	}

	files, ok := form.File[dto.MULTIPART_FORM_NAME]
	if !ok {
		return util.ResponseErr(c, message.ERR_RETRIEVE_FILES, nil)
	}

	err = cc.chatService.HandleUploadFileAttachments(ctx, messageId, files)
	if err != nil {
		return util.ResponseErr(c, message.ERR_SEND_FILES, err)
	}

	return util.ResponseOk(c, message.SUCCESS_SEND_FILES, nil)

}
