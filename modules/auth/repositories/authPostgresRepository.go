package repositories

import (
	"template-golang/database"
	"template-golang/modules/auth/entities"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type authPostgresRepository struct {
	db        database.Database
	validator *validator.Validate
}

func ProvideAuthRepository(db database.Database) *authPostgresRepository {
	return &authPostgresRepository{db: db,
		validator: validator.New(),
	}
}

func (r *authPostgresRepository) UpsertData(in *entities.Auth) error {
	if err := r.validator.Struct(in); err != nil {
		log.Errorf("UpsertData validation failed: %v", err)
		return err
	}

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
