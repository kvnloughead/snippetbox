package models

import (
	"database/sql"
	"os"
	"testing"
)

// Reads script from the provided file path and executes it, closing the pool and calling t.Fatal() in case of error.
func execScript(t *testing.T, db *sql.DB, scriptPath string) {
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
}

// Establishes an sql.DB connection pool for test DB, and runs the setup.sql script to create the tables and a single user document. Also registers a cleanup function that closes the connection pool and runs teardown.sql when the calling test is finished running.
func newTestDB(t *testing.T) *sql.DB {
	// The multiStatements parameter is needed to run our setup and teardown sql scripts.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	execScript(t, db, "./testdata/setup.sql")

	// Cleanup function closes the connection pool and runs the teardown script.
	// This will be called when the test that called newTestDB is finished.
	t.Cleanup(func() {
		defer db.Close()
		execScript(t, db, "./testdata/teardown.sql")
	})

	return db
}
