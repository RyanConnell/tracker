package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var settingsFilePath string = "server.conf"

var defaultSettings map[string]string = map[string]string{
	"ip":                   "localhost",
	"port":                 "8080",
	"google_client_id":     "GOOGLE_CLIENT_ID",
	"google_client_secret": "GOOGLE_CLIENT_SECRET",
}

// Loads the settings if they exist, or creates the defaults if not.
func loadOrCreateSettings() (map[string]string, error) {
	if _, err := os.Stat(settingsFilePath); err == nil {
		return readSettingsFromFile(settingsFilePath)
	}

	// Create our file since it doesn't exist.
	file, err := os.Create(settingsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for key, value := range defaultSettings {
		_, err := file.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Settings file created: %s\n", settingsFilePath)
	return defaultSettings, nil
}

// Reads the settings from a file on disk.
func readSettingsFromFile(filename string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	settings := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		s := strings.SplitN(line, ":", 2)
		if len(s) == 2 {
			settings[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
		}
	}
	return settings, nil
}
