package main

import (
	"tracker/server"
	"tracker/trackable/show"
)

func main() {
	apis := map[string]server.API{
		"api/show": &show.API{},
	}

	backend, err := server.NewBackend(apis)
	if err != nil {
		panic(err)
	}
	if err = server.Serve(backend.Port()); err != nil {
		panic(err)
	}
}
