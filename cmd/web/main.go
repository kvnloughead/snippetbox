package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Serve files out of ./ui/static directory. The path given should be relative
	// to the project's root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Register file server as handler for all URL paths that start with /static/.
	// Must strip the /static prefix before the request hits the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
