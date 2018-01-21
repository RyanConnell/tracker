package show

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"tracker/server/page"

	"github.com/gorilla/mux"
)

// API implemnts server.API
type API struct {
	name    string
	handler Handler
}

var api *API

func (_ *API) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(fmt.Sprintf("/%s/", subdomain), defaultRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/{id:[0-9]+}", subdomain), getRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/list", subdomain), listRequest)

	http.Handle(fmt.Sprintf("/%s/", subdomain), rtr)
}

func (a *API) Init() error {
	fmt.Println("Show APi Initialised")
	api = a
	api.handler = Handler{}
	api.handler.Init()
	return nil
}

func defaultRequest(w http.ResponseWriter, r *http.Request) {
	p := page.Page{[]byte("Show API landing page - Perhaps serve a README here?")}
	p.ServePage(w)
}

func getRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		serveError(err, w, r)
		return
	}

	show, err := api.handler.Get(id)
	if err != nil {
		serveError(err, w, r)
		return
	}

	body, err := json.Marshal(show)
	p := page.Page{body}

	p.ServePage(w)
}

func listRequest(w http.ResponseWriter, r *http.Request) {
	list := api.handler.GetList()
	body, err := json.Marshal(list)
	if err != nil {
		serveError(err, w, r)
		return
	}
	p := page.Page{body}
	p.ServePage(w)
}

func serveError(err error, w http.ResponseWriter, r *http.Request) {
	p := page.Page{[]byte(err.Error())}
	p.ServePage(w)
}
