package model

import (
	"errors"
	"net/mail"
)

type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (r EmailRequest) Validate() error {
	if r.From == "" {
		return errors.New("from is required")
	}
	if _, err := mail.ParseAddress(r.From); err != nil {
		return errors.New("from is not a valid email address")
	}
	if r.To == "" {
		return errors.New("to is required")
	}
	if _, err := mail.ParseAddress(r.To); err != nil {
		return errors.New("to is not a valid email address")
	}
	return nil
}

type SendEmailResponse struct {
	Success bool   `json:"success"`
	JobID   string `json:"job_id,omitempty"`
	Error   string `json:"error,omitempty"`
}
