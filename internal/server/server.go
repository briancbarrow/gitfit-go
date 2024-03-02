package server

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	// "github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/libsqlstore"
	"github.com/alexedwards/scs/v2"
	database "github.com/briancbarrow/gitfit-go/internal/database/db"
	"github.com/briancbarrow/gitfit-go/internal/models"
	"github.com/briancbarrow/gitfit-go/internal/prettylog"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/stytchapi"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger          *slog.Logger
	users           *models.UserModel
	formDecoder     *form.Decoder
	sessionManager  *scs.SessionManager
	stytchAPIClient *stytchapi.API
	queries         *database.Queries
}

func openDB() (*sql.DB, error) {
	var dbUrl = os.Getenv("MAIN_SQL_URL")
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}

	// create connection to check for errors
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewServer(isProd bool) *http.Server {
	var filepath string
	if isProd {
		filepath = ".env"
	} else {
		filepath = ".env.local"
	}

	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr := flag.String("addr", ":4000", "HTTP network address")

	// TODO: Look into if there is a better way to do prettylog
	logger := slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		AddSource: true,
	}))
	formDecoder := form.NewDecoder()

	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	stytchAPIClient, err := stytchapi.NewClient(
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
	)
	if err != nil {
		panic(err)
	}

	queries := database.New(db)

	sessionManager := scs.New()
	sessionManager.Store = libsqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		users:           &models.UserModel{DB: db},
		logger:          logger,
		formDecoder:     formDecoder,
		sessionManager:  sessionManager,
		stytchAPIClient: stytchAPIClient,
		queries:         queries,
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
