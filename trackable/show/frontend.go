package show

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"tracker/date"
	"tracker/server/auth"
	"tracker/server/host"
	"tracker/web"
	"tracker/web/templates"

	"github.com/gorilla/mux"
)

const DEVMODE = false

// Frontend implemnts server.Frontend
type Frontend struct {
	name       string
	host       *host.Host
	apiHost    *host.Host
	handler    Handler
	templates  *template.Template
	httpClient *http.Client
}

func (f *Frontend) RegisterHandlers(subdomain string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(fmt.Sprintf("/%s/", subdomain), f.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/request", subdomain), f.addShowRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/schedule", subdomain), f.scheduleRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/{type:[a-z]+}", subdomain), f.listRequest)
	rtr.HandleFunc(fmt.Sprintf("/%s/{id:[0-9]+}", subdomain), f.detailRequest)
	// TODO: Add slug based detail view

	http.Handle(fmt.Sprintf("/%s/", subdomain), rtr)
}

func (f *Frontend) Init(serverHost, apiHost *host.Host) (err error) {
	fmt.Println("Show Frontend Initialised")
	f.host = serverHost
	f.apiHost = apiHost
	f.httpClient = http.DefaultClient

	// Define all template functions
	funcMap := template.FuncMap{
		"mod":          templates.Mod,
		"doubleDigits": templates.DoubleDigits,
	}

	f.templates, err = template.New("").
		Funcs(funcMap).
		ParseFS(web.Templates, "**/**.html")

	if err != nil {
		return fmt.Errorf("unable to parse templates: %w", err)
	}

	log.Println("defined templates:", f.templates.DefinedTemplates())

	return nil
}

// Reload is only used for debugging/dev purposes. It will reinitialize the frontend each time it's
// called. This helps with development as we don't have to restart the server to see updates in
// the templates
func (f *Frontend) Reload() {
	if DEVMODE {
		_ = f.Init(f.host, f.apiHost)
	}
}

type ListRequestData struct {
	Title string

	ShowList
	User auth.User
}

func (f *Frontend) listRequest(w http.ResponseWriter, r *http.Request) {
	f.Reload()

	params := mux.Vars(r)
	listType, ok := params["type"]
	if !ok {
		listType = "all"
	}

	u := fmt.Sprintf("/api/show/get/list/%s", listType)

	var showList ShowList
	if err := f.get(r.Context(), u, &showList); err != nil {
		serveError(err, w, r)
		return
	}

	user, err := auth.CurrentUser(r)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	}

	data := ListRequestData{
		Title:    fmt.Sprintf("Show Tracker - %s", strings.Title(listType)),
		ShowList: showList,
		User:     user,
	}

	fmt.Printf("\tTemplate: User=%v\n", user)

	if err = f.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		serveError(err, w, r)
	}
}

type DetailsRequestData struct {
	Title string

	ShowFull
	User auth.User
}

func (f *Frontend) detailRequest(w http.ResponseWriter, r *http.Request) {
	f.Reload()

	params := mux.Vars(r)
	id := params["id"]

	u := fmt.Sprintf("/api/show/get/%s", id)

	var showDetails ShowFull
	if err := f.get(r.Context(), u, &showDetails); err != nil {
		serveError(err, w, r)
		return
	}

	user, err := auth.CurrentUser(r)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	}

	data := DetailsRequestData{
		Title:    fmt.Sprintf("Show Tracker - %s", showDetails.Name),
		ShowFull: showDetails,
		User:     user,
	}

	fmt.Printf("\tTemplate: User=%v\n", user)

	if err = f.templates.ExecuteTemplate(w, "detail.html", data); err != nil {
		serveError(err, w, r)
	}
}

type ScheduleRequestData struct {
	Title string

	Schedule
	User auth.User
}

func (f *Frontend) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	f.Reload()

	curDate := date.CurrentDate()
	startDate := curDate.Minus(7 + curDate.Weekday())
	endDate := startDate.Plus((7 * 7) - 1)

	u := fmt.Sprintf("/api/show/get/schedule/%s/%s", startDate, endDate)

	var schedule Schedule
	if err := f.get(r.Context(), u, &schedule); err != nil {
		serveError(err, w, r)
		return
	}

	user, err := auth.CurrentUser(r)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	}

	data := ScheduleRequestData{
		Title:    "Show Tracker - Schedule",
		Schedule: schedule,
		User:     user,
	}

	fmt.Printf("\tTemplate: User=%v\n", user)

	if err = f.templates.ExecuteTemplate(w, "schedule.html", data); err != nil {
		serveError(err, w, r)
	}
}

func (f *Frontend) loginRequest(w http.ResponseWriter, r *http.Request) {
	f.Reload()

	err := f.templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		serveError(err, w, r)
	}
}

type AddShowRequestData struct {
	Title string

	User auth.User
}

func (f *Frontend) addShowRequest(w http.ResponseWriter, r *http.Request) {
	f.Reload()

	user, err := auth.CurrentUser(r)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	}

	data := AddShowRequestData{
		Title: "Show Tracker - Add Tracking",

		User: user,
	}

	if err = f.templates.ExecuteTemplate(w, "add_show.html", data); err != nil {
		serveError(err, w, r)
	}
}

// get the given URL and unmarshal data into res. The url specified must be
// prefixed with a /
func (f *Frontend) get(ctx context.Context, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", f.apiHost.Address(), url),
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	res, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to request: %w", err)
	}
	// accept 200 and 300s, ignore 100s
	if c := res.StatusCode / 100; c > 3 {
		return errors.New("bad response")
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}

	return nil
}
