package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/kvnloughead/snippetbox/internal"
)

func TestPing(t *testing.T) {
	// ResponseRecorder will write our responses without sending them over HTTP.
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		// Marks test as faild, logs the error, and stops the test execution.
		t.Fatal(err)
	}

	// Call ping handler with ResponseRecorder and our new HTTP request.
	ping(rr, r)

	// rr.Result() returns the response.
	response := rr.Result()

	// Verify the status code.
	assert.Equal(t, response.StatusCode, http.StatusOK)

	// Verify the body.
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		// Marks test as faild, logs the error, and stops the test execution.
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
