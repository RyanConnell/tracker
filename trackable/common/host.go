package common

import (
	"fmt"
	"strconv"
)

type Host struct {
	ip       string
	port     int
	settings map[string]string
}

func (h *Host) Init(settings map[string]string) error {
	var ok bool
	if h.ip, ok = settings["ip"]; !ok {
		fmt.Println("Host IP not set - Defaulting to 'localhost'")
		h.ip = "localhost"
		settings["ip"] = "localhost"
	}

	if portStr, ok := settings["port"]; !ok {
		fmt.Println("Host Port not set - Defaulting to 8080")
		h.port = 8080
		settings["port"] = "8080"
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
