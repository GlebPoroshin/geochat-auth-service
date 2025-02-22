type UserRepository interface {
    FindByLoginOrEmail(ctx context.Context, loginOrEmail string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
    UpdateVerificationStatus(ctx context.Context, userID string, verified bool) error
    UpdatePassword(ctx context.Context, userID string, hashedPassword string) error
}

type VerificationRepository interface {
    Create(ctx context.Context, code *models.VerificationCode) error
    FindValid(ctx context.Context, userID, code, codeType string) (*models.VerificationCode, error)
    Delete(ctx context.Context, code *models.VerificationCode) error
} 