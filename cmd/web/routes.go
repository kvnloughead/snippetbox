package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

/*
Returns a servemux that serves files from ./ui/static and contains the following routes:
  - GET  /										display the home page
  - GET  /snippet/view/:id    display a specific snippet
  - GET  /snippet/create      display form to create snippets
  - POST /snippet/create      create a new snippet
  - GET  /static/*filepath    serve a static file
*/
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Use our app.notFound method instead of httprouter's built-in 404 handler.
	router.NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.notFound(w)
		})

	// Serve static files out of ./ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		http.StripPrefix("/static", fileServer))

	// Routes
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.viewSnippet)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.createSnippet)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.createSnippetPost)

	// Initialize chain of standard pre-request middlewares.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
