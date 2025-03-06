package router

import (
	"github.com/GlebPoroshin/geochat-auth-service/internal/api/handlers"
	"github.com/GlebPoroshin/geochat-auth-service/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Не требуют токена в хедере
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify-registration", authHandler.VerifyRegistration)
	auth.Post("/login", authHandler.Login)
	auth.Post("/password-reset", authHandler.InitiatePasswordReset)
	auth.Post("/verify-reset-code", authHandler.VerifyPasswordResetCode)
	auth.Post("/reset-password", authHandler.ResetPassword)
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
		Title:       "Auth Service API Documentation",
	}))

	// Требуют токен в хедере
	protected := auth.Use(middleware.RequireAuth)
	protected.Get("/refresh", authHandler.RefreshToken)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
		})
	})
}
