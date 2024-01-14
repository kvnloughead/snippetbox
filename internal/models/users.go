package models

import (
	"database/sql"
	"time"
)

// Type representing our user document.
type User struct {
	ID              int
	Name            string
	Email           string
	Hashed_password []byte // a bcrypt hash
	Created         time.Time
}

// A wrapper for our sql.DB connection pool.
// Contains methods for interacting with the users collection.
type UserModel struct {
	DB *sql.DB
}

// Inserts a new user user the DB.
// Returns the ID of the inserted record or an error.
func (m *UserModel) Insert(
	name string,
	email string,
	password string,
) error {
	return nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
