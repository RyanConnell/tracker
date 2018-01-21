package server

import (
	"fmt"
	"net/http"

	"tracker/trackable/show"
)

var apis []API = make([]API, 0)
var frontends []Frontend = make([]Frontend, 0)

type API interface {
	RegisterHandlers(subdomain string)
	Init() error
}

type Frontend interface {
	RegisterHandlers(subdomain string)
	Init() error
}

func Launch(port int) {
	fmt.Printf("Starting server on port %d\n", port)

	// Register the Show API.
	show_api := &show.API{}
	show_api.RegisterHandlers("api/show")
	apis = append(apis, show_api)

	show_frontend := &show.Frontend{}
	show_frontend.RegisterHandlers("show")
	frontends = append(frontends, show_frontend)

	// Register public files such as CSS, JS, and Images.
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("templates/public"))))

	for _, api := range apis {
		api.Init()
	}

	for _, frontend := range frontends {
		frontend.Init()
	}

	fmt.Printf("%d frontends registered\n", len(frontends))
	fmt.Printf("%d apis registered\n", len(apis))

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
