package server

import (
	"fmt"

	"tracker/trackable/common"
)

type API interface {
	Init(*common.Host) error
	RegisterHandlers(subdomain string)
}

type Backend struct {
	apis map[string]API
	host *common.Host
}

func NewBackend(apis map[string]API) (*Backend, error) {
	settings, err := loadOrCreateSettings()
	if err != nil {
		return nil, err
	}

	host := &common.Host{}
	if err = host.Init(settings); err != nil {
		return nil, err
	}

	for subdomain, api := range apis {
		fmt.Printf("Registering api handler for %q\n", subdomain)
		api.RegisterHandlers(subdomain)
		if err = api.Init(host); err != nil {
			return nil, err
		}
	}

	return &Backend{apis: apis, host: host}, nil
}

func (b *Backend) Port() int {
	return b.host.Port()
}
