package main

import (
	"fmt"
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
