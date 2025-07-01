package main

import (
	"flag"
	"os"
	"strconv"
	"template-golang/config"
	"template-golang/database"
	"template-golang/pkg/logger"
)

func main() {
	var action = flag.String("action", "up", "Migration action: up, down, version")
	var steps = flag.String("steps", "1", "Number of steps for rollback (only used with down action)")
	flag.Parse()

	conf := config.Provide(&config.ConfigOption{
		ConfigPath: ".",
	})
	db := database.NewPostgres(conf)

	switch *action {
	case "up":
		if err := db.MigrateUp(); err != nil {
			logger.Errorf("Failed to run migrations: %v", err)
			os.Exit(1)
		}
		logger.Info("Migrations completed successfully")

	case "down":
		stepCount, err := strconv.Atoi(*steps)
		if err != nil {
			logger.Errorf("Invalid steps steps: %v err: %v", *steps, err)
			os.Exit(1)
		}
		if err := db.MigrateDown(stepCount); err != nil {
			logger.Errorf("Failed to rollback migrations: %v", err)
			os.Exit(1)
		}
		logger.Infof("Rollback of %d steps completed successfully", stepCount)

	case "version":
		version, dirty, err := db.GetVersion()
		if err != nil {
			logger.Errorf("Failed to get migration version: %v", err)
			os.Exit(1)
		}
		logger.Infof("Current migration version: %d (dirty: %t)", version, dirty)

	default:
		logger.Errorf("Unknown action: %s. Available actions: up, down, version", *action)
		os.Exit(1)
	}
}
