package server

import (
	"encoding/json"
	"fmt"
	"github.com/nixmaldonado/blazeMailer/config"
	"github.com/nixmaldonado/blazeMailer/models"
	"log"
	"net/http"
	"net/smtp"
)

func HandleRequests() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/email/send", sendEmail)
	log.Println("listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping"))
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	var p models.SendEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	cfg := config.Spec
	password := cfg.SMTPPassword
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	// Create authentication
	auth := smtp.PlainAuth("", p.From, password, smtpHost)

	// Format message
	mail := fmt.Sprintf("To: %v\r\nSubject: %v\r\n\r\n%v\r\n", p.To, p.Subject, p.Body)

	// Send actual message
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, p.From, []string{p.To}, []byte(mail)); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(models.SendEmailResponse{Success: true, Recipient: p.To})
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.SendEmailResponse{Success: false, Error: err.Error()})
}
