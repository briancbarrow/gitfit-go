package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/briancbarrow/gitfit-go/web/ui"
	"github.com/briancbarrow/gitfit-go/web/ui/pages"
	"github.com/julienschmidt/httprouter"
)

func (app *application) RegisterRoutes() http.Handler {
	r := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Handler(http.MethodGet, "/css/*filepath", fileServer)
	r.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// r.HandlerFunc(http.MethodGet, "/", app.helloWorldHandler)
	r.Handler(http.MethodGet, "/", templ.Handler(pages.Home()))
	r.Handler(http.MethodGet, "/login", templ.Handler(pages.LoginPage(pages.UserLoginForm{})))
	r.HandlerFunc(http.MethodPost, "/login", app.userLoginPost)
	r.Handler(http.MethodGet, "/register", templ.Handler(pages.Register(pages.UserRegisterForm{})))
	r.HandlerFunc(http.MethodPost, "/register", app.userRegisterPost)
	r.Handler(http.MethodGet, "/dashboard", templ.Handler(pages.Dashboard()))

	return r
}
