package controller

import (
	"net/http"
	"time"
	"whatsupp-backend/dto"
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
		return err
	}

	user, err := u.userService.HandleLogin(ctx, req)
	if err != nil {
		return err
	}

	accToken, err := util.GenerateJWT(jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})

	if err != nil {
		return err
	}

	cookies := new(http.Cookie)
	cookies.Name = "X-ACC-TOKEN"
	cookies.Value = accToken
	cookies.HttpOnly = true
	cookies.Path = "/"
	cookies.Secure = true

	c.SetCookie(cookies)
	return c.JSON(200, "sukses login")
}

func (u *UserController) Register(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.RegisterRequest)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	return u.userService.HandleRegister(ctx, req)
}

func (u *UserController) UpdateUser(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(dto.UpdateUserRequest)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	claims := util.GetClaims(c)

	err = u.userService.HandleUpdateUser(ctx, req, claims)
	if err != nil {
		return err
	}

	return c.JSON(200, "sukses update user")
}
