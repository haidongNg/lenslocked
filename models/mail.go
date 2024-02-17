package models

import (
	"fmt"
	"net/smtp"
)

const DefaultSender = "support@lenslocked.dev"

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type EmailService struct {
	DefaultSender string
	Auth          *smtp.Auth
}

func NewMailService(config SMTPConfig) *EmailService {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	es := EmailService{
		DefaultSender: DefaultSender,
		Auth:          &auth,
	}
	return &es
}

func (es *EmailService) Send(email Email) error {
	to := []string{email.To}
	msg := []byte("From: " + email.From + "\r\n" +
		"To: " + email.To + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" + email.Plaintext + "\r\n")
	err := smtp.SendMail("sandbox.smtp.mailtrap.io:25", *es.Auth, es.DefaultSender, to, msg)

	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (es *EmailService) ForgotPassword(to, requestUrl string) error {
	email := Email{
		From:      DefaultSender,
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset password, please visit following link: " + requestUrl,
		HTML:      `<p>To reset password, please visit following link: <a></a></p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}

	return err
}
