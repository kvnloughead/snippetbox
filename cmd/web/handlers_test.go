package main

import (
	"net/http"
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
