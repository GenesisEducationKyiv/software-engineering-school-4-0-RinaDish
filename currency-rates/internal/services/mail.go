package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"go.uber.org/zap"
)

const serverName = "smtp.gmail.com"
const serverPort = "587"

type Mail struct {
	from   string
	pass   string
	client *smtp.Client
	l      *zap.SugaredLogger
}

func NewEmail(from string, pass string, l *zap.SugaredLogger) (*Mail, error) {
	m := &Mail{
		from: from,
		pass: pass,
		l:    l,
	}

	if err := m.initClient(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Mail) initClient() error {
	client, err := smtp.Dial(serverName + ":" + serverPort)
	if err != nil {
		return err
	}

	if err = client.StartTLS(&tls.Config{
		ServerName: serverName,
	}); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.from, m.pass, serverName)
	if err = client.Auth(auth); err != nil {
		return err
	}

	m.client = client
	
	return nil
}

func (m *Mail) Send(to, body string) {
	if m.client == nil {
		m.l.Infof("SMTP client is not connected")
		return
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Dollar Rate\n\n%s", m.from, to, body)

	if err := m.client.Mail(m.from); err != nil {
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
