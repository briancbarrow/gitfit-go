package tenancy

import (
	"database/sql"
	"embed"
	"encoding/csv"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed *.sql *.csv
var embedFiles embed.FS

func InsertDataFromCSV(dbUrl string) error {
	csvFile := "exercises.csv"
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO exercises(id, name, target_area) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	file, err := embedFiles.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		_, err = stmt.Exec(id, record[1], record[2])
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func CreateTenantTables(dbUrl string) error {

	// https://github.com/libsql/libsql-client-go/#open-a-connection-to-sqld
	// libsql://[your-database].turso.io?authToken=[your-auth-token]
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}
	goose.SetBaseFS(embedFiles)
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}
