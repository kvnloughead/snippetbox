package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/kvnloughead/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// Gather templates into a string slice.
	// Note that the filepath used is relative to the root of the project.
	// This means that you need to run the server from the root directory.
	// files := []string{
	// 	// base template must come first
	// 	"./ui/html/base.tmpl", "./ui/html/pages/home.tmpl", "./ui/html/partials/nav.tmpl",
	// }

	// // Parse the template set or return a 500 error.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// // Write template content to response body. We need to specify our base
	// // template as the second argument.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
}

// View page for the snippet with the given ID.
// If there's no matching snippet a 404 NotFound response is sent.
func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template set or return a 500 error.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	/* Pass the the template data as final parameter to ts.ExecuteTemplate. We create a struct to contain all data necessary in the template, for now this is just the snippet.

	The data will be represented by the `.` character inside the templates. The type of the dot will be the type of data that's passed to it.

	You can then access data within the templates via a familiar sort of syntax. For instance `.Title` expands to the snippets title. */
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
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
