package main

import (
	"fmt"
	"os"

	"tracker/server"
	"tracker/trackable/common"
)

var settingsFilePath string = "server.conf"

var defaultSettings map[string]string = map[string]string{
	"ip":                   "localhost",
	"port":                 "8080",
	"google_client_id":     "GOOGLE_CLIENT_ID",
	"google_client_secret": "GOOGLE_CLIENT_SECRET",
}

func main() {
	settings, err := loadSettings(settingsFilePath)
	if err != nil {
		panic(err)
	}

	host := &common.Host{}
	err = host.Init(settings)
	if err != nil {
		panic(err)
	}

	server.Launch(host, settings)
}

func loadSettings(path string) (map[string]string, error) {
	if common.FileExists(path) {
		return common.LoadSettings(path), nil
	}

	// Create our file since it doesn't exist.
	file, err := os.Create(path)
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
