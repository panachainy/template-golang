package main

import (
	"template-golang/config"
)

func main() {
	conf := config.GetConfig()
	print(conf)
	// db := database.NewPostgresDatabase(conf)
	// server.NewEchoServer(conf, db).Start()
}
