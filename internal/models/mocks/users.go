package mocks

import (
	"time"

	"github.com/kvnloughead/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@mail.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	if email == "testuser@mail.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (models.User, error) {
	if id == 1 {
		u := models.User{
			ID:      1,
			Name:    "User",
			Email:   "testuser@mail.com",
			Created: time.Now(),
		}
		return u, nil
	} else {
		return models.User{}, models.ErrNoRecord
	}
}

func (m *UserModel) PasswordUpdate(id int, password string) error {
	return nil
}
