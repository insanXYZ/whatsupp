package middleware

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

var (
	HasJWT echo.MiddlewareFunc
)

func InitMiddleware() {
	HasJWT = echojwt.WithConfig(echojwt.Config{
		SuccessHandler: func(c *echo.Context) error {
			c.Set("user", c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims))
			return nil
		},
		SigningKey:    []byte(os.Getenv("JWT_SECRET_KEY")),
		SigningMethod: echojwt.AlgorithmHS256,
	})
}
