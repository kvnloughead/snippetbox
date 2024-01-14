package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kvnloughead/snippetbox/internal/models"
	"github.com/kvnloughead/snippetbox/internal/validator"
)

// A struct for passing data to our templateData struct. Contains all form
// fields, plus an embedded validator struct. The tags instruct our application wide form decoder on how to map struct fields to markup.
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // "-" tells formDecoder to ignore the field
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
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display signup form")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "submit signup form")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display login form")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "submit login form")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout user")
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = snippetCreateForm{Expires: 365}
	app.render(w, r, http.StatusOK, "create.tmpl", templateData)
}

/*
Inserts a new record into the database. If successful, redirects the user to
the corresponding page with a 303 status code.

If one or more fields are invalid, the form is rendered again with a 422 status
code, displaying the appropriate error messages.

If we were using http.ServeMux, we would have to check the method in this handler.
*/
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create an instance of our form struct and decode it with the app.decodePostForm.
	// This automatically parses the values passed as the second argument into the
	// corresponding struct fields, making appropriate data conversions.
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate all form fields.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field can't be blank.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This can't contain more than 100 characters.")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field can't be blank.")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365.")

	// If there are any validation errors, render the page again with the errors.
	if !form.Valid() {
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

	// Assign text to session data with the key "flash". The data is stored in the
	// request's context. If there is no current session, a new one will be created.
	// The flash is added to our template data via the newTemplateData function.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// Redirect to page containing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
