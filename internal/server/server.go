package server

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/briancbarrow/gitfit-go/internal/models"
	"github.com/briancbarrow/gitfit-go/internal/prettylog"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger         *slog.Logger
	users          *models.UserModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewServer() *http.Server {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "local.db", "SQLite DB name")

	logger := slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		AddSource: true,
	}))
	formDecoder := form.NewDecoder()

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		users:          &models.UserModel{DB: db},
		logger:         logger,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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
