package service

import (
	"context"
	"errors"
	"github.com/GlebPoroshin/geochat-auth-service/internal/models"
	"github.com/GlebPoroshin/geochat-auth-service/internal/repository"
	sharedJWT "github.com/GlebPoroshin/geochat-shared/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

type authService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.VerificationRepository
	emailService     EmailService
	jwtSecret        string
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.VerificationRepository,
	emailService EmailService,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		emailService:     emailService,
		jwtSecret:        jwtSecret,
	}
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

func (s *authService) Register(ctx context.Context, login, email, password string) (string, error) {
	if _, err := s.userRepo.FindByLoginOrEmail(ctx, login); err == nil {
		return "", errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		ID:       uuid.New().String(),
		Login:    login,
		Email:    email,
		Password: string(hashedPassword),
		Verified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	code := generateVerificationCode()
	verificationCode := &models.VerificationCode{
		UserID:    user.ID,
		Code:      code,
		Type:      "registration",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.verificationRepo.Create(ctx, verificationCode); err != nil {
		return "", err
	}

	if err := s.emailService.SendVerificationCode(email, code); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *authService) VerifyRegistration(ctx context.Context, userID, code string) error {
	verificationCode, err := s.verificationRepo.FindValid(ctx, userID, code, "registration")
	if err != nil {
		return errors.New("invalid or expired verification code")
	}

	if err := s.userRepo.UpdateVerificationStatus(ctx, userID, true); err != nil {
		return err
	}

	return s.verificationRepo.Delete(ctx, verificationCode)
}

func (s *authService) Login(ctx context.Context, loginOrEmail, password string) (AuthResponse, error) {
	user, err := s.userRepo.FindByLoginOrEmail(ctx, loginOrEmail)
	if err != nil {
		return AuthResponse{}, errors.New("user not found")
	}

	if !user.Verified {
		return AuthResponse{}, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return AuthResponse{}, errors.New("invalid password")
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) InitiatePasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByLoginOrEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	code := generateVerificationCode()
	verificationCode := &models.VerificationCode{
		UserID:    user.ID,
		Code:      code,
		Type:      "password_reset",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.verificationRepo.Create(ctx, verificationCode); err != nil {
		return err
	}

	return s.emailService.SendPasswordResetCode(email, code)
}

func (s *authService) VerifyPasswordResetCode(ctx context.Context, email, code string) error {
	user, err := s.userRepo.FindByLoginOrEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	_, err = s.verificationRepo.FindValid(ctx, user.ID, code, "password_reset")
	if err != nil {
		return errors.New("invalid or expired verification code")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, email, newPassword string) error {
	user, err := s.userRepo.FindByLoginOrEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword))
}

func (s *authService) RefreshToken(_ context.Context, userID, refreshToken string) (AuthResponse, error) {
	claims, err := sharedJWT.Validate(refreshToken)
	if err != nil {
		return AuthResponse{}, errors.New("invalid refresh token")
	}

	if claims.Subject != userID {
		return AuthResponse{}, errors.New("token does not match user")
	}

	accessToken, newRefreshToken, err := s.generateTokens(userID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
