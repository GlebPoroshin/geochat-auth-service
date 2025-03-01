package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
	JWTSecret     string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		JWTSecret:     os.Getenv("JWT_SECRET"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.DBHost == "" || c.DBUser == "" || c.DBPassword == "" || c.DBName == "" {
		return fmt.Errorf("database configuration is incomplete")
	}
	if c.SMTPHost == "" || c.SMTPUsername == "" || c.SMTPPassword == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}
	return nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
