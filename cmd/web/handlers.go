package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/kvnloughead/snippetbox/internal/models"
)

// A struct for passing data to our templateData struct. Contains all form
// fields, plus a map for potential error messages.
type createSnippetForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

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
	templateData.Form = createSnippetForm{Expires: 365}
	app.render(w, r, http.StatusOK, "create.tmpl", templateData)
}

/*
Inserts a new record into the database. If successful, redirects the user to
the corresponding page with a 303 status code.

If one or more fields are invalid, the form is rendered again with a 422 status code, displaying the appropriate error messages.

If we were using http.ServeMux, we would have to check the method in this handler.
*/
func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {

	// r.ParseForm() populates r.Form and r.PostForm and validates response body.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get the 'expires' value, converting it to an 'int'.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create form struct containing field values and a map for potential errors.
	form := createSnippetForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Validate title field. If invalid, an error is added to form.FieldErrors.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field can't be blank."
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field can't contain more than 100 characters."
	}

	// Validate content field. If invalid, an error is added to form.FieldErrors.
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field can't be blank."
	}

	// Validate expires field. If invalid, an error is added to form.FieldErrors.
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7, or 365."
	}

	// If there are any validation errors, render the page again with the errors.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Insert new record or respond with a server error.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect to page containing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
