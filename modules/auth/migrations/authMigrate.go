package main

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/auth/entities"
)

// TODO: migrate to API
func main() {
	conf := config.Provide()
	db := database.Provide(conf)
	authMigrate(db)
}

func authMigrate(db database.Database) {
	err := db.GetDb().AutoMigrate(&entities.Auth{}, &entities.AuthMethod{})
	if err != nil {
		panic(err)
	}

	// Add any initial auth data if needed
	// Example: Create default admin user
	// db.GetDb().CreateInBatches([]entities.User{
	//     {Email: "admin@example.com", Role: "admin"},
	// }, 10)
}
