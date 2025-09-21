// apps/cp-api/internal/app/app.go
package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"xeed/apps/cp-api/internal/config"
)

func Run() error {
	// 1) Load config
	cfg := config.FromEnv()

	// 2) Context untuk graceful shutdown (Ctrl+C / SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 3) Build seluruh dependency & handler (di wire.go)
	handler, cleanup, err := buildHTTP(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	// 4) Start HTTP server
	srv := &http.Server{
		Addr:              cfg.Addr, // ex: ":8080"
		Handler:           handler,  // dari routers
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("[cp-api] listening on %s", cfg.Addr)

	// 5) Graceful shutdown
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.Println("[cp-api] shutting down...")
		shCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		errCh <- srv.Shutdown(shCtx)
	}()

	// 6) Serve (blocking)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-errCh
}
