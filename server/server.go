package server

import (
	"fmt"
	"net/http"

	"tracker/server/auth"
	"tracker/trackable/common"
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
	Init(host *common.Host) error
}

func Launch(host *common.Host, settings map[string]string) {
	fmt.Printf("Starting server on %s\n", host.Address())

	// Register the Show API.
	show_api := &show.API{}
	show_api.RegisterHandlers("api/show")
	apis = append(apis, show_api)

	show_frontend := &show.Frontend{}
	show_frontend.RegisterHandlers("show")
	frontends = append(frontends, show_frontend)

	// Register the auth frontend.
	auth_api := &auth.API{}
	auth_api.Init(host, settings)

	// Register public files such as CSS, JS, and Images.
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("templates/public"))))

	for _, api := range apis {
		api.Init()
	}

	for _, frontend := range frontends {
		frontend.Init(host)
	}

	fmt.Printf("%d frontends registered\n", len(frontends))
	fmt.Printf("%d apis registered\n", len(apis))

	err := http.ListenAndServe(fmt.Sprintf(":%d", host.Port()), nil)
	fmt.Printf("Error encountered; %v", err)
}
