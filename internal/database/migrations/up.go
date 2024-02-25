package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed *.sql
var embedMigrations embed.FS

func main() {
	isProd := flag.Bool("isProd", false, "Determines if the server is in production mode")
	direction := flag.String("direction", "up", "Run up migrations")
	flag.Parse()
	fmt.Println("Starting migrations", *direction)
	var filepath string
	if *isProd {
		filepath = ".env"
	} else {
		filepath = ".env.local"
	}
	fmt.Println("Loading .env file", filepath)
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatal("failed to load env")
	}
	// https://github.com/libsql/libsql-client-go/#open-a-connection-to-sqld
	// libsql://[your-database].turso.io?authToken=[your-auth-token]
	sqlURL := os.Getenv("MAIN_SQL_URL")
	if sqlURL == "" {
		log.Fatal("No MAIN_SQL_URL set in .env")
	}
	db, err := sql.Open("libsql", sqlURL)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatal(err)
	}
	if *direction == "down" {
		if err := goose.Down(db, "."); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := goose.Up(db, "."); err != nil {
			log.Fatal(err)
		}
	}
}
