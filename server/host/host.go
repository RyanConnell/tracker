package host

import (
	"fmt"
)

type Host struct {
	hostname string
	port     int
}

func NewHost(hostname string, port int) *Host {
	return &Host{hostname: hostname, port: port}
}

func (h *Host) Hostname() string {
	return h.hostname
}

func (h *Host) Port() int {
	return h.port
}

func (h *Host) Address() string {
	return fmt.Sprintf("http://%s:%d", h.hostname, h.port)
}
