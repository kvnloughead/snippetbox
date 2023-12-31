package main

import "net/http"

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
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
