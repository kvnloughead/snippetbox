package models

import (
	"testing"

	assert "github.com/kvnloughead/snippetbox/internal"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Existing ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existing ID",
			userID: 2,
			want:   false,
		},
	}

	for _, sub := range tests {
		t.Run(sub.name, func(t *testing.T) {
			// Each test sets runs the setup and teardown scripts.
			db := newTestDB(t)
			m := UserModel{db}

			exists, err := m.Exists(sub.userID)

			assert.Equal(t, exists, sub.want)
			assert.IsNil(t, err)
		})
	}

}
