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

	// r.HandlerFunc(http.MethodGet, "/", app.helloWorldHandler)
	r.Handler(http.MethodGet, "/", dynamic.Then(templ.Handler(pages.Home())))
	r.Handler(http.MethodGet, "/login", dynamic.Then(templ.Handler(pages.LoginPage(pages.UserLoginForm{}))))
	r.Handler(http.MethodPost, "/login", dynamic.ThenFunc(app.userLoginPostStytch))
	r.Handler(http.MethodGet, "/register", dynamic.Then(templ.Handler(pages.Register(pages.UserRegisterForm{}))))
	r.Handler(http.MethodPost, "/register", dynamic.ThenFunc(app.userRegisterPostStytch))
	r.Handler(http.MethodGet, "/dashboard", dynamic.Then(templ.Handler(pages.Dashboard())))

	standard := alice.New(app.recoverPanic, app.logRequest, ui.SecureHeaders)

	return standard.Then(r)
}
