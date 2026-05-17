package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer struct {
	host      string
	port      string
	user      string
	pass      string
	fromEmail string
}

func New(host, port, user, pass, fromEmail string) *Mailer {
	return &Mailer{
		host:      host,
		port:      port,
		user:      user,
		pass:      pass,
		fromEmail: fromEmail,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", m.fromEmail),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	return smtp.SendMail(addr, auth, m.fromEmail, []string{to}, []byte(msg))
}