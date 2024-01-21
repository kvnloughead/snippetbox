package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/kvnloughead/snippetbox/ui"
)

/*
Returns a servemux that serves files from ./ui/static and contains the following routes:

Static unprotected routes
  - GET  /static/*filepath    serve a static file

Dynamic unprotected routes:
  - GET  /										display the home page
  - GET  /about								display the about page
  - GET  /ping 							  responses with 200 OK
  - GET  /snippet/view/:id    display a specific snippet
  - GET  /user/signup					display the signup form
  - POST /user/signup					create a new user
  - GET  /user/login					display the login form
  - POST /user/login					authenticate and login a user

Protected routes (only available to authenticated users):
  - POST /user/logout         logout the user
  - GET  /snippet/create      display form to create snippets
  - POST /snippet/create      create a new snippet
*/
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Use our app.notFound method instead of httprouter's built-in 404 handler.
	router.NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.notFound(w)
		})

	// Serve static files out of embedded filesystem ui.Files.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(
		http.MethodGet,
		"/static/*filepath",
		fileServer,
	)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Middleware chain for dynamic routes only (not static files).
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Dynamic routes are wrapped in our dynamic middleware. Note that since
	// ThenFunc returns an http.Handler, we need to use router.Handler instead of
	// router.HandlerFunc.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Middleware chain for protected routes. Includes all middleware from dynamic
	// chain, as well as app.requireAuthentication.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// Initialize chain of standard pre-request middlewares.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
