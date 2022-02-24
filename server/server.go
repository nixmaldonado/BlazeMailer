package server

import (
	"github.com/nixmaldonado/blazeMailer/config"
	"log"
	"net/http"
	"net/smtp"
)

func handleRequests() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/email/send", sendEmail)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	cfg := config.Spec
	password := cfg.SMTPPassword
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

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
