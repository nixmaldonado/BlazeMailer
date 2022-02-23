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
	message := []byte("My super secret message.")

	pwd := config.Spec.SMTPPassword
	host := config.Spec.SMTPHost
	port := config.Spec.SMTPPort

	// Create authentication
	auth := smtp.PlainAuth("", "TOOD", pwd, host)

	// Send actual message
	err := smtp.SendMail(host+":"+port, auth, "from", []string{"to"}, message)
	if err != nil {
		log.Fatal(err)
	}
}
