package service

import "context"

type AuthService interface {
	Register(ctx context.Context, login, email, password string) (string, error)
	VerifyRegistration(ctx context.Context, userID, code string) error
	Login(ctx context.Context, loginOrEmail, password string) (AuthResponse, error)
	InitiatePasswordReset(ctx context.Context, email string) error
	VerifyPasswordResetCode(ctx context.Context, email, code string) error
	ResetPassword(ctx context.Context, email, newPassword string) error
	RefreshToken(ctx context.Context, userID, refreshToken string) (AuthResponse, error)
}

type AuthResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type EmailService interface {
	SendVerificationCode(email, code string) error
	SendPasswordResetCode(email, code string) error
}
