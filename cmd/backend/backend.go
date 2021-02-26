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
	apis := map[string]server.API{
		"api/show": &show.API{},
	}

	backend, err := server.NewBackend(apis)
	if err != nil {
		return fmt.Errorf("unable to initialize backend server: %w", err)
	}

	if err := server.Serve(backend.Port()); err != nil {
		return fmt.Errorf("unable to serve on port %d: %w", backend.Port(), err)
	}

	return nil
}
