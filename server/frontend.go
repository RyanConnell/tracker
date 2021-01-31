package server

import (
	"net/http"

	"tracker/server/host"
)

type WebFrontend interface {
	Init(*host.Host) error
	RegisterHandlers(subdomain string)
}

type Frontend struct {
	frontends map[string]WebFrontend
	host      *host.Host
}

func NewFrontend(frontends map[string]WebFrontend) (*Frontend, error) {
	settings, err := loadOrCreateSettings()
	if err != nil {
		return nil, err
	}

	host := &host.Host{}
	if err = host.Init(settings, "ip", "port"); err != nil {
		return nil, err
	}

	for subdomain, f := range frontends {
		f.RegisterHandlers(subdomain)
		if err = f.Init(host); err != nil {
			return nil, err
		}
	}

	// Register public files such as CSS, JS, and Images.
	http.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir("templates/public"))))

	// Register our landing page.
	http.HandleFunc("/", landingPage)

	return &Frontend{frontends: frontends, host: host}, nil
}

func (f *Frontend) Port() int {
	return f.host.Port()
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/show", http.StatusSeeOther)
}
