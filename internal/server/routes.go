package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/briancbarrow/gitfit-go/ui"
	"github.com/briancbarrow/gitfit-go/ui/pages"
	"github.com/julienschmidt/httprouter"
)

func (app *application) RegisterRoutes() http.Handler {
	r := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	r.Handler(http.MethodGet, "/css/*filepath", fileServer)
	r.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// r.HandlerFunc(http.MethodGet, "/", app.helloWorldHandler)
	r.Handler(http.MethodGet, "/", templ.Handler(pages.Home()))
	r.Handler(http.MethodGet, "/login", templ.Handler(pages.LoginPage()))
	r.Handler(http.MethodGet, "/register", templ.Handler(pages.Register()))

	return r
}

func (app *application) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Write(jsonResp)
}
