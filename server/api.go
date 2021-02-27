package server

import (
	"fmt"

	"tracker/server/host"
)

type API interface {
	Init(*host.Host) error
	RegisterHandlers(subdomain string)
}

type Backend struct {
	apis map[string]API
	host *host.Host
}

func NewBackend(settings *Settings, apis map[string]API) (*Backend, error) {
	host := host.NewHost(settings.Hostname, settings.APIPort)
	for subdomain, api := range apis {
		fmt.Printf("Registering api handler for %q\n", subdomain)
		api.RegisterHandlers(subdomain)
		if err := api.Init(host); err != nil {
			return nil, err
		}
	}

	return &Backend{apis: apis, host: host}, nil
}

func (b *Backend) Port() int {
	return b.host.Port()
}
