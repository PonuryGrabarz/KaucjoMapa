package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
)

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

var _ port.IMailSender = (*SMTPSender)(nil)

func NewSMTPSender(host, port, username, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) SendPasswordRecoveryCode(_ context.Context, to, code string) error {
	subject := "KaucjoMapa - kod resetowania hasla"
	body := fmt.Sprintf(
		"Twoj kod resetowania hasla to: %s\r\n\r\nKod jest wazny przez 10 minut.\r\nJesli to nie Ty zglaszales reset hasla, zignoruj te wiadomosc.\r\n",
		code,
	)

	msg := strings.Join([]string{
		"From: " + s.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, []byte(msg))
}
