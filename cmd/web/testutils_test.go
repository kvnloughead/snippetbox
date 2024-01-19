package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Returns an application struct with mocked dependencies for testing.
// Currently, logger is the only included dependency, because some middlewares depend on it.
func newTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// Custom server struct for testing.
type testServer struct {
	*httptest.Server
}

// Create HTTPS server serving the specified handler, for testing on random port.
func newTestServer(t *testing.T, h http.Handler) testServer {
	ts := httptest.NewTLSServer(h)
	return testServer{ts}
}

// Creates a test client to send requests to our ts. And sends a GET request to the given endpoint. Returns the status code, headers and body of the response.
func (ts *testServer) get(t *testing.T, endpoint string) (int, http.Header, string) {
	response, err := ts.Client().Get(ts.URL + endpoint)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return response.StatusCode, response.Header, string(body)
}
