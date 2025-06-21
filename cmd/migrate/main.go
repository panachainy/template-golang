package main

import (
	"flag"
	"os"
	"strconv"
	"template-golang/config"
	"template-golang/database"

	"github.com/labstack/gommon/log"
)

func main() {
	var action = flag.String("action", "up", "Migration action: up, down, version")
	var steps = flag.String("steps", "1", "Number of steps for rollback (only used with down action)")
	flag.Parse()

	conf := config.Provide()
	db := database.Provide(conf)
	migrationManager := database.ProvideMigrationManager(db, conf)

	switch *action {
	case "up":
		if err := migrationManager.RunMigrations(); err != nil {
			log.Errorf("Failed to run migrations: %v", err)
			os.Exit(1)
		}
		log.Info("Migrations completed successfully")

	case "down":
		stepCount, err := strconv.Atoi(*steps)
		if err != nil {
			log.Errorf("Invalid steps steps: %v err: %v", *steps, err)
			os.Exit(1)
		}
		if err := migrationManager.RollbackMigrations(stepCount); err != nil {
			log.Errorf("Failed to rollback migrations: %v", err)
			os.Exit(1)
		}
		log.Infof("Rollback of %d steps completed successfully", stepCount)

	case "version":
		version, dirty, err := migrationManager.GetVersion()
		if err != nil {
			log.Errorf("Failed to get migration version: %v", err)
			os.Exit(1)
		}
		log.Infof("Current migration version: %d (dirty: %t)", version, dirty)

	default:
		log.Errorf("Unknown action: %s. Available actions: up, down, version", *action)
		os.Exit(1)
	}
}
