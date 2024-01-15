package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// Generate hash from the password with bcrypt.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// The query to be executed. Query statements allow for '?' as placeholders.
	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Execute query. Exec accepts variadic values for the query placeholders.
	_, err = m.DB.Exec(query, name, email, string(hash))
	if err != nil {
		// Handle duplicate email errors.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// If errors.As() is true, the error will be assigned to mySQLError.
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate a user on login by comparing the plain text password to the
// user's stored hashed password. If the email or password is incorrect, an
// ErrInvalidCredentials error is returned.
func (m *UserModel) Authenticate(email string, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := "SELECT id, hashed_password FROM users WHERE email = ?"

	// QueryRow returns the first matching row. Scan copies the columns of the
	// matched row into the specified locations. Scan returns ErrNoRows if no
	// match was found.
	err := m.DB.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Compare password to hash. If they don't match return ErrInvalidCredentials.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// If password is correct, return the user's ID.
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
