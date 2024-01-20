package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/kvnloughead/snippetbox/internal/models/mocks"
)

// Regex to capture CSRF from input on signup page.
var csrfTokenRX = regexp.MustCompile(`<input type=["']hidden["'] name=["']csrf_token["'] value=["'](.+)["'] ?/?>`)

func extractCSRFToken(t *testing.T, body string) string {
	// FindStringSubmatch returns an array with the matched pattern at index 0 and any captured values at subsequent indexes.
	matches := csrfTokenRX.FindStringSubmatch(body)
	t.Logf(matches[0])
	if len(matches) < 2 {
		t.Fatal("no csrf found in body")
	}

	// It's necessary to unescape the token, because html/template library will automatically escape special characters that may be included in the token, such as '+'.
	return html.UnescapeString((string(matches[1])))
}

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

	body := readBodyAsString(t, response)

	return response.StatusCode, response.Header, body
}

// Reads the body of an HTTP response and returns it as a string
// whitespace-trimmed string. Closes the response's body when returning.
//
// Errors occuring when reading the body are logged to test output (when they
// are run with the -v flag).
func readBodyAsString(t *testing.T, r *http.Response) string {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes.TrimSpace(body))
}
