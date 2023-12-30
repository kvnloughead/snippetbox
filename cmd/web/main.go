package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", "4000", "HTTP Network Address")
	flag.Parse()

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

	log.Printf("starting server on :%s", *addr)

	err := http.ListenAndServe(":"+*addr, mux)
	log.Fatal(err)
}
