package main

import (
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
