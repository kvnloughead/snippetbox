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

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
