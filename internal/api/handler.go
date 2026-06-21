package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nixmaldonado/blazeMailer/internal/model"
	"github.com/nixmaldonado/blazeMailer/internal/queue"
)

type Server struct {
	queue queue.Queue
}

func NewServer(q queue.Queue) *Server { return &Server{queue: q} }

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
	job := model.Job{
		ID:         uuid.NewString(),
		Email:      req,
		EnqueuedAt: time.Now(),
	}
	if err := s.queue.Enqueue(r.Context(), job); err != nil {
		if errors.Is(err, queue.ErrQueueFull) || errors.Is(err, queue.ErrQueueClosed) {
			writeErr(w, http.StatusServiceUnavailable, "queue unavailable, retry later")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, model.SendEmailResponse{Success: true, JobID: job.ID})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, model.SendEmailResponse{Success: false, Error: msg})
}
