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

func (r *authPostgresRepository) UpsertData(in *entities.Auth) error {
	result := r.db.GetDb().Save(in)

	if result.Error != nil {
		log.Errorf("UpsertAuth: %v", result.Error)
		return result.Error
	}

	log.Debugf("UpsertAuth: %v", result.RowsAffected)
	return nil
}

func (r *authPostgresRepository) Gets(limit int) ([]*entities.Auth, error) {
	var auths []*entities.Auth
	// result := r.db.GetDb().Model(&entities.Auth{}).Limit(limit).Find(&auths)
	result := r.db.GetDb().Limit(limit).Find(&auths)

	if result.Error != nil {
		log.Errorf("Gets: %v", result.Error)
		return nil, result.Error
	}

	log.Debugf("Gets: %v rows retrieved", result.RowsAffected)
	return auths, nil
}
