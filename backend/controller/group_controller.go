package controller

import (
	"whatsupp-backend/dto"
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
	claims := util.GetClaims(c)
	req := new(dto.SearchGroupRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_GROUPS, err)
	}

	results, err := g.groupService.HandleLists(ctx, claims, req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LIST_GROUPS, err)
	}

	return util.ResponseOk(c, message.ERR_LIST_GROUPS, results)
}
