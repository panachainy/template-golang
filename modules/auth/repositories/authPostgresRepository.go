package repositories

import (
	"template-golang/database"
	"template-golang/modules/auth/entities"

	"github.com/labstack/gommon/log"
)

type authPostgresRepository struct {
	db database.Database
}

func ProvideAuthRepository(db database.Database) *authPostgresRepository {
	return &authPostgresRepository{db: db}
}

func (r *authPostgresRepository) InsertData(in *entities.Auth) error {
	result := r.db.GetDb().Create(in)

	if result.Error != nil {
		log.Errorf("InsertAuth: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertAuth: %v", result.RowsAffected)
	return nil
}
