package repository

import (
	"context"

	"gorm.io/gorm"

	"geochat-auth-service/internal/models"
)

type UserRepository interface {
	FindByLoginOrEmail(ctx context.Context, loginOrEmail string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	UpdateVerificationStatus(ctx context.Context, userID string, verified bool) error
	UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByLoginOrEmail(ctx context.Context, loginOrEmail string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("login = ? OR email = ?", loginOrEmail, loginOrEmail).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateVerificationStatus(ctx context.Context, userID string, verified bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("verified", verified).Error
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
} 