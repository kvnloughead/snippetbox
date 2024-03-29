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

//
// Basic Handlers (ping, home, about)
//

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Displays home page in response to GET /. If we were using http.ServeMux we
// would have to check the URL, but with httprouter.Router, "/" is exclusive.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// Displays about page in response to GET /about.
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about.tmpl", data)
}

//
// Snippet handlers
//

// Struct containing form fields for the /snippet/create form.
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // "-" tells formDecoder to ignore the field
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{Expires: 365}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

/*
Inserts a new record into the database. If successful, redirects the user to
the corresponding page with a 303 status code.

If one or more fields are invalid, the form is rendered again with a 422 status
code, displaying the appropriate error messages.

If we were using http.ServeMux, we would have to check the method in this handler.
*/
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create an instance of our form struct and decode it with decodePostForm.
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
	app.sessionManager.Put(r.Context(), string(flash), "Snippet successfully created!")

	// Redirect to page containing the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

//
// User handlers

// Struct containing form fields for the /user/signup form.
type userSignupForm struct {
	Name                string     `form:"name"`
	Email               string     `form:"email"`
	Password            string     `form:"password"`
	validator.Validator `form:"-"` // "-" tells formDecoder to ignore the field
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field can't be blank.")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field can't be blank.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field can't be blank.")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "That email is already in use.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), string(flash), "Your signup was successful, please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Struct containing form fields for the /user/login form.
type userLoginForm struct {
	Email               string     `form:"email"`
	Password            string     `form:"password"`
	validator.Validator `form:"-"` // "-" tells formDecoder to ignore the field
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field can't be blank.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field can't be blank.")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Try to authenticate user. If the user's credentials are invalid, the
	// login page is re-rendered with a non-field error.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnauthorized, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
	}

	// When authentication state or privilege levels change, the session ID should
	// be changed, via the RenewToken method.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Add user ID to session data to indicate their logged in status.
	app.sessionManager.Put(r.Context(), string(authenticatedUserID), id)

	app.sessionManager.Put(r.Context(), string(flash), "Login successful.")

	dest := app.sessionManager.PopString(r.Context(), string(redirectAfterLogin))
	if dest != "" {
		http.Redirect(w, r, dest, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// When authentication state or privilege levels change, the session ID should
	// be changed, via the RenewToken method.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// "Logout" user by removing the authenticatedUserID, and flash success.
	app.sessionManager.Remove(r.Context(), string(authenticatedUserID))
	app.sessionManager.Put(r.Context(), string(flash), "You have succesfully logged out.")

	// Redirect to home.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//
// Account handlers
//

type accountPasswordUpdateForm struct {
	CurrentPassword     string     `form:"currentPassword"`
	NewPassword         string     `form:"newPassword"`
	ConfirmPassword     string     `form:"confirmPassword"`
	validator.Validator `form:"-"` // "-" tells formDecoder to ignore the field
}

func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}
	app.render(w, r, http.StatusOK, "password.tmpl", data)
}

func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field can't be blank.")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirmPassword", "This field can't be blank.")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "Password must be at least 8 characters.")
	form.CheckField(validator.MinChars(form.ConfirmPassword, 8), "confirmPassword", "Password must be at least 8 characters.")
	form.CheckField(form.NewPassword == form.ConfirmPassword, "confirmPassword", "Passwords don't match.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "password.tmpl", data)
		return
	}

	// Get ID from session data to retrieve user's data.
	id := app.sessionManager.GetInt(r.Context(), string(authenticatedUserID))
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Verify that user entered the correct password.
	_, err = app.users.Authenticate(user.Email, form.CurrentPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Password is incorrect.")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnauthorized, "password.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// If validated and authenticated, update the password.
	err = app.users.PasswordUpdate(id, form.NewPassword)
	if err != nil {
		form.AddNonFieldError("Failed to update password.")
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), string(flash), "Password successfully updated.")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

// Displays account page in response to GET /account/view.
func (app *application) accountView(w http.ResponseWriter, r *http.Request) {

	id := app.sessionManager.GetInt(r.Context(), string(authenticatedUserID))
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, r, http.StatusOK, "account.tmpl", data)
}
