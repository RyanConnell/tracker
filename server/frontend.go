package server

import (
	"net/http"

	"tracker/server/host"
)

type WebFrontend interface {
	Init(serverHost, apiHost *host.Host) error
	RegisterHandlers(subdomain string)
}

type Frontend struct {
	frontends  map[string]WebFrontend
	serverHost *host.Host
}

func NewFrontend(settings *Settings, frontends map[string]WebFrontend) (*Frontend, error) {
	serverHost := host.NewHost(settings.Hostname, settings.Port)
	apiHost := host.NewHost(settings.APIHostname, settings.APIPort)

	for subdomain, f := range frontends {
		f.RegisterHandlers(subdomain)
		if err := f.Init(serverHost, apiHost); err != nil {
			return nil, err
		}
	}

	// Register public files such as CSS, JS, and Images.
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("templates/public"))))

	// Register our landing page.
	http.HandleFunc("/", landingPage)

	return &Frontend{frontends: frontends, serverHost: serverHost}, nil
}

func (f *Frontend) Port() int {
	return f.serverHost.Port()
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/show", http.StatusSeeOther)
}
