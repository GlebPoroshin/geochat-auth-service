package models

import "time"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Login     string    `json:"login" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VerificationCode struct {
	UserID    string    `json:"user_id" gorm:"primaryKey"`
	Code      string    `json:"code"`
	Type      string    `json:"type"` // "registration" or "password_reset"
	ExpiresAt time.Time `json:"expires_at"`
}
