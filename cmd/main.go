package main

import (
	"github.com/nixmaldonado/blazeMailer/config"
	"log"
	"net/smtp"
)

func main() {
	config.InitConfig()

	password := config.Spec.SMTPPassword
	smtpHost := config.Spec.SMTPHost
	smtpPort := config.Spec.SMTPPort

	from := "nixsm@yahoo.com"
	to := "nixmaldonado@gmail.com"
	msg := []byte("To: nixmaldonado@gmail.com\r\n" +
		"Subject: Blaze Mail Hello World!\r\n" +
		"\r\n" +
		"This is the email body Blaze the world!\r\n")

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	recipients := []string{to}

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, msg)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func initMailer() {
	// TODO init server
}
