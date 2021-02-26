package main

import (
	"fmt"
	"log"
	"os"

	"tracker/server"
	"tracker/trackable/show"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	apis := map[string]server.WebFrontend{
		"show": &show.Frontend{},
	}

	frontend, err := server.NewFrontend(apis)
	if err != nil {
		return fmt.Errorf("unable to create a new frontend: %w", err)
	}
	if err = server.Serve(frontend.Port()); err != nil {
		return fmt.Errorf("unable to serve on port %d: %w", frontend.Port(), err)
	}

	return nil
}
