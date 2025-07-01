package repositories

import (
	"template-golang/database"
	"template-golang/modules/auth/entities"
	"template-golang/pkg/logger"

	"github.com/go-playground/validator/v10"
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
		logger.Errorf("UpsertData validation failed: %v", err)
		return err
	}

	// Start a transaction to ensure data consistency
	tx := r.db.GetDb().Begin()
	if tx.Error != nil {
		logger.Errorf("UpsertData: failed to begin transaction: %v", tx.Error)
		return tx.Error
	}

	// Store AuthMethods temporarily and clear them from the Auth struct
	authMethods := in.AuthMethods
	in.AuthMethods = nil

	// First, save the Auth record without AuthMethods
	result := tx.Save(in)
	if result.Error != nil {
		tx.Rollback()
		logger.Errorf("UpsertAuth: %v", result.Error)
		return result.Error
	}

	// Now save the AuthMethods with the correct AuthID
	for i := range authMethods {
		authMethods[i].AuthID = in.ID
		result = tx.Save(&authMethods[i])
		if result.Error != nil {
			tx.Rollback()
			logger.Errorf("UpsertAuthMethods: %v", result.Error)
			return result.Error
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.Errorf("UpsertData: failed to commit transaction: %v", err)
		return err
	}

	// Restore AuthMethods to the original struct
	in.AuthMethods = authMethods

	logger.Debugf("UpsertAuth: successfully saved auth and %d auth methods", len(authMethods))
	return nil
}

func (r *authPostgresRepository) Gets(limit int) ([]*entities.Auth, error) {
	var auths []*entities.Auth
	// result := r.db.GetDb().Model(&entities.Auth{}).Limit(limit).Find(&auths)
	result := r.db.GetDb().Limit(limit).Find(&auths)

	if result.Error != nil {
		logger.Errorf("Gets: %v", result.Error)
		return nil, result.Error
	}

	logger.Debugf("Gets: %v rows retrieved", result.RowsAffected)
	return auths, nil
}

func (r *authPostgresRepository) GetUserByUserID(userID string) (*entities.Auth, error) {
	var authM entities.AuthMethod
	result := r.db.GetDb().Where("user_id = ?", userID).First(&authM)

	if result.Error != nil {
		logger.Errorf("GetUserByUserID: %v", result.Error)
		return nil, result.Error
	}

	// query auth by authM
	var auth entities.Auth
	result = r.db.GetDb().Where("id = ?", authM.AuthID).First(&auth)
	if result.Error != nil {
		logger.Errorf("GetUserByUserID: %v", result.Error)
		return nil, result.Error
	}

	if auth.ID == "" {
		logger.Errorf("GetUserByUserID: auth not found for user_id %s", userID)
		return nil, nil // or return an error if preferred
	}

	logger.Debugf("GetUserByUserID: found auth for user_id %s", userID)
	return &auth, nil
}

func (r *authPostgresRepository) GetAuthIDMethodIDByUserID(userID string) (*GetAuthIdMethodIdResponse, error) {
	var authM entities.AuthMethod
	result := r.db.GetDb().Where("user_id = ?", userID).First(&authM)

	if result.Error != nil {
		// can be in case first login
		logger.Errorf("GetAuthIDMethodIDByUserID: %v", result.Error)
		return nil, result.Error
	}

	response := &GetAuthIdMethodIdResponse{
		AuthID:   authM.AuthID,
		MethodID: authM.ID, // assuming AuthMethod has an ID field
	}

	logger.Debugf("GetAuthIDMethodIDByUserID: found auth_id %s and method_id %s for user_id %s",
		response.AuthID, response.MethodID, userID)
	return response, nil
}
