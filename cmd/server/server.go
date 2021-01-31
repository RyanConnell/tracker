package main

import (
	"tracker/server"
)

func main() {
	if err := server.Launch(); err != nil {
		panic(err)
	}
}
