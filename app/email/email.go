package email

import (
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/wneessen/go-mail"
)

type Email struct {
	Client    *mail.Client
	EmailFrom string
}

func New(cfg *config.Config) (*Email, error) {
	var client *mail.Client
	var err error
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		client, err = mail.NewClient(
			cfg.SMTPHost,
			mail.WithPort(cfg.SMTPPort),
			mail.WithTLSPortPolicy(mail.NoTLS),
		)
	} else {
		client, err = mail.NewClient(
			cfg.SMTPHost,
			mail.WithPort(cfg.SMTPPort),
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.SMTPUsername),
			mail.WithPassword(cfg.SMTPPassword),
		)
	}
	if err != nil {
		return nil, err
	}

	return &Email{
		Client:    client,
		EmailFrom: cfg.EmailFrom,
	}, nil
}

func (e *Email) Send(to string, subject string, text string) error {
	message := mail.NewMsg()
	if err := message.From(e.EmailFrom); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, text)

	return e.Client.DialAndSend(message)
}
