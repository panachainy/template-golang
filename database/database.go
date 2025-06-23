//go:generate mockgen -source=database.go -destination=../mock/mock_database.go -package=mock

package database

import "gorm.io/gorm"

type Database interface {
	GetDb() *gorm.DB
	MigrateUp() error
	MigrateDown(steps int) error
	GetVersion() (uint, bool, error)
	Close() error
}
