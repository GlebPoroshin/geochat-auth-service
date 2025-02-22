package email

import (
	"fmt"
	"net/smtp"
)

type emailService struct {
	config *config.Config
}

func NewEmailService(config *config.Config) service.EmailService {
	return &emailService{
		config: config,
	}
}

func (s *emailService) SendVerificationCode(email, code string) error {
	subject := "Email Verification"
	body := fmt.Sprintf("Your verification code is: %s\nIt will expire in 15 minutes.", code)
	return s.sendEmail(email, subject, body)
}

func (s *emailService) SendPasswordResetCode(email, code string) error {
	subject := "Password Reset"
	body := fmt.Sprintf("Your password reset code is: %s\nIt will expire in 15 minutes.", code)
	return s.sendEmail(email, subject, body)
}

func (s *emailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.config.SMTPUsername, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.SMTPUsername, []string{to}, []byte(msg))
} 