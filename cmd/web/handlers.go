package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kvnloughead/snippetbox/internal/models"
)

// Displays home page in response to GET /. If we were using http.ServeMux we
// would have to check the URL, but with httprouter.Router, "/" is exclusive.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", templateData)
}

// View page for the snippet with the given ID.
// If there's no matching snippet a 404 NotFound response is sent.
func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	// Params are stored by httprouter in the request context.
	params := httprouter.ParamsFromContext(r.Context())

	// Once parsed, params are available by params.ByName().
	id, err := strconv.Atoi(params.ByName("id"))
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

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", templateData)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl", templateData)
}

/*
Inserts a new record into the database. If successful, redirects the user to
the corresponding page with a 303 status code.

If we were using http.ServeMux, we would have to check the method in this handler.
*/
func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {

	// r.ParseForm() populates r.Form and r.PostForm and validates response body.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get the values from the form, converting 'expires' to an 'int'.
	// PostForm.Get() returns only the first value for a given input, (by name).
	// To access all values (for instance, of several checkboxes), you can access
	// the underlying r.PostForm map, which has type map[string][]string.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Insert new record or respond with a server error.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect to page containing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
