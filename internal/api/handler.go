package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nixmaldonado/blazeMailer/internal/delivery"
	"github.com/nixmaldonado/blazeMailer/internal/model"
)

type Server struct {
	sender delivery.Sender
}

func NewServer(s delivery.Sender) *Server { return &Server{sender: s} }

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", s.ping)
	r.Post("/v1/emails", s.sendEmail)
	return r
}

func (s *Server) ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ping"))
}

func (s *Server) sendEmail(w http.ResponseWriter, r *http.Request) {
	var req model.EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.sender.Send(r.Context(), req); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, model.SendEmailResponse{Success: true})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, model.SendEmailResponse{Success: false, Error: msg})
}
