package models

import (
	"database/sql"
	"time"
)

// Type of data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// This wrapper for our sql.DB connection pool contains a number of helpful
// methods for interacting with the DB.
type SnippetModel struct {
	DB *sql.DB
}

// Inserts a new snippet into the DB. Returns the ID of the inserted record or
// an error.
func (m *SnippetModel) Insert(
	title string,
	content string,
	expires int) (int, error) {

	// Query statements allow `?` for placeholders. Values for placeholders are
	// supplied as arguments to sql.DB.Exec.
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// sql.DB.Exec accepts the query statement, and variadic placeholder values.
	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get ID of the inserted record as an int64.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a snippet by its ID.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
