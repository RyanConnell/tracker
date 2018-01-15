package server

import (
	"fmt"
	"net/http"

	"tracker/trackable/show"
)

var apis []API = make([]API, 0)

type API interface {
	RegisterHandlers(subdomain string)
	Init() error
}

func Launch(port int) {
	fmt.Printf("Starting server on port %d\n", port)

	// Register the Show API.
	show_api := &show.API{}
	show_api.RegisterHandlers("show")
	apis = append(apis, show_api)

	for _, api := range apis {
		api.Init()
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Printf("Error encountered; %v", err)
}

/*
func getRequest(w http.ResponseWriter, r *http.Request) {
	show := &show.Show{
		Name:     "Arrow",
		Finished: false,
	}

	body, err := json.Marshal(show)
	if err != nil {
		renderError(fmt.Sprintf("%v", err), w)
	}

	page := Page{body}
	page.servePage(w)
}
*/
