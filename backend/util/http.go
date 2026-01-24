package util

import (
	"whatsupp-backend/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func GetClaims(c *echo.Context) jwt.MapClaims {
	return c.Get("user").(jwt.MapClaims)
}

func ResponseOk(c *echo.Context, message string, data any, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, 200)
	}

	response := dto.Response{
		Message: message,
		Data:    data,
	}

	return c.JSON(statusCode[0], response)
}

func ResponseErr(c *echo.Context, message string, err any, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, 400)
	}

	response := dto.Response{
		Message: message,
		Error:   err,
	}

	return c.JSON(statusCode[0], response)
}
