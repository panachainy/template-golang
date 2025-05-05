package main

import (
	"template-golang/config"
	"template-golang/server"
)

func main() {
	conf := config.GetConfig()
	server.NewGinServer(conf).Start()
}
