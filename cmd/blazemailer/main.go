package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nixmaldonado/blazeMailer/internal/api"
	"github.com/nixmaldonado/blazeMailer/internal/config"
	"github.com/nixmaldonado/blazeMailer/internal/delivery"
	"github.com/nixmaldonado/blazeMailer/internal/queue"
	"github.com/nixmaldonado/blazeMailer/internal/worker"
)

func main() {
	config.InitConfig()
	q := queue.NewMemory(1024)
	sender := delivery.SMTPSender{
		Host:     config.Spec.SMTPHost,
		Port:     config.Spec.SMTPPort,
		Password: config.Spec.SMTPPassword,
	}
	pool := worker.NewPool(q, sender, 8)

	// The pool runs under its own context so that on shutdown workers keep
	// draining buffered jobs until the queue is closed, rather than being
	// cancelled mid-flight. poolCancel is the hard-stop fallback.
	poolCtx, poolCancel := context.WithCancel(context.Background())
	defer poolCancel()
	poolDone := make(chan struct{})
	go func() { pool.Run(poolCtx); close(poolDone) }()

	srv := &http.Server{Addr: ":3000", Handler: api.NewServer(q).Routes()}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Println("listening on :3000 (memory mode)")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("shutting down...")

	// 1. Stop accepting new requests.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}

	// 2. Close the queue so workers drain remaining buffered jobs and exit.
	q.Close()

	// 3. Wait for the pool to finish draining, with a hard-stop fallback.
	select {
	case <-poolDone:
		log.Println("workers drained")
	case <-time.After(10 * time.Second):
		log.Println("drain timeout; forcing worker stop")
		poolCancel()
		<-poolDone
	}
}
