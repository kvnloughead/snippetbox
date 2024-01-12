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

	// Middleware chain for dynamic routes only (not static files). Currently this
	// chain only includes our session manager middleware.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Dynamic routes are wrapped in our dynamic middleware. Note that since
	// ThenFunc returns an http.Handler, we need to use router.Handler instead of
	// router.HandlerFunc.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.viewSnippet))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.createSnippet))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.createSnippetPost))

	// Initialize chain of standard pre-request middlewares.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
