package util

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type Claims struct {
	Sub   int
	Email string
}

func GenerateJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func GetClaims(c *echo.Context) *Claims {
	mapClaims := c.Get("user").(jwt.MapClaims)
	return &Claims{
		Sub:   int(mapClaims["sub"].(float64)),
		Email: mapClaims["email"].(string),
	}
}
