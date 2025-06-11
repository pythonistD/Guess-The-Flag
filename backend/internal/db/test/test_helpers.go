package test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"strings"
	"testing"
)

// newTestDB connect to db and migrate for tests
func newTestDB(t *testing.T) *sqlx.DB {
	connStr := "postgres://postgres:postgres@localhost:5432/guess_the_flag?sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	migrationsDir := "../../../migrations"
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("goose: failed to set dialect: %v", err)
	}
	if err := goose.Up(db.DB, migrationsDir); err != nil {
		t.Fatalf("goose: failed to apply migrations: %v", err)
	}
	return db
}

func clearTables(t *testing.T, db *sqlx.DB, tables []string) {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", strings.Join(tables, ","))
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}
