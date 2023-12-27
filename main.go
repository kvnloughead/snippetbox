package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

/*
Handler for GET /.

Renders home page. If request URL doesn't exactly match "/", sends
a 404 NotFound response.
*/
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

/*
Handler for GET /snippet/view?id=

If ID isn't a positive integer, returns a 400 BadRequest.
*/
func viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Viewing snippet number %d\n", id)
}

/*
Handler for POST /snippet/create.

If request method isn't POST, returns 405 response.
*/
func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Creating snippet..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
