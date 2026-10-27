// Command searchd is the search engine's HTTP server.
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/pratham-rajesh/search-engine-from-scratch-in-GO/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("searchd failed", "err", err)
		os.Exit(1)
	}
}

// run holds the real startup logic. It returns errors instead of exiting so
// that deferred cleanup runs; os.Exit in main skips defers.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: newMux(),
		// Without this, a client that dribbles headers can hold a connection open forever.
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("searchd listening", "addr", srv.Addr, "log_level", cfg.LogLevel)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

// newMux builds the route table. It is separate from run so tests can
// exercise the routes without starting a real server.
func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}
