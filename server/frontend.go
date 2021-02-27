package server

import (
	"net/http"

	"tracker/server/host"
	"tracker/web"
)

type WebFrontend interface {
	Init(serverHost, apiHost *host.Host) error
	RegisterHandlers(subdomain string)
}

type Frontend struct {
	frontends  map[string]WebFrontend
	serverHost *host.Host
}

func NewFrontend(frontends map[string]WebFrontend) (*Frontend, error) {
	settings, err := NewSettings()
	if err != nil {
		return nil, err
	}

	serverHost := host.NewHost(settings.Hostname, settings.Port)
	apiHost := host.NewHost(settings.APIHostname, settings.APIPort)

	for subdomain, f := range frontends {
		f.RegisterHandlers(subdomain)
		if err = f.Init(serverHost, apiHost); err != nil {
			return nil, err
		}
	}

	http.Handle("/public/", http.FileServer(http.FS(web.Static)))

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
