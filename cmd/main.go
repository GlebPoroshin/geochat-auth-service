// @title Auth Service API
// @version 1.0
// @description API сервиса авторизации
// @host localhost:8081
// @BasePath /auth
// @schemes http
package main

import (
	"github.com/GlebPoroshin/geochat-auth-service/docs"
	"github.com/GlebPoroshin/geochat-auth-service/internal/api/handlers"
	"github.com/GlebPoroshin/geochat-auth-service/internal/api/router"
	"github.com/GlebPoroshin/geochat-auth-service/internal/config"
	"github.com/GlebPoroshin/geochat-auth-service/internal/models"
	"github.com/GlebPoroshin/geochat-auth-service/internal/repository"
	"github.com/GlebPoroshin/geochat-auth-service/internal/service"
	sharedJWT "github.com/GlebPoroshin/geochat-shared/jwt"
	"github.com/gofiber/fiber/v2"

	_ "github.com/gofiber/swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	sharedJWT.Init(cfg.JWTSecret)

	// Инициализация Swagger
	docs.SwaggerInfo.Title = "Auth Service API"
	docs.SwaggerInfo.Description = "API сервиса авторизации"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.ServerAddress
	docs.SwaggerInfo.BasePath = "/auth"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)

	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(userRepo, verificationRepo, emailService, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New()
	router.SetupRoutes(app, authHandler)

	log.Fatal(app.Listen(cfg.ServerAddress))
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Миграция моделей
	if err := db.AutoMigrate(&models.User{}, &models.VerificationCode{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	log.Println("Database migration completed successfully")
	return db, nil
}
