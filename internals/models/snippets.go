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

// Insert a new snippet into the DB.
func (m *SnippetModel) Insert(
	title string, content string, expires int,
) (int, error) {
	return 0, nil
}

// Get a snippet by its ID.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
