package main

import (
	"tracker/server"
	"tracker/trackable/show"
)

func main() {
	apis := map[string]server.WebFrontend{
		"show": &show.Frontend{},
	}

	frontend, err := server.NewFrontend(apis)
	if err != nil {
		panic(err)
	}
	if err = server.Serve(frontend.Port()); err != nil {
		panic(err)
	}
}
