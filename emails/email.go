package emails

import (
	"github.com/bosskrub9992/file-upload-service/configs"

	"github.com/cenkalti/backoff/v4"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	cfg    *configs.Config
	dialer *gomail.Dialer
}

func NewSender(cfg *configs.Config, dialer *gomail.Dialer) *Sender {
	return &Sender{
		cfg:    cfg,
		dialer: dialer,
	}
}

func (s Sender) Send(to, subject, content string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.cfg.Email.From)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", content)
	if err := s.dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func (s Sender) SendUploadSuccess(to string) error {
	err := backoff.Retry(
		func() error {
			return s.Send(to, "upload succeeded", "upload succeeded")
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3),
	)
	if err != nil {
		return err
	}
	return nil
}
