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

type UserModelInterface interface {
	Authenticate(email string, password string) (int, error)
	Get(id int) (User, error)
	Exists(id int) (bool, error)
	Insert(name, email, password string) error
	PasswordUpdate(id int, password string) error
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

// Get a user by its ID.
// If no matching snippet is found, a models.ErrNoRecord error is returned.
func (m *UserModel) Get(id int) (User, error) {
	query := `SELECT name, email, created FROM users
	WHERE id = ?`

	// Executes a query statement that will return no more than one row.
	// Accepts the query statement and a variadic list of placeholder values.
	row := m.DB.QueryRow(query, id)

	// Declare an empty user struct and populate it from the returned row.
	// If no rows were found, an sql.ErrNoRows error is returned.
	// If multiple rows were found, the first row is used.
	var u User
	err := row.Scan(&u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		} else {
			return User{}, err
		}
	}

	return u, nil
}

// Returns true if a user with the given ID is found in the database.
//
// In normal circumstances the error returned will always be nil, because the sql EXISTS statement always returns a row, even when there is a match.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}

// Inserts a new user user the DB.
// Returns the ID of the inserted record or an error.
func (m *UserModel) Insert(name, email, password string) error {
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

// Generates a hash from the supplied password and updates it in the DB.
// The password is not validated, so make sure that it is valid before calling.
func (m *UserModel) PasswordUpdate(id int, password string) error {
	// Generate hash from the password with bcrypt.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET hashed_password = ? WHERE id = ?`
	m.DB.Exec(stmt, hash, id)
	return nil
}
