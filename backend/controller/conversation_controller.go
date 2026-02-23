package controller

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/dto/message"
	"whatsupp-backend/service"
	"whatsupp-backend/util"

	"github.com/labstack/echo/v5"
)

type ConversationController struct {
	conversationService *service.ConversationService
}

func NewConversationController(conversationService *service.ConversationService) *ConversationController {
	return &ConversationController{
		conversationService: conversationService,
	}
}

func (cc *ConversationController) Lists(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)
	req := new(dto.SearchConversationRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_CONVERSATION, err)
	}

	results, err := cc.conversationService.HandleFindConversations(ctx, claims, req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_CONVERSATION, err)
	}

	return util.ResponseOk(c, message.SUCCESS_LIST_CONVERSATION, results)
}

func (cc *ConversationController) LoadRecentConversations(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)

	recentConversations, err := cc.conversationService.HandleLoadRecentConversations(ctx, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_RECENT_CONVERSATION, err)
	}

	return util.ResponseOk(c, message.SUCCESS_LIST_RECENT_CONVERSATION, recentConversations)
}

func (cc *ConversationController) CreateGroupConversation(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)
	req := new(dto.CreateGroupConversationRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	file, err := c.FormFile("image")
	if err == nil {
		req.Image = file
	}

	newConversation, err := cc.conversationService.HandleCreateGroupConversation(ctx, req, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_CREATE_GROUP_CONVERSATION, err)
	}

	return util.ResponseOk(c, message.SUCCESS_CREATE_GROUP_CONVERSATION, converter.ConversationEntityToDto(newConversation))

}

func (cc *ConversationController) JoinGroupConversation(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)
	req := new(dto.JoinGroupConversationRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	isJoin, err := cc.conversationService.HandleJoinGroupConversation(ctx, req, claims)
	if err != nil {
		errMessage := message.ERR_JOIN_GROUP
		if !isJoin {
			errMessage = message.ERR_LEAVE_GROUP
		}

		return util.ResponseErr(c, errMessage, err)
	}

	successMessage := message.SUCCESS_JOIN_GROUP
	if !isJoin {
		successMessage = message.SUCCESS_LEAVE_GROUP
	}

	return util.ResponseOk(c, successMessage, nil)

}

func (cc *ConversationController) ListMembersConversation(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.ListMembersConversationRequest)
	claims := util.GetClaims(c)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	members, err := cc.conversationService.HandleListMembersConversation(ctx, req, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_MEMBERS_CONVERSATION, err)
	}

	response := converter.MemberEntitiesToDto(members)

	return util.ResponseOk(c, message.SUCCESS_LIST_MEMBERS_CONVERSATION, response)

}
