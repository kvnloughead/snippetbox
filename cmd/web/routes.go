package main

import (
	"net/http"

	"github.com/justinas/alice"
)

/*
Returns a servemux that serves files from ./ui/static and contains the following routes:
  - GET  /
  - GET  /snippet/view
  - POST /snippet/create
*/
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Serve static files out of ./ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.viewSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Initialize standard set of pre-request middlewares.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
