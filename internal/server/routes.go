package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/briancbarrow/gitfit-go/cmd/web/ui"
	"github.com/briancbarrow/gitfit-go/cmd/web/ui/pages"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) RegisterRoutes() http.Handler {
	r := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Handler(http.MethodGet, "/css/*filepath", fileServer)
	r.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	r.Handler(http.MethodGet, "/", dynamic.Then(templ.Handler(pages.Home(ui.TemplateData{}))))

	notProtected := dynamic.Append(app.redirectIfAuthenticated)
	r.Handler(http.MethodGet, "/login", notProtected.ThenFunc(app.loginGet))
	r.Handler(http.MethodGet, "/register", notProtected.ThenFunc(app.registerGet))

	r.Handler(http.MethodPost, "/login", dynamic.ThenFunc(app.userLoginPostStytch))
	r.Handler(http.MethodPost, "/register", dynamic.ThenFunc(app.userRegisterPostStytch))

	protected := dynamic.Append(app.requireAuthentication)
	r.Handler(http.MethodGet, "/dashboard", protected.ThenFunc(app.dashboardGet))

	standard := alice.New(app.recoverPanic, app.logRequest, ui.SecureHeaders)

	return standard.Then(r)
}
