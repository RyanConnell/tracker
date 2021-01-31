package host

import (
	"fmt"
	"strconv"
)

type Host struct {
	ip       string
	port     int
	settings map[string]string
}

func (h *Host) Init(settings map[string]string, ipKey, portKey string) error {
	var ok bool
	if h.ip, ok = settings[ipKey]; !ok {
		fmt.Printf("%q not set - Defaulting to 'localhost'\n", ipKey)
		h.ip = "localhost"
		settings[ipKey] = "localhost"
	}

	if portStr, ok := settings[portKey]; !ok {
		fmt.Printf("%q not set - Defaulting to 8080\n", portKey)
		h.port = 8080
		settings[portKey] = "8080"
	} else {
		var err error
		if h.port, err = strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("Unable to parse port: %v", portStr)
		}
	}

	h.settings = settings
	return nil
}

func (h *Host) IP() string {
	return h.ip
}

func (h *Host) Port() int {
	return h.port
}

func (h *Host) Address() string {
	return fmt.Sprintf("http://%s:%d", h.ip, h.port)
}
