package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Gather templates into a string slice.
	// Note that the filepath used is relative to the root of the project.
	// This means that you need to run the server from the root directory.
	files := []string{
		// base template must come first
		"./ui/html/base.tmpl", "./ui/html/pages/home.tmpl", "./ui/html/partials/nav.tmpl",
	}

	// Parse the template set or return a 500 error.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Write template content to response body. We need to specify our base
	// template as the second argument.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

/*
Inserts a new record into the database. If successful, redirects the user to
the corresponding page with a 303 status code.
*/
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Temporary dummy data.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Insert new record or respond with a server error.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect to page containing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
