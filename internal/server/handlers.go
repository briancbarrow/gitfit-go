package server

import (
	"errors"
	"net/http"

	"github.com/briancbarrow/gitfit-go/internal/models"
	"github.com/briancbarrow/gitfit-go/internal/validator"
	"github.com/briancbarrow/gitfit-go/web/ui/pages"
)

type templateData struct {
	Form any
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	form := &pages.UserRegisterForm{}

	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.NotBlankFull(form.Name, "name")
	form.NotBlankFull(form.Email, "email")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.NotBlankFull(form.Password, "password")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field mus be at least 8 characters long")
	// form.check
	if !form.Valid() {
		var data templateData
		data.Form = form
		pages.Register(*form).Render(r.Context(), w)
		return
	}
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			var data templateData
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			pages.Register(*form).Render(r.Context(), w)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	form := &pages.UserLoginForm{}

	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.NotBlankFull(form.Email, "email")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.NotBlankFull(form.Password, "password")

	_, err = app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password do not match")
			var data templateData
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			pages.LoginPage(*form).Render(r.Context(), w)
			return
		}
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
