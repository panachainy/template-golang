package main

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/cockroach/entities"
)

// TODO: migrate to API
func main() {
	conf := config.Provide()
	db := database.Provide(conf)
	cockroachMigrate(db)
}

func cockroachMigrate(db database.Database) {
	err := db.GetDb().AutoMigrate(&entities.Cockroach{})
	if err != nil {
		panic(err)
	}

	db.GetDb().CreateInBatches([]entities.Cockroach{
		{Amount: 1},
		{Amount: 2},
		{Amount: 2},
		{Amount: 5},
		{Amount: 3},
	}, 10)
}
