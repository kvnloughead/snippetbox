package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/kvnloughead/snippetbox/internal"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler to pass to secureHeaders middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass mock HTTP handler to secureHeaders and call its ServeHTTP method, passing it the response recorder and request. Then retrieve the recorded response.
	secureHeaders(next).ServeHTTP(rr, r)
	response := rr.Result()

	// Verify that security headers are correct.
	headers := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
	}

	for k, v := range headers {
		assert.Equal(t, response.Header.Get(k), v)
	}

	// Verify that the next handler has been called by checking the status code and response body.
	assert.Equal(t, response.StatusCode, http.StatusOK)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
