package middleware

import (
	sharedJWT "github.com/GlebPoroshin/geochat-shared/jwt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	claims, err := sharedJWT.Validate(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Добавляем userID в контекст запроса
	c.Locals("userID", claims.Subject)
	return c.Next()
}
