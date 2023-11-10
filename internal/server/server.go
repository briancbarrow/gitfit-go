package server

import (
	"flag"
	"log/slog"
	"net/http"
	"time"

	"github.com/briancbarrow/gitfit-go/internal/prettylog"
)

type application struct {
	logger *slog.Logger
}

func NewServer() *http.Server {
	addr := flag.String("addr", ":4000", "HTTP network address")

	logger := slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		AddSource: true,
	}))
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("port", *addr))

	// Declare Server config
	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
