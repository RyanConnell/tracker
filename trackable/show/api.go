package show

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tracker/server/page"
)

// API implemnts server.API
type API struct {
	name    string
	handler Handler
}

var api *API

func (_ *API) RegisterHandlers(subdomain string) {
	http.HandleFunc(fmt.Sprintf("/%s/", subdomain), defaultRequest)
	http.HandleFunc(fmt.Sprintf("/%s/get/", subdomain), getRequest)
	http.HandleFunc(fmt.Sprintf("/%s/get/list", subdomain), listRequest)
}

func (a *API) Init() error {
	fmt.Println("Show APi Initialised")
	api = a
	api.handler = Handler{}
	api.handler.Init()
	return nil
}

func defaultRequest(w http.ResponseWriter, r *http.Request) {
	p := page.Page{[]byte("<h1>Shows</h1>")}
	p.ServePage(w)
}

func getRequest(w http.ResponseWriter, r *http.Request) {
	show, err := api.handler.Get(1)
	if err != nil {
		serveError(err, w, r)
		return
	}

	body, err := json.Marshal(show)
	p := page.Page{body}

	p.ServePage(w)
}

func listRequest(w http.ResponseWriter, r *http.Request) {
	list, err := api.handler.GetList(5, 0)
	if err != nil {
		serveError(err, w, r)
		return
	}
	body, err := json.Marshal(list)
	p := page.Page{body}
	p.ServePage(w)
}

func serveError(err error, w http.ResponseWriter, r *http.Request) {
	p := page.Page{[]byte(err.Error())}
	p.ServePage(w)
}
