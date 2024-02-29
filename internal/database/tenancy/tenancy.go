package tenancy

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateRandomDBString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func CreateTenantDB(databaseID string) (bool, error) {
	fmt.Println("Creating tenant DB", databaseID)
	body := []byte(fmt.Sprintf(`{
		"name":  "%s",
		"group": "default"
		}`, databaseID))

	org := os.Getenv("TURSO_ORG")
	fmt.Println("Org", org)
	fmt.Println("BODY", string(body))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.turso.tech/v1/organizations/%s/databases", org), bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request", err)
		slog.Error("Error creating request", "err", err)
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TURSO_API_TOKEN")))
	fmt.Println("Request", req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		slog.Error("Error sending request", "err", err)
		return false, err
	}
	fmt.Println("Response", resp)
	if resp.StatusCode == 200 {
		slog.Info("Tenant DB created successfully")
		return true, nil
	}
	return false, nil
}

func CheckIfTenantDBExists(databaseID string) (bool, error) {
	org := os.Getenv("TURSO_ORG")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.turso.tech/v1/organizations/%s/databases/%s", org, databaseID), nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		slog.Error("Error creating request", "err", err)
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TURSO_API_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		slog.Error("Error sending request", "err", err)
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func SeedTenantDB(databaseID string) error {
	dbUrl := fmt.Sprintf("libsql://%s-briancbarrow.turso.io?authToken=%s", databaseID, os.Getenv("TURSO_DB_TOKEN"))
	fmt.Println("Running migrations on", dbUrl)
	err := CreateTenantTables(dbUrl)
	if err != nil {
		return err
	}
	err = InsertDataFromCSV(dbUrl)
	if err != nil {
		return err
	}
	return nil
}

func OpenTenantDB(databaseID string) (*sql.DB, error) {
	authToken := os.Getenv("TURSO_DB_TOKEN")
	dbUrl := fmt.Sprintf("libsql://%s-briancbarrow.turso.io?authToken=%s", databaseID, authToken)
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
