package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pavankumarhfl-hub/MeshGate/internal/config"
	"github.com/pavankumarhfl-hub/MeshGate/internal/gateway"
)

func main() {
	cfg := config.FromEnv()
	gw := gateway.New(cfg)

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           gw.Handler(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("meshgate listening on %s", cfg.ListenAddr)
		if err := gw.Run(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil { log.Fatal(err) }
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
			_ = server.Close()
		}
		log.Printf("meshgate stopped")
	}

	_ = time.Second // keeps shutdown timing visible in generated docs/examples
}
