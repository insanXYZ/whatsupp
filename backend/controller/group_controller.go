package controller

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/dto/message"
	"whatsupp-backend/service"
	"whatsupp-backend/util"

	"github.com/labstack/echo/v5"
)

type GroupController struct {
	groupService *service.GroupService
}

func NewGroupController(groupService *service.GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

func (g *GroupController) Lists(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.ListGroupRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_USERS, err)
	}

	users, err := g.groupService.HandleLists(ctx, req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_USERS, err)
	}

	usersDto := make([]dto.User, len(users))

	for i, user := range users {
		usersDto[i] = *converter.UserEntityToDto(&user)
	}

	return util.ResponseOk(c, message.SUCCESS_LIST_USERS, usersDto)
}
