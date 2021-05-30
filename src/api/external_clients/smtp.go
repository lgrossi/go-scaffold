package external_clients

import (
	"github.com/lgrossi/go-scaffold/src/configs"
	"github.com/lgrossi/go-scaffold/src/logger"
	gomail "gopkg.in/mail.v2"
)

const (
	EnvKeySmtpHost     = "SMTP_HOST"
	EnvKeySmtpPort     = "SMTP_PORT"
	EnvKeySmtpUser     = "SMTP_USER"
	EnvKeySmtpPassword = "SMTP_PASS"
)

func SendEmail(to, subject, body string) {
	m := gomail.NewMessage()
	from := configs.GetEnvStr(EnvKeySmtpUser)

	m.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      {to},
		"Subject": {subject},
	})

	m.SetBody(
		"text/html",
		body,
	)

	d := gomail.NewDialer(
		configs.GetEnvStr(EnvKeySmtpHost, "smtp.gmail.com"),
		configs.GetEnvInt(EnvKeySmtpPort, 465),
		from,
		configs.GetEnvStr(EnvKeySmtpPassword),
	)

	if err := d.DialAndSend(m); err != nil {
		logger.Panic(err)
	}
}
