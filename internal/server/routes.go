package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/briancbarrow/gitfit-go/ui"
	"github.com/julienschmidt/httprouter"
)

func (app *application) RegisterRoutes() http.Handler {
	component := ui.HelloTempl("Brian")
	r := httprouter.New()
	r.HandlerFunc(http.MethodGet, "/", app.helloWorldHandler)
	r.Handler(http.MethodGet, "/web", templ.Handler(component))

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
