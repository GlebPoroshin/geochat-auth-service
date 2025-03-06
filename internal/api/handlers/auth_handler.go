package handlers

import (
	"github.com/GlebPoroshin/geochat-auth-service/internal/api/dto"
	"github.com/GlebPoroshin/geochat-auth-service/internal/service"
	sharedJWT "github.com/GlebPoroshin/geochat-shared/jwt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Создает нового пользователя по логину, email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} map[string]interface{} "Успешная регистрация"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, err := h.authService.Register(c.Context(), req.Login, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": userID,
		"message": "Registration successful. Please check your email for verification code.",
	})
}

// Login выполняет вход пользователя
// @Summary Вход пользователя
// @Description Авторизует пользователя по логину или email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Успешный вход"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.authService.Login(c.Context(), req.LoginOrEmail, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

// VerifyRegistration подтверждает регистрацию
// @Summary Подтверждение регистрации
// @Description Подтверждает email пользователя с помощью кода верификации
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyRegistrationRequest true "Данные для верификации"
// @Success 200 {object} map[string]string "Email подтвержден"
// @Failure 400 {object} map[string]string "Ошибка в коде верификации"
// @Router /auth/verify-registration [post]
func (h *AuthHandler) VerifyRegistration(c *fiber.Ctx) error {
	var req dto.VerifyRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.authService.VerifyRegistration(c.Context(), req.UserID, req.Code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
	})
}

// InitiatePasswordReset отправляет код для сброса пароля
// @Summary Запрос на сброс пароля
// @Description Отправляет код сброса пароля на email пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.PasswordResetRequest true "Email пользователя"
// @Success 200 {object} map[string]string "Код сброса отправлен"
// @Failure 400 {object} map[string]string "Ошибка при отправке"
// @Router /auth/password-reset [post]
func (h *AuthHandler) InitiatePasswordReset(c *fiber.Ctx) error {
	var req dto.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.authService.InitiatePasswordReset(c.Context(), req.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset code sent to email",
	})
}

// VerifyPasswordResetCode проверяет код сброса пароля
// @Summary Проверка кода сброса пароля
// @Description Проверяет корректность введенного кода для сброса пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyPasswordResetRequest true "Email и код сброса"
// @Success 200 {object} map[string]string "Код сброса подтвержден"
// @Failure 400 {object} map[string]string "Неверный код"
// @Router /auth/verify-reset-code [post]
func (h *AuthHandler) VerifyPasswordResetCode(c *fiber.Ctx) error {
	var req dto.VerifyPasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.authService.VerifyPasswordResetCode(c.Context(), req.Email, req.Code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reset code verified successfully",
	})
}

// ResetPassword сбрасывает пароль
// @Summary Сброс пароля
// @Description Устанавливает новый пароль после верификации кода
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Email и новый пароль"
// @Success 200 {object} map[string]string "Пароль успешно сброшен"
// @Failure 400 {object} map[string]string "Ошибка при сбросе пароля"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.authService.ResetPassword(c.Context(), req.Email, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// RefreshToken обновляет токен доступа
// @Summary Обновление токена
// @Description Обновляет токен доступа, используя refresh-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer RefreshToken"
// @Success 200 {object} map[string]interface{} "Новый токен"
// @Failure 401 {object} map[string]interface{} "Ошибка авторизации"
// @Router /auth/refresh [get]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	parts := strings.Split(refreshToken, " ")

	claims, _ := sharedJWT.Validate(parts[1])
	userID := claims.Subject

	response, err := h.authService.RefreshToken(c.Context(), userID, parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
