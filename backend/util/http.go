package util

import (
	"net/http"
	"whatsupp-backend/dto"

	"github.com/labstack/echo/v5"
)

func ResponseOk(c *echo.Context, message string, data any, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusOK)
	}

	response := dto.Response{
		Message: message,
		Data:    data,
	}

	return c.JSON(statusCode[0], response)
}

func ResponseErr(c *echo.Context, message string, err any, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusBadRequest)
	}

	response := dto.Response{
		Message: message,
		Error:   err,
	}

	return c.JSON(statusCode[0], response)
}

func ResponseErrInternal(c *echo.Context, err any) error {
	response := dto.Response{
		Message: "Internal Server Error",
		Error:   err,
	}

	return c.JSON(http.StatusInternalServerError, response)
}
