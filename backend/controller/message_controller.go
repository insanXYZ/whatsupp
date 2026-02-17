package controller

import (
	"fmt"
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/dto/message"
	"whatsupp-backend/service"
	"whatsupp-backend/util"

	"github.com/labstack/echo/v5"
)

type MessageController struct {
	messageService *service.MessageService
}

func NewChatController(messageService *service.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

func (mc *MessageController) UpgradeWs(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)

	err := mc.messageService.HandleUpgradeWs(ctx, claims, c.Response(), c.Request())
	if err != nil {
		fmt.Println("error ws:", err.Error())
		return err
	}

	return nil
}

func (mc *MessageController) UploadFileAttachments(c *echo.Context) error {
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

	err = mc.messageService.HandleUploadFileAttachments(ctx, messageId, files)
	if err != nil {
		return util.ResponseErr(c, message.ERR_SEND_FILES, err)
	}

	return util.ResponseOk(c, message.SUCCESS_SEND_FILES, nil)

}

func (mc *MessageController) GetMessages(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)
	req := new(dto.GetMessageRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	messages, err := mc.messageService.HandleGetMessages(ctx, req.ConversationId, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_GET_MESSAGES, err)
	}

	messageDtos := make([]*dto.ItemGetMessagesResponse, len(messages))

	for i, message := range messages {
		messageDtos[i] = converter.MessageEntitytoItemGetMessagesResponseDto(message, claims.Sub)
	}

	response := &dto.GetMessagesResponse{
		ConversationId: req.ConversationId,
		Message:        messageDtos,
	}

	return util.ResponseOk(c, message.SUCCESS_GET_MESSAGES, response)
}
