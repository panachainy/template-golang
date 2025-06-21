package repositories

import (
	"template-golang/database"
	"template-golang/modules/auth/entities"
	"template-golang/modules/auth/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) database.Database {
	// Create SQLite database with in-memory DSN
	sqliteDB, err := database.ProvideSqliteDatabaseWithMigrationPath(":memory:", false, "file:../../../db/migrations")
	assert.NoError(t, err)

	err = sqliteDB.MigrateUp()
	assert.NoError(t, err)

	return sqliteDB
}

func TestProvideAuthRepository(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	assert.NotNil(t, repo)
	assert.Equal(t, testDB, repo.db)
}

func TestAuthRepository_UpsertData_CreateNewAuthRecord(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		ID:       "auth-123",
		Username: "",
		Password: "",
		Email:    "testuser@example.com",
		Role:     models.RoleUser,
		Active:   true,
		AuthMethods: []entities.AuthMethod{
			{
				ID:                "method-123",
				Provider:          entities.ProviderLocal,
				ProviderID:        "user123",
				UserID:            "user123",
				Name:              "Test User",
				FirstName:         "Test",
				LastName:          "User",
				AccessToken:       "token123",
				RefreshToken:      "refresh123",
				IDToken:           "id123",
				AccessTokenSecret: "secret123",
			},
		},
	}

	err := repo.UpsertData(auth)

	assert.NoError(t, err)
	assert.Equal(t, "auth-123", auth.ID)
	assert.Equal(t, "testuser@example.com", auth.Email)
	assert.Equal(t, models.RoleUser, auth.Role)
	assert.True(t, auth.Active)
	assert.Len(t, auth.AuthMethods, 1)
	assert.Equal(t, entities.ProviderLocal, auth.AuthMethods[0].Provider)
	assert.Equal(t, "user123", auth.AuthMethods[0].ProviderID)
	assert.Equal(t, "user123", auth.AuthMethods[0].UserID)
	assert.Equal(t, "Test User", auth.AuthMethods[0].Name)
	assert.Equal(t, "Test", auth.AuthMethods[0].FirstName)
	assert.Equal(t, "User", auth.AuthMethods[0].LastName)
	assert.Equal(t, "token123", auth.AuthMethods[0].AccessToken)
	assert.Equal(t, "refresh123", auth.AuthMethods[0].RefreshToken)
	assert.Equal(t, "id123", auth.AuthMethods[0].IDToken)
	assert.Equal(t, "secret123", auth.AuthMethods[0].AccessTokenSecret)
}

func TestAuthRepository_UpsertData_UpdateExistingAuthRecord(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	// Create existing record first
	existingAuth := &entities.Auth{
		ID:       "auth-456",
		Username: "olduser",
		Password: "oldpassword",
		Email:    "old@example.com",
		Role:     models.RoleUser,
		Active:   true,
	}
	testDB.GetDb().Create(existingAuth)

	newAuthData := &entities.Auth{
		ID:       existingAuth.ID, // Use existing ID to update
		Username: "existinguser",
		Password: "hashedpassword",
		Email:    "existing@example.com",
		Role:     models.RoleAdmin,
		Active:   false,
		AuthMethods: []entities.AuthMethod{
			{
				ID:           "method-456",
				Provider:     entities.ProviderFirebase,
				ProviderID:   "google123",
				AccessToken:  "newtoken",
				RefreshToken: "newrefresh",
			},
		},
	}

	err := repo.UpsertData(newAuthData)
	assert.NoError(t, err)

	var retrievesData []*entities.Auth
	testDB.GetDb().Find(&retrievesData)

	// Verify only one record exists (updated, not created new)
	assert.Len(t, retrievesData, 1)

	// Verify the record was updated with new data
	updatedAuth := retrievesData[0]
	assert.Equal(t, existingAuth.ID, updatedAuth.ID) // Same ID as existing record
	assert.Equal(t, "existinguser", updatedAuth.Username)
	assert.Equal(t, "hashedpassword", updatedAuth.Password)
	assert.Equal(t, "existing@example.com", updatedAuth.Email)
	assert.Equal(t, models.RoleAdmin, updatedAuth.Role)
	assert.False(t, updatedAuth.Active)

	// Verify auth methods were updated
	var authMethods []entities.AuthMethod
	testDB.GetDb().Where("auth_id = ?", updatedAuth.ID).Find(&authMethods)
	assert.Len(t, authMethods, 1)
	assert.Equal(t, entities.ProviderFirebase, authMethods[0].Provider)
	assert.Equal(t, "google123", authMethods[0].ProviderID)
	assert.Equal(t, "newtoken", authMethods[0].AccessToken)
	assert.Equal(t, "newrefresh", authMethods[0].RefreshToken)
}

func TestAuthRepository_UpsertData_CreateAuthWithMultipleAuthMethods(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		ID:       "auth-789",
		Username: "multiuser",
		Email:    "multi@example.com",
		Role:     models.RoleUser,
		Active:   true,
		AuthMethods: []entities.AuthMethod{
			{
				ID:          "method-789-1",
				Provider:    entities.ProviderLocal,
				ProviderID:  "local789",
				UserID:      "user789",
				Name:        "Multi User",
				FirstName:   "Multi",
				LastName:    "User",
				AccessToken: "localtoken",
			},
			{
				ID:           "method-789-2",
				Provider:     entities.ProviderFirebase,
				ProviderID:   "firebase789",
				UserID:       "user789",
				Name:         "Multi User",
				FirstName:    "Multi",
				LastName:     "User",
				AccessToken:  "firebasetoken",
				RefreshToken: "firebaserefresh",
			},
		},
	}

	err := repo.UpsertData(auth)

	assert.NoError(t, err)
	assert.Equal(t, "auth-789", auth.ID)
	assert.Len(t, auth.AuthMethods, 2)

	// Check first auth method
	localMethod := auth.AuthMethods[0]
	assert.Equal(t, entities.ProviderLocal, localMethod.Provider)
	assert.Equal(t, "local789", localMethod.ProviderID)
	assert.Equal(t, "user789", localMethod.UserID)

	// Check second auth method
	firebaseMethod := auth.AuthMethods[1]
	assert.Equal(t, entities.ProviderFirebase, firebaseMethod.Provider)
	assert.Equal(t, "firebase789", firebaseMethod.ProviderID)
	assert.Equal(t, "user789", firebaseMethod.UserID)
}

func TestAuthRepository_UpsertData_CreateAuthWithEmptyRequiredFieldsShouldFail(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		ID:       "auth-invalid",
		Username: "", // Empty required field
		Email:    "invalid@example.com",
		Role:     models.RoleUser,
		Active:   true,
	}

	err := repo.UpsertData(auth)
	assert.Error(t, err)
}
