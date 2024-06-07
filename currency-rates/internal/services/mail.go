package services

import (
	"net/smtp"

	"go.uber.org/zap"
)

type Mail struct {
	from string
	pass string
	l *zap.SugaredLogger
}

func NewEmail(from string, pass string, l *zap.SugaredLogger) *Mail {
	return &Mail{
		from: from,
		pass: pass,
		l: l,
	}
}

func (m Mail) Send(to, body string) {
	pass := m.pass
	from := m.from

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Dollar Rate\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		m.l.Infof("smtp error: %s", err)
		return
	}
	
	m.l.Infof("email sended to %s", to)
}