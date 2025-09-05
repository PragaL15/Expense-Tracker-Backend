package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			return fiber.ErrUnauthorized
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}
		if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
			return fiber.ErrUnauthorized
		}
		c.Locals("user_id", claims["sub"])
		return c.Next()
	}
}
