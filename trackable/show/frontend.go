package show

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"tracker/templates"
	"tracker/trackable/common"

	"github.com/gorilla/mux"
)

const DEVMODE = true

// Frontend implemnts server.Frontend
type Frontend struct {
	name      string
	handler   Handler
	templates *template.Template
}

var frontend *Frontend

func (f *Frontend) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(fmt.Sprintf("/%s/", subdomain), f.indexRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/{id:[0-9]+}", subdomain), f.detailRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/schedule", subdomain), f.scheduleRequest)

	http.Handle(fmt.Sprintf("/%s/", subdomain), rtr)
}

func (f *Frontend) Init() error {
	fmt.Println("Show Frontend Initialised")
	frontend = f

	// Define all template functions
	funcMap := template.FuncMap{
		"mod":          templates.Mod,
		"doubleDigits": templates.DoubleDigits,
	}

	// Load all templates
	f.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob(
		"templates/shows/*.html"))

	return nil
}

func (f *Frontend) indexRequest(w http.ResponseWriter, r *http.Request) {
	if DEVMODE {
		f.Init()
	}

	apiURL := "http://localhost:8080/api/show/"
	resp, err := http.Get(fmt.Sprintf("%sget/list", apiURL))
	if err != nil {
		fmt.Println(err)
		return
	}

	decode := json.NewDecoder(resp.Body)
	var jsonRep ShowList
	decode.Decode(&jsonRep)

	err = f.templates.ExecuteTemplate(w, "index.html", jsonRep)
	if err != nil {
		serveError(err, w, r)
	}
}

func (f *Frontend) detailRequest(w http.ResponseWriter, r *http.Request) {
	if DEVMODE {
		f.Init()
	}

	params := mux.Vars(r)
	id := params["id"]

	apiURL := "http://localhost:8080/api/show"
	resp, err := http.Get(fmt.Sprintf("%s/get/%s", apiURL, id))
	if err != nil {
		serveError(err, w, r)
		return
	}

	decode := json.NewDecoder(resp.Body)
	var jsonRep ShowFull
	decode.Decode(&jsonRep)

	err = f.templates.ExecuteTemplate(w, "detail.html", jsonRep)
	if err != nil {
		serveError(err, w, r)
	}
}

func (f *Frontend) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	if DEVMODE {
		f.Init()
	}

	curDate := common.CurrentDate()
	startDate := curDate.Minus(7 + curDate.Weekday())
	endDate := startDate.Plus((7 * 7) - 1)

	apiURL := "http://localhost:8080/api/show/"
	url := fmt.Sprintf("%sget/schedule/%s/%s", apiURL, startDate, endDate)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	decode := json.NewDecoder(resp.Body)
	var jsonRep Schedule
	decode.Decode(&jsonRep)

	err = f.templates.ExecuteTemplate(w, "schedule.html", jsonRep)
	if err != nil {
		serveError(err, w, r)
	}
}
