package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nixmaldonado/blazeMailer/internal/model"
)

type fakeSender struct {
	sent []model.EmailRequest
	err  error
}

func (f *fakeSender) Send(_ context.Context, e model.EmailRequest) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, e)
	return nil
}

func TestSendEmail_OK(t *testing.T) {
	fs := &fakeSender{}
	srv := NewServer(fs)
	body, _ := json.Marshal(model.EmailRequest{From: "a@b.com", To: "c@d.com", Subject: "hi", Body: "yo"})
	req := httptest.NewRequest(http.MethodPost, "/v1/emails", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if len(fs.sent) != 1 {
		t.Fatalf("sent %d emails, want 1", len(fs.sent))
	}
}

func TestSendEmail_ValidationError(t *testing.T) {
	srv := NewServer(&fakeSender{})
	body := []byte(`{"to":"c@d.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/emails", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
