package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"geochat-auth-service/internal/models"
)

type VerificationRepository interface {
	Create(ctx context.Context, code *models.VerificationCode) error
	FindValid(ctx context.Context, userID, code, codeType string) (*models.VerificationCode, error)
	Delete(ctx context.Context, code *models.VerificationCode) error
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) Create(ctx context.Context, code *models.VerificationCode) error {
	return r.db.Create(code).Error
}

func (r *verificationRepository) FindValid(ctx context.Context, userID, code, codeType string) (*models.VerificationCode, error) {
	var verificationCode models.VerificationCode
	if err := r.db.Where(
		"user_id = ? AND code = ? AND type = ? AND expires_at > ?",
		userId, code, codeType, time.Now(),
	).First(&verificationCode).Error; err != nil {
		return nil, err
	}
	return &verificationCode, nil
}

func (r *verificationRepository) Delete(ctx context.Context, code *models.VerificationCode) error {
	return r.db.Delete(code).Error
} 