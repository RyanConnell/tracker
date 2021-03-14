package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tracker/internal/frontend"
	"tracker/internal/httpserver"
	"tracker/web"

	oldserver "tracker/server"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

// Config used for the frontend.
// TODO: Add default values once we deprecate file settings.
type Config struct {
	Debug       bool   `envconfig:"DEBUG"`
	Port        int    `envconfig:"PORT"`
	BackendAddr string `split_words:"true"`
}

func main() {
	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg Config
	if err := envconfig.Process("frontend", &cfg); err != nil {
		return fmt.Errorf("unable to process config: %w", err)
	}

	var log *zap.Logger
	if cfg.Debug {
		log, _ = zap.NewDevelopment()
	} else {
		log, _ = zap.NewProduction()
	}

	// For compatibility reasons, fill in options from the settings
	settings, err := oldserver.NewSettings()
	if err != nil {
		return fmt.Errorf("unable to parse settings: %w", err)
	}
	if cfg.Port == 0 {
		cfg.Port = settings.Port
	}
	if cfg.BackendAddr == "" {
		cfg.BackendAddr = fmt.Sprintf("%s:%d", settings.APIHostname, settings.APIPort)
	}

	// Initialize the show frontend
	showFrontend, err := frontend.NewShow(cfg.BackendAddr)
	if err != nil {
		return fmt.Errorf("unable to init show frontend: %w", err)
	}

	s := httpserver.NewServer(map[string]httpserver.Component{
		"/show":   showFrontend,
		"/public": frontend.NewStatic(web.Static),
		"/":       frontend.NewRedirect(http.StatusTemporaryRedirect, "/show/"),
	},
		httpserver.Logger(log),
	)

	return s.Run(cfg.Port)
}
