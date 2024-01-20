package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/kvnloughead/snippetbox/internal/models/mocks"
)

// Returns an application struct with mocked dependencies for testing.
// Currently, logger is the only included dependency, because some middlewares depend on it.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// Same settings as production, but without the mysql store. Without a store, sessions will be stored in-memory.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// Custom server struct for testing.
type testServer struct {
	*httptest.Server
}

// Create HTTPS server serving the specified handler, for testing on random port.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize a test server and cookie jar.
	ts := httptest.NewTLSServer(h)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add cookie jar to test server client. Now response cookies will be stored
	// and sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following. This function will be recalled whenever a 3xx response is received by the client. Return ErrUseLastResponse forces client to immediately return the received response.
	ts.Client().CheckRedirect =
		func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}

	return &testServer{ts}
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
