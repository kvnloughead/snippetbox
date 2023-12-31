package main

import "net/http"

/*
Returns a servemux that serves files from ./ui/static and contains the following routes:
  - GET  /
  - GET  /snippet/view
  - POST /snippet/create
*/
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files out of ./ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
