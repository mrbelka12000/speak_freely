package mail

import (
	"fmt"

	gomail "gopkg.in/mail.v2"

	"github.com/mrbelka12000/linguo_sphere_backend/pkg/config"
)

type (
	Client struct {
		smtpHost string
		smtpPort int
		smtpUser string
		smtpPass string
	}

	Request struct {
		To      string
		Body    string
		Subject string
	}
)

func New(cfg config.Config) *Client {
	return &Client{
		smtpUser: cfg.SenderEmail,
		smtpHost: cfg.SMTPHost,
		smtpPort: cfg.SMTPPort,
		smtpPass: cfg.SMTPPassword,
	}
}

func (c *Client) Send(req Request) error {
	message := gomail.NewMessage()
	message.SetHeader("From", c.smtpUser)
	message.SetHeader("To", req.To)
	message.SetHeader("Subject", req.Subject)
	message.SetBody("text/plain", req.Body)

	dialer := gomail.NewDialer(c.smtpHost, c.smtpPort, c.smtpUser, c.smtpPass)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("dial mail fail: %v", err)
	}

	return nil
}
