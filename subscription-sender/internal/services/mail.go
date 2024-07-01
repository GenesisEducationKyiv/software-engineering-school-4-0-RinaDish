package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/RinaDish/subscription-sender/tools"
)

const serverName = "smtp.gmail.com"
const serverPort = "587"

type Mail struct {
	from   	string
	pass   	string
	client 	*smtp.Client
	logger  tools.Logger
}

func NewEmail(from string, pass string, logger tools.Logger) (*Mail, error) {
	m := &Mail{
		from: from,
		pass: pass,
		logger:	logger,
	}

	if err := m.initClient(); err != nil {
		return nil, err
	}

	return m, nil
}

func (mail *Mail) initClient() error {
	client, err := smtp.Dial(serverName + ":" + serverPort)
	if err != nil {
		return err
	}

	if err = client.StartTLS(&tls.Config{
		ServerName: serverName,
	}); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", mail.from, mail.pass, serverName)
	if err = client.Auth(auth); err != nil {
		return err
	}

	mail.client = client
	
	return nil
}

func (mail *Mail) Send(to, body string) {
	if mail.client == nil {
		mail.logger.Infof("SMTP client is not connected")
		return
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Dollar Rate\n\n%s", mail.from, to, body)

	if err := mail.client.Mail(mail.from); err != nil {
		mail.logger.Info(err)
		return
	}

	if err := mail.client.Rcpt(to); err != nil {
		mail.logger.Info(err)
		return
	}

	w, err := mail.client.Data()
	if err != nil {
		mail.logger.Info(err)
		return
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		mail.logger.Info(err)
		return
	}

	err = w.Close()
	if err != nil {
		mail.logger.Info(err)
		return
	}

	mail.logger.Infof("email sent to %s", to)
}
