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

func (r *authPostgresRepository) GetUserByUserID(userID string) (*entities.Auth, error) {
	var authM entities.AuthMethod
	result := r.db.GetDb().Where("user_id = ?", userID).First(&authM)

	if result.Error != nil {
		log.Errorf("GetUserByUserID: %v", result.Error)
		return nil, result.Error
	}

	// query auth by authM
	var auth entities.Auth
	result = r.db.GetDb().Where("id = ?", authM.AuthID).First(&auth)
	if result.Error != nil {
		log.Errorf("GetUserByUserID: %v", result.Error)
		return nil, result.Error
	}

	if auth.ID == "" {
		log.Errorf("GetUserByUserID: auth not found for user_id %s", userID)
		return nil, nil // or return an error if preferred
	}

	log.Debugf("GetUserByUserID: found auth for user_id %s", userID)
	return &auth, nil
}

func (r *authPostgresRepository) GetAuthIDMethodIDByUserID(userID string) (*GetAuthIdMethodIdResponse, error) {
	var authM entities.AuthMethod
	result := r.db.GetDb().Where("user_id = ?", userID).First(&authM)

	if result.Error != nil {
		// can be in case first login
		log.Errorf("GetAuthIDMethodIDByUserID: %v", result.Error)
		return nil, result.Error
	}

	response := &GetAuthIdMethodIdResponse{
		AuthID:   authM.AuthID,
		MethodID: authM.ID, // assuming AuthMethod has an ID field
	}

	log.Debugf("GetAuthIDMethodIDByUserID: found auth_id %s and method_id %s for user_id %s",
		response.AuthID, response.MethodID, userID)
	return response, nil
}
