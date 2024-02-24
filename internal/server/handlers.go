package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/briancbarrow/gitfit-go/cmd/web/ui"
	"github.com/briancbarrow/gitfit-go/cmd/web/ui/pages"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/passwords"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/users"
	"github.com/stytchauth/stytch-go/v11/stytch/stytcherror"
)

func (app *application) dashboardGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	pages.Dashboard(data).Render(r.Context(), w)
}

func (app *application) loginGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	form := &ui.UserLoginForm{}
	pages.LoginPage(*form, data).Render(r.Context(), w)
}

func (app *application) registerGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	form := &ui.UserRegisterForm{}
	pages.Register(*form, data).Render(r.Context(), w)
}

func (app *application) userRegisterPostStytch(w http.ResponseWriter, r *http.Request) {
	form := &ui.UserRegisterForm{}
	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	params := &passwords.CreateParams{
		Name: &users.Name{
			FirstName: form.FirstName,
			LastName:  form.LastName,
		},
		Email:    form.Email,
		Password: form.Password,
	}

	response, err := app.stytchAPIClient.Passwords.Create(context.Background(), params)
	if err != nil {

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

			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			// err is not of type *stytcherror.Error
			form.AddNonFieldError("Error with registration. Please try again later.")
		}
		data := app.newTemplateData(r)
		pages.Register(*form, data).Render(r.Context(), w)
		return
	}
	// err = app.users.Insert(form.FirstName, form.LastName, form.Email, response.UserID)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *application) userLoginPostStytch(w http.ResponseWriter, r *http.Request) {
	form := &ui.UserLoginForm{}

	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	params := &passwords.AuthenticateParams{
		Email:                  form.Email,
		Password:               form.Password,
		SessionDurationMinutes: 60,
	}

	stytchResponse, err := app.stytchAPIClient.Passwords.Authenticate(context.Background(), params)
	if err != nil {
		fmt.Printf("ERROR IS TYPE: %T", err)
		if stytchErr, ok := err.(stytcherror.Error); ok {

			form.AddNonFieldError(string(stytchErr.ErrorMessage))

			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			form.AddNonFieldError("Error logging in. Please try again later.")
		}
		data := app.newTemplateData(r)
		pages.LoginPage(*form, data).Render(r.Context(), w)
		return
	}
	app.sessionManager.Put(r.Context(), "toast", "LOGGED IN Success")
	app.sessionManager.Put(r.Context(), "authenticatedUserID", stytchResponse.UserID)
	app.sessionManager.Put(r.Context(), "stytchSessionToken", stytchResponse.SessionToken)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) userLogoutPostStytch(w http.ResponseWriter, r *http.Request) {
	var stytchSessionToken string
	if _, ok := app.sessionManager.Get(r.Context(), "stytchSessionToken").(string); ok {
		stytchSessionToken = app.sessionManager.Get(r.Context(), "stytchSessionToken").(string)
	} else {
		stytchSessionToken = ""
	}
	_, err := app.stytchAPIClient.Sessions.Revoke(r.Context(), &sessions.RevokeParams{
		SessionToken: stytchSessionToken,
	})
	if err != nil {
		app.serverError(w, r, err)
		app.sessionManager.Remove(r.Context(), "stytchSessionToken")
	}
	app.sessionManager.Put(r.Context(), "toast", "You've been logged out successfully!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
