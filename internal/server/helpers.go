package server

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/briancbarrow/gitfit-go/cmd/web/ui"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"github.com/stytchauth/stytch-go/v11/stytch/consumer/sessions"
)

func (app *application) newTemplateData(r *http.Request) ui.TemplateData {
	return ui.TemplateData{
		Toast:           app.sessionManager.PopString(r.Context(), "toast"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	var stytchSessionToken string
	if _, ok := app.sessionManager.Get(r.Context(), "stytchSessionToken").(string); ok {
		stytchSessionToken = app.sessionManager.Get(r.Context(), "stytchSessionToken").(string)
	} else {
		stytchSessionToken = ""
	}
	resp, err := app.stytchAPIClient.Sessions.Authenticate(r.Context(), &sessions.AuthenticateParams{
		SessionToken: stytchSessionToken,
	})
	if err != nil {
		// TODO: Need to handle if stych is down and throwing errors rather than just returning false
		return false
	}
	if resp.StatusCode >= 200 {
		return true
	} else {
		return false
	}
}
