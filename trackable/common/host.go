package common

import (
	"fmt"
	"os"
)

type Host struct {
	ip       string
	port     int
	settings map[string]string
}

func (h *Host) Init(path string) error {
	if FileExists(path) {
		return h.init(LoadSettings(path))
	}
	err := h.init(map[string]string{})
	if err != nil {
		return err
	}
	return h.writeToFile(path)
}

func (h *Host) init(settings map[string]string) error {
	var ok bool
	if h.ip, ok = settings["ip"]; !ok {
		fmt.Println("Host IP not set - Defaulting to 'localhost'")
		h.ip = "localhost"
		settings["ip"] = "localhost"
	}

	if portStr, ok := settings["port"]; !ok {
		fmt.Println("Host Port not set - Defaulting to 80")
		h.port = 80
		settings["port"] = "80"
	} else {
		var err error
		if h.port, err = StringToInt(portStr); err != nil {
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

func (h *Host) writeToFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range h.settings {
		_, err := file.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		if err != nil {
			return err
		}
	}
	return nil
}
