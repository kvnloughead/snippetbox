package models

import (
	"database/sql"
	"errors"
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
// If no matching snippet is found, a models.ErrNoRecord is returned.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Executes a query statement that will return no more than one row.
	// Accepts the query statement and a variadic list of placeholder values.
	row := m.DB.QueryRow(query, id)

	// Declare an empty snippet and populate it from the row returned by QueryRow.
	// If multiple rows were found, the first row is used.
	// If no rows were found, an sql.ErrNoRows error is returned.
	var s Snippet
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Query will return an sql.Rows result set containing 10 latest entries.
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // don't defer closing until after handling the error

	// Iterate through result set, calling rows.Scan on each row. Create a snippet
	// Create a snippet for each row and add it to the snippets slice.
	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// rows.Err() contains any errors that occurred during iteration, including
	// including errors that wouldn't be returned by rows.Scan().
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
