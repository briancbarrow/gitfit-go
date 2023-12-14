package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/briancbarrow/gitfit-go/cmd/web/ui/pages"
	"github.com/briancbarrow/gitfit-go/internal/models"
	"github.com/briancbarrow/gitfit-go/internal/validator"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

type templateData struct {
	Form any
}

func (app *application) userRegisterPostStytch(w http.ResponseWriter, r *http.Request) {
	form := &pages.UserRegisterForm{}

	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	params := &passwords.CreateParams{
		Email:    form.Email,
		Password: form.Password,
	}

	_, err = app.stytchAPIClient.Passwords.Create(context.Background(), params)
	if err != nil {
		fmt.Printf("ERROR IS TYPE: %T", err)
		if stytchErr, ok := err.(stytcherror.Error); ok {
			// Now stytchErr is of type *stytcherror.Error and you can access its fields
			if stytchErr.ErrorType == "weak_password" {
				form.AddFieldError("password", string(stytchErr.ErrorMessage))

				strengthCheckParams := &passwords.StrengthCheckParams{
					Email:    form.Email,
					Password: form.Password,
				}
				strengthCheck, err := app.stytchAPIClient.Passwords.StrengthCheck(context.Background(), strengthCheckParams)
				if err != nil {
					fmt.Print("Error getting strength Check")
				}
				for _, suggestion := range strengthCheck.Feedback.Suggestions {
					form.AddNonFieldError(suggestion)
				}
			} else {
				form.AddNonFieldError(string(stytchErr.ErrorMessage))
			}

			var data templateData
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			pages.Register(*form).Render(r.Context(), w)
		} else {
			// err is not of type *stytcherror.Error
			fmt.Printf("\n\nOther error: %v\n", err)
		}
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
