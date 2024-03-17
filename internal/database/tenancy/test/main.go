package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/briancbarrow/gitfit-go/internal/database/tenancy"
	"github.com/joho/godotenv"
)

func main() {
	dbId := flag.String("dbId", "", "The database id to run the migrations on")
	godotenv.Load()
	flag.Parse()
	authToken := os.Getenv("TURSO_DB_TOKEN")

	dbUrl := fmt.Sprintf("libsql://%s-briancbarrow.turso.io?authToken=%s", *dbId, authToken)

	err := tenancy.CreateTenantTables(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	err = tenancy.InsertDataFromCSV(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
}
