package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nixmaldonado/blazeMailer/internal/model"
	"github.com/nixmaldonado/blazeMailer/internal/queue"
)

type fakeQueue struct {
	jobs []model.Job
	full bool
}

func (f *fakeQueue) Enqueue(_ context.Context, j model.Job) error {
	if f.full {
		return queue.ErrQueueFull
	}
	f.jobs = append(f.jobs, j)
	return nil
}
func (f *fakeQueue) Dequeue(context.Context) (model.Job, error)          { return model.Job{}, nil }
func (f *fakeQueue) DeadLetter(context.Context, model.Job, string) error { return nil }
func (f *fakeQueue) Depth(context.Context) (int, error)                  { return len(f.jobs), nil }
func (f *fakeQueue) Close() error                                        { return nil }

func postEmail(srv *Server, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/emails", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	return rec
}

func validBody() []byte {
	b, _ := json.Marshal(model.EmailRequest{From: "a@b.com", To: "c@d.com", Subject: "hi", Body: "yo"})
	return b
}

func TestSendEmail_Accepted(t *testing.T) {
	fq := &fakeQueue{}
	srv := NewServer(fq)
	rec := postEmail(srv, validBody())
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}
	var resp model.SendEmailResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.JobID == "" {
		t.Fatal("expected non-empty job_id")
	}
	if len(fq.jobs) != 1 {
		t.Fatalf("enqueued %d jobs, want 1", len(fq.jobs))
	}
}

func TestSendEmail_ValidationError(t *testing.T) {
	srv := NewServer(&fakeQueue{})
	rec := postEmail(srv, []byte(`{"to":"c@d.com"}`))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestSendEmail_QueueFull(t *testing.T) {
	srv := NewServer(&fakeQueue{full: true})
	rec := postEmail(srv, validBody())
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}
