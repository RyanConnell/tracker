package show

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"tracker/server/host"
	"tracker/server/page"

	"github.com/gorilla/mux"
)

// API implemnts server.API
type API struct {
	name    string
	handler Handler
	host    *host.Host
}

func (a *API) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(fmt.Sprintf("/%s/", subdomain), a.defaultRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/{id:[0-9]+}", subdomain), a.getRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/list", subdomain), a.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/list/{type:[a-z]*}", subdomain), a.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/schedule/{start:[0-9-]+}/{end:[0-9-]+}", subdomain),
		a.scheduleRequest)

	http.Handle(fmt.Sprintf("/%s/", subdomain), rtr)
}

func (a *API) Init(*host.Host) error {
	fmt.Println("Show APi Initialised")
	a.handler.Init()
	return nil
}

func (a *API) defaultRequest(w http.ResponseWriter, r *http.Request) {
	p := page.Page{Body: []byte("Show API landing page - Perhaps serve a README here?")}
	p.ServePage(w)
}

func (a *API) getRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		serveError(err, w, r)
		return
	}

	show, err := a.handler.Get(id)
	if err != nil {
		serveError(err, w, r)
		return
	}

	body, err := json.Marshal(show)
	if err != nil {
		serveError(err, w, r)
		return
	}
	p := page.Page{Body: body}
	p.ServePage(w)
}

func (a *API) listRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	listType, ok := params["type"]
	if !ok {
		listType = "all"
	}

	list, err := a.handler.GetList(listType)
	if err != nil {
		serveError(err, w, r)
		return
	}
	body, err := json.Marshal(list)
	if err != nil {
		serveError(err, w, r)
		return
	}
	p := page.Page{Body: body}
	p.ServePage(w)
}

func (a *API) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	schedule, err := a.handler.GetSchedule(params["start"], params["end"])
	if err != nil {
		serveError(err, w, r)
	}
	body, err := json.Marshal(schedule)
	if err != nil {
		serveError(err, w, r)
		return
	}
	p := page.Page{Body: body}
	p.ServePage(w)
}

func serveError(err error, w http.ResponseWriter, r *http.Request) {
	p := page.Page{Body: []byte(fmt.Sprintf("Error occured: %v", err.Error()))}
	p.ServePage(w)
}
