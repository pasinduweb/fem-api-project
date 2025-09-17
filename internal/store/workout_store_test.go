package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5434 sslmode=disable")

	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// Migratoins for the test_db
	err = Migrate(db, "../../migrations/")

	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE users, workouts, workout_entries CASCADE`)

	if err != nil {
		t.Fatalf("truncating tables %v", err)
	}

	return db
}
