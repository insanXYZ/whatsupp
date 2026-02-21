package controller

import (
	"time"
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/dto/message"
	"whatsupp-backend/service"
	"whatsupp-backend/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) Login(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.LoginRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	user, err := u.userService.HandleLogin(ctx, req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_LOGIN, err)
	}

	accToken, err := util.GenerateJWT(jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	if err != nil {
		return util.ResponseErrInternal(c, err)
	}

	cookie := util.CreateCookie(dto.COOKIES_ACC_TOKEN_NAME, accToken)

	c.SetCookie(cookie)
	return util.ResponseOk(c, message.SUCCESS_LOGIN, nil)
}

func (u *UserController) Register(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.RegisterRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	err = u.userService.HandleRegister(ctx, req)

	if err != nil {
		return util.ResponseErr(c, message.ERR_REGISTER, err)
	}

	return util.ResponseOk(c, message.SUCCESS_REGISTER, nil)
}

func (u *UserController) UpdateMe(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.UpdateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return util.ResponseErr(c, message.ERR_BIND_REQ, err)
	}

	claims := util.GetClaims(c)

	user, err := u.userService.HandleUpdateUser(ctx, req, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_UPDATE_ME, err)
	}

	return util.ResponseOk(c, message.SUCCESS_UPDATE_ME, converter.UserEntityToDto(user))
}

func (u *UserController) Me(c *echo.Context) error {
	ctx := c.Request().Context()
	claims := util.GetClaims(c)
	user, err := u.userService.HandleMe(ctx, claims)
	if err != nil {
		return util.ResponseErr(c, message.ERR_GET_ME, err)
	}

	return util.ResponseOk(c, message.SUCCESS_GET_ME, converter.UserEntityToDto(user))
}

func (u *UserController) Logout(c *echo.Context) error {

	cookie := util.DeleteCookie(dto.COOKIES_ACC_TOKEN_NAME)
	c.SetCookie(cookie)

	return util.ResponseOk(c, message.SUCCESS_LOGOUT, nil)

}
