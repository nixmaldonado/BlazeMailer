package delivery

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/nixmaldonado/blazeMailer/internal/model"
)

type Sender interface {
	Send(ctx context.Context, email model.EmailRequest) error
}

type SMTPSender struct {
	Host     string
	Port     string
	Password string
}

func (s SMTPSender) Send(_ context.Context, e model.EmailRequest) error {
	auth := smtp.PlainAuth("", e.From, s.Password, s.Host)
	msg := fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v\r\n", e.From, e.To, e.Subject, e.Body)
	return smtp.SendMail(s.Host+":"+s.Port, auth, e.From, []string{e.To}, []byte(msg))
}
