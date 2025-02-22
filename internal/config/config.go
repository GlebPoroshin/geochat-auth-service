package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	JWTSecret    string
	SMTPUsername string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     int
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret:    os.Getenv("JWT_SECRET"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
	}, nil
}
