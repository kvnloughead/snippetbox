package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	assert "github.com/kvnloughead/snippetbox/internal"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	// Create HTTPS server for testing on random port.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	// Verify the status code and body.
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)

	// Create HTTPS server for testing on random port.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, header, _ := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		formTag := regexp.MustCompile(`<form .* action="/snippet/create" method="POST">`)
		// GET login page and extract CSRF token.
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		// Fill out login form.
		form := url.Values{}
		form.Add("email", "testuser@mail.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", csrfToken)

		// Send POST /user/login request, and check that the authenticated user is shown the create snippet form.
		code, header, _ := ts.post(t, "/user/login", form)
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/snippet/create")

		// GET /snippet/create and verify response.
		code, _, body = ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContainsMatch(t, body, formTag)
	})
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		id       string
		wantCode int
		wantBody string
	}{
		{
			name:     "Existing",
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "This is a mock snippet.",
		},
		{
			name:     "Non-existing ID",
			id:       "999",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Zero ID",
			id:       "0",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			id:       "-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Non-integer ID",
			id:       "1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Non-number ID",
			id:       "foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			id:       "",
			wantCode: http.StatusNotFound,
		},
	}

	for _, sub := range tests {
		t.Run(sub.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, "/snippet/view/"+fmt.Sprint(sub.id))
			assert.Equal(t, statusCode, sub.wantCode)
			if sub.wantBody != "" {
				assert.StringContains(t, body, sub.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET request to /user/signup and extract the CSRF token from the response's body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Name"
		validPassword = "pa$$word"
		validEmail    = "email@mail.com"
	)

	// When request fails, the form tag should still be present on the page.
	formTag := regexp.MustCompile(`<form .* action="/user/signup" method="POST" novalidate>`)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  *regexp.Regexp
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "badToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bad@email.",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pass",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@mail.com",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
		},
	}

	for _, sub := range tests {
		t.Run(sub.name, func(t *testing.T) {
			// Populate form values.
			form := url.Values{}
			form.Add("name", sub.userName)
			form.Add("email", sub.userEmail)
			form.Add("password", sub.userPassword)
			form.Add("csrf_token", sub.csrfToken)

			code, _, body := ts.post(t, "/user/signup", form)

			assert.Equal(t, code, sub.wantCode)
			if sub.wantFormTag != nil {
				assert.StringContainsMatch(t, body, sub.wantFormTag)
			}
		})
	}
}
