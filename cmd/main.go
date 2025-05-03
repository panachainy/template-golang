package main

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	server.NewGinServer(conf, db).Start()
}
