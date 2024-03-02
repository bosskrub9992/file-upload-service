package main

import (
	"log/slog"

	"github.com/bosskrub9992/file-upload-service/emails"

	"github.com/bosskrub9992/file-upload-service/configs"

	"gopkg.in/gomail.v2"
)

func main() {
	cfg := configs.New()
	dialer := gomail.NewDialer(
		cfg.Secret.Email.Host,
		cfg.Secret.Email.Port,
		cfg.Secret.Email.Username,
		cfg.Secret.Email.Password,
	)
	sender := emails.NewSender(cfg, dialer)

	if err := sender.Send("to@email.com", "testSubject", "testContent"); err != nil {
		slog.Error(err.Error())
		return
	}
}
