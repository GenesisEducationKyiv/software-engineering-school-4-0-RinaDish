package services

import (
	"crypto/tls"
	"net/smtp"

	"go.uber.org/zap"
)

const serverName = "smtp.gmail.com"
type Mail struct {
	from string
	pass string
	client   *smtp.Client
	l *zap.SugaredLogger
}

func NewEmail(from string, pass string, l *zap.SugaredLogger) (*Mail, error) {
	client, err := smtp.Dial("smtp.gmail.com:587")
	if err != nil {
		return nil, err
	}

	tlsconfig := &tls.Config{
		ServerName: serverName,
	}
	
	if err = client.StartTLS(tlsconfig); err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", from, pass, serverName)
	if err = client.Auth(auth); err != nil {
		return nil, err
	}
	
	return &Mail{
		from: from,
		pass: pass,
		client: client,
		l: l,
	}, nil
}

func (m Mail) Send(to, body string) {
	if m.client == nil {
		m.l.Infof("SMTP client is not connected")
		return
	}

	from := m.from

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Dollar Rate\n\n" +
		body

	if err := m.client.Mail(from); err != nil {
		m.l.Info(err)
		return
	}

	if err := m.client.Rcpt(to); err != nil {
		m.l.Info(err)
		return
	}

	w, err := m.client.Data()
	if err != nil {
		m.l.Info(err)
		return
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		m.l.Info(err)
		return
	}

	err = w.Close()
	if err != nil {
		m.l.Info(err)
		return
	}
	
	m.l.Infof("email sent to %s", to)
}