package main

import (
	"github.com/nixmaldonado/blazeMailer/config"
	"github.com/nixmaldonado/blazeMailer/server"
)

func main() {
	config.InitConfig()
	server.HandleRequests()
}
