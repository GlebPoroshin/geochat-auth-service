package dto

type RegisterRequest struct {
    Login    string `json:"login"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    LoginOrEmail string `json:"login_or_email"`
    Password     string `json:"password"`
}

type VerifyRegistrationRequest struct {
    UserID string `json:"user_id"`
    Code   string `json:"code"`
}

type PasswordResetRequest struct {
    Email string `json:"email"`
}

type VerifyPasswordResetRequest struct {
    Email string `json:"email"`
    Code  string `json:"code"`
}

type ResetPasswordRequest struct {
    Email       string `json:"email"`
    NewPassword string `json:"new_password"`
}

type RefreshTokenRequest struct {
    UserID       string `json:"user_id"`
    RefreshToken string `json:"refresh_token"`
} 