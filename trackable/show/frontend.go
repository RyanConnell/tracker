package show

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Frontend implemnts server.Frontend
type Frontend struct {
	name    string
	handler Handler
}

var frontend *Frontend

func (_ *Frontend) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/show/", indexRequest)
	rtr.HandleFunc("/show/{id:[0-9]+}", detailRequest)

	http.Handle("/", rtr)
}

func (f *Frontend) Init() error {
	fmt.Println("Show Frontend Initialised")
	frontend = f
	return nil
}

func indexRequest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/shows/index.html", "templates/shows/navbar.html")

	if err != nil {
		fmt.Println(err)
		return
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

	tmpl.Execute(w, jsonRep)
}

func detailRequest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/shows/detail.html", "templates/shows/navbar.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	apiURL := "http://localhost:8080/api/show"
	resp, err := http.Get(fmt.Sprintf("%s/get/%s", apiURL, id))
	if err != nil {
		fmt.Println(err)
		return
	}

	decode := json.NewDecoder(resp.Body)
	var jsonRep Show
	decode.Decode(&jsonRep)

	tmpl.Execute(w, jsonRep)
}
