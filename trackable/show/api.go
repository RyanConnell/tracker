package show

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"tracker/server/auth"
	"tracker/server/page"

	"github.com/gorilla/mux"
)

// API implemnts server.API
type API struct {
	name    string
	handler Handler
}

var api *API

func (a *API) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(fmt.Sprintf("/%s/", subdomain), a.defaultRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/{id:[0-9]+}", subdomain), a.getRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/list", subdomain), a.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/list/{type:[a-z]*}", subdomain), a.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/get/schedule/{start:[0-9-]+}/{end:[0-9-]+}", subdomain),
		a.scheduleRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/request/addShow", subdomain), a.addShowRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/track", subdomain), a.trackShow)

	http.Handle(fmt.Sprintf("/%s/", subdomain), rtr)
}

func (a *API) Init() error {
	fmt.Println("Show APi Initialised")
	a.handler = Handler{}
	a.handler.Init()
	return nil
}

func (a *API) defaultRequest(w http.ResponseWriter, r *http.Request) {
	p := page.Page{[]byte("Show API landing page - Perhaps serve a README here?")}
	p.ServePage(w)
}

func (a *API) getRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		serveError(err, w)
		return
	}

	show, err := a.handler.Get(id)
	if err != nil {
		serveError(err, w)
		return
	}

	body, err := json.Marshal(show)
	if err != nil {
		serveError(err, w)
		return
	}
	p := page.Page{body}
	p.ServePage(w)
}

func (a *API) listRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	listType, ok := params["type"]
	if !ok {
		listType = "all"
	}

	user, err := auth.CurrentUser(r)
	if err != nil {
		serveRequestResult(false, err, w)
		return
	}
	fmt.Printf("[Backend] [Show] listRequest: CurrentUser=%v\n", user)

	list, err := a.handler.GetList(listType, user)
	if err != nil {
		serveError(err, w)
		return
	}

	body, err := json.Marshal(list)
	if err != nil {
		serveError(err, w)
		return
	}
	p := page.Page{body}
	p.ServePage(w)
}

func (a *API) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	schedule, err := a.handler.GetSchedule(params["start"], params["end"])
	if err != nil {
		serveError(err, w)
	}
	body, err := json.Marshal(schedule)
	if err != nil {
		serveError(err, w)
		return
	}
	p := page.Page{body}
	p.ServePage(w)
}

func (a *API) addShowRequest(w http.ResponseWriter, r *http.Request) {
	user, err := auth.CurrentUser(r)
	if err != nil {
		serveRequestResult(false, err, w)
		return
	}
	fmt.Printf("[Backend] [Show] addShowRequest: CurrentUser=%v\n", user)
	title := r.FormValue("title")
	wikipedia := r.FormValue("wikipedia")
	trailer := r.FormValue("trailer")
	coverImg := r.FormValue("img")

	success, err := a.handler.RequestShow(user, title, wikipedia, trailer, coverImg)
	serveRequestResult(success, err, w)
}

func (a *API) trackShow(w http.ResponseWriter, r *http.Request) {
	user, err := auth.CurrentUser(r)
	if err != nil {
		serveRequestResult(false, err, w)
		return
	}
	fmt.Printf("[Backend] [Show] trackShowRequest: CurrentUser=%v\n", user)
	show_id := r.FormValue("show_id")
	state := r.FormValue("state")

	success, err := a.handler.TrackShow(user, show_id, state)
	serveRequestResult(success, err, w)
}

func serveError(err error, w http.ResponseWriter) {
	p := page.Page{[]byte(fmt.Sprintf("Error occured: %v", err.Error()))}
	p.ServePage(w)
}

func serveRequestResult(success bool, err error, w http.ResponseWriter) {
	result := struct {
		Success bool
		Error   string
	}{success, ""}

	if err != nil {
		result.Error = err.Error()
	}

	body, err := json.Marshal(result)
	if err != nil {
		serveError(err, w)
		return
	}
	p := page.Page{body}
	p.ServePage(w)
}
