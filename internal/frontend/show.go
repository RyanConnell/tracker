package frontend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"tracker/date"
	"tracker/internal/httpserver"
	"tracker/server/auth"
	"tracker/trackable/show"
	"tracker/web"
	"tracker/web/templates"
)

type ShowFrontend struct {
	funcs     template.FuncMap
	templates *template.Template

	apiAddr    string
	httpClient *http.Client
}

// NewFrontend creates a new frontend with default values, which can be
// overridden using the options. The components allow to extend the
// frontend with the given paths of the route.
func NewShow(apiAddr string, opts ...ShowOption) (*ShowFrontend, error) {
	f := &ShowFrontend{
		apiAddr:    apiAddr,
		httpClient: http.DefaultClient,
		funcs: template.FuncMap{
			"mod":          templates.Mod,
			"doubleDigits": templates.DoubleDigits,
		},
	}

	for _, opt := range opts {
		if err := opt(f); err != nil {
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	if f.templates == nil {
		if err := TemplateFS(web.Templates, "**/**.html")(f); err != nil {
			return nil, fmt.Errorf("unable to apply default template: %w", err)
		}
	}

	return f, nil
}

func (f *ShowFrontend) RegisterHandlers(r *mux.Router) {
	r.Path("/request").
		HandlerFunc(f.addShowRequest)
	r.Path("/schedule").
		HandlerFunc(f.scheduleRequest)
	r.Path("/{type:[a-z]+}").
		HandlerFunc(f.listRequest)
	r.Path("/{id:[0-9]+}").
		HandlerFunc(f.detailRequest)
	r.Path("/").
		HandlerFunc(f.listRequest)
	r.Path("").
		HandlerFunc(f.listRequest)
}

type ListRequestData struct {
	Title string

	show.ShowList
	User auth.User
}

func (f *ShowFrontend) listRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	listType, ok := params["type"]
	if !ok {
		listType = "all"
	}

	u := fmt.Sprintf("/api/show/get/list/%s", listType)

	var showList show.ShowList
	if err := f.get(r.Context(), u, &showList); err != nil {
		httpserver.ServeError(err, w)
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
		httpserver.ServeError(err, w)
	}
}

type DetailsRequestData struct {
	Title string

	show.ShowFull
	User auth.User
}

func (f *ShowFrontend) detailRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	u := fmt.Sprintf("/api/show/get/%s", id)

	var showDetails show.ShowFull
	if err := f.get(r.Context(), u, &showDetails); err != nil {
		httpserver.ServeError(err, w)
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
		httpserver.ServeError(err, w)
	}
}

type ScheduleRequestData struct {
	Title string

	show.Schedule
	User auth.User
}

func (f *ShowFrontend) scheduleRequest(w http.ResponseWriter, r *http.Request) {
	curDate := date.CurrentDate()
	startDate := curDate.Minus(7 + curDate.Weekday())
	endDate := startDate.Plus((7 * 7) - 1)

	u := fmt.Sprintf("/api/show/get/schedule/%s/%s", startDate, endDate)

	var schedule show.Schedule
	if err := f.get(r.Context(), u, &schedule); err != nil {
		httpserver.ServeError(err, w)
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
		httpserver.ServeError(err, w)
	}
}

func (f *ShowFrontend) loginRequest(w http.ResponseWriter, r *http.Request) {
	err := f.templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		httpserver.ServeError(err, w)
	}
}

type AddShowRequestData struct {
	Title string

	User auth.User
}

func (f *ShowFrontend) addShowRequest(w http.ResponseWriter, r *http.Request) {
	user, err := auth.CurrentUser(r)
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
	}

	data := AddShowRequestData{
		Title: "Show Tracker - Add Tracking",

		User: user,
	}

	if err = f.templates.ExecuteTemplate(w, "add_show.html", data); err != nil {
		httpserver.ServeError(err, w)
	}
}

// get the given URL and unmarshal data into res. The url specified must be
// prefixed with a /
func (f *ShowFrontend) get(ctx context.Context, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", f.apiAddr, url),
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

// ShowOption allows to modify the frontend.
type ShowOption func(*ShowFrontend) error

// HTTPClient allows to override the default client used for requests.
func HTTPClient(c *http.Client) ShowOption {
	return func(f *ShowFrontend) error {
		f.httpClient = c
		return nil
	}
}

// TemplateFS allows to override the fs.FS used for loading templates
func TemplateFS(tfs fs.FS, pattern string) ShowOption {
	return func(f *ShowFrontend) (err error) {
		f.templates, err = template.New("").
			Funcs(f.funcs).
			ParseFS(tfs, pattern)
		if err != nil {
			return fmt.Errorf("unable to parse templates from provided fs: %w", err)
		}

		return nil
	}
}
