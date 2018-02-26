package main

import (
	"tracker/server"
	"tracker/trackable/common"
)

var hostFile string = "host.conf"

func main() {
	host := &common.Host{}
	err := host.Init(hostFile)
	if err != nil {
		panic(err)
	}

	server.Launch(host)
}
