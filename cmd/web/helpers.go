package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

/*
Writes an Error level log entry, including the request method and uri, and
returns a 500 Internal Server Error.
*/
func (app *application) serverError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = debug.Stack()
	)

	// Log error with stack trace.
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	// Send http error response to client.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	// If in development, output stack trace in human readable form to stderr.
	if os.Getenv("GO_ENV") != "production" {
		debug.PrintStack()
		app.logger.Error("IN DEV")
	}
}

/*
Sends a status code and corresponding description to the user. Uses
http.StatusText to generate the standard description.
*/
func (app *application) clientError(
	w http.ResponseWriter,
	status int,
) {
	http.Error(w, http.StatusText(status), status)
}

// A wrapper around clientError to easily send 404 Not Found responses.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Retrieves appropriate template from the cache, based on the page string.
// If no entry exists in the cache, a 500 server error is returned.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Write template to a buffer, instead of immediately to the ResponseWriter.
	// If there's an error, return a server error instead of a 200 response.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If template is written to buffer with no errors, we are safe to send a 200
	// response to the client.
	w.WriteHeader(status)
	buf.WriteTo(w) // write contents of buffer to http.ResponseWriter
}
