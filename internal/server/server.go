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
}

func openDB(dsn string) (*sql.DB, error) {
	var dbUrl = os.Getenv("MAIN_SQL_URL")
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}
	var token = "from Go token"
	var b = "from Go B data"
	var expiry = time.Time{}
	_, err = db.Exec("REPLACE INTO sessions (token, data, expiry) VALUES (?, ?, julianday($3))", token, b, expiry.UTC().Format("2006-01-02T15:04:05.999"))
	fmt.Println("AFTER EXEC", token)
	if err != nil {
		fmt.Println("GOT TO ERROR", err)
	}

	// db, err := sql.Open("sqlite3", dsn)
	// if err != nil {
	// 	return nil, err
	// }

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewServer() *http.Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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

	stytchAPIClient, err := stytchapi.NewClient(
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
	)
	if err != nil {
		panic(err)
	}

	// List all users for current application

	sessionManager := scs.New()
	sessionManager.Store = libsqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		users:           &models.UserModel{DB: db},
		logger:          logger,
		formDecoder:     formDecoder,
		sessionManager:  sessionManager,
		stytchAPIClient: stytchAPIClient,
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
