package repositories

import (
	"template-golang/modules/auth/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// testDatabase implements the Database interface for testing
type testDatabase struct {
	db *gorm.DB
}

func (t *testDatabase) GetDb() *gorm.DB {
	return t.db
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *testDatabase {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress logs during tests
	})
	assert.NoError(t, err)

	// Auto-migrate the real auth entities for testing
	err = db.AutoMigrate(&entities.Auth{}, &entities.AuthMethod{})
	assert.NoError(t, err)

	return &testDatabase{db: db}
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
		UserID:    "user123",
		Username:  "",
		Password:  "",
		Email:     "testuser@example.com",
		Role:      "user",
		Active:    true,
		Name:      "Test User",
		FirstName: "Test",
		LastName:  "User",
		AuthMethods: []entities.AuthMethod{
			{
				Provider:          entities.ProviderLocal,
				ProviderID:        "user123",
				AccessToken:       "token123",
				RefreshToken:      "refresh123",
				IDToken:           "id123",
				AccessTokenSecret: "secret123",
			},
		},
	}

	err := repo.UpsertData(auth)

	assert.NoError(t, err)
	assert.NotZero(t, auth.ID)
	assert.Equal(t, "user123", auth.UserID)
	assert.Equal(t, "testuser@example.com", auth.Email)
	assert.Equal(t, "user", auth.Role)
	assert.True(t, auth.Active)
	assert.Equal(t, "Test User", auth.Name)
	assert.Equal(t, "Test", auth.FirstName)
	assert.Equal(t, "User", auth.LastName)
	assert.Len(t, auth.AuthMethods, 1)
	assert.Equal(t, entities.ProviderLocal, auth.AuthMethods[0].Provider)
	assert.Equal(t, "user123", auth.AuthMethods[0].ProviderID)
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
		UserID:    "user456",
		Username:  "olduser",
		Password:  "oldpassword",
		Email:     "old@example.com",
		Role:      "user",
		Active:    true,
		Name:      "Old User",
		FirstName: "Old",
		LastName:  "User",
	}
	testDB.GetDb().Create(existingAuth)

	auth := &entities.Auth{
		UserID:    "user456",
		Username:  "existinguser",
		Password:  "hashedpassword",
		Email:     "existing@example.com",
		Role:      "admin",
		Active:    false,
		Name:      "Updated User",
		FirstName: "Updated",
		LastName:  "User",
		AuthMethods: []entities.AuthMethod{
			{
				Provider:     entities.ProviderFirebase,
				ProviderID:   "google123",
				AccessToken:  "newtoken",
				RefreshToken: "newrefresh",
			},
		},
	}

	err := repo.UpsertData(auth)

	assert.NoError(t, err)
	assert.NotZero(t, auth.ID)
	assert.Equal(t, "user456", auth.UserID)
	assert.Equal(t, "existinguser", auth.Username)
	assert.Equal(t, "existing@example.com", auth.Email)
	assert.Equal(t, "admin", auth.Role)
	assert.False(t, auth.Active)
	assert.Equal(t, "Updated User", auth.Name)
	assert.Len(t, auth.AuthMethods, 1)
	assert.Equal(t, entities.ProviderFirebase, auth.AuthMethods[0].Provider)
}

func TestAuthRepository_UpsertData_CreateAuthWithMultipleAuthMethods(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		UserID:    "user789",
		Username:  "multiuser",
		Email:     "multi@example.com",
		Role:      "user",
		Active:    true,
		Name:      "Multi User",
		FirstName: "Multi",
		LastName:  "User",
		AuthMethods: []entities.AuthMethod{
			{
				Provider:    entities.ProviderLocal,
				ProviderID:  "local789",
				AccessToken: "localtoken",
			},
			{
				Provider:     entities.ProviderLocal,
				ProviderID:   "github789",
				AccessToken:  "githubtoken",
				RefreshToken: "githubrefresh",
			},
		},
	}

	err := repo.UpsertData(auth)

	assert.NoError(t, err)
	assert.NotZero(t, auth.ID)
	assert.Equal(t, "user789", auth.UserID)
	assert.Len(t, auth.AuthMethods, 2)

	// Check first auth method
	localMethod := auth.AuthMethods[0]
	assert.Equal(t, entities.ProviderLocal, localMethod.Provider)
	assert.Equal(t, "local789", localMethod.ProviderID)

	// Check second auth method
	githubMethod := auth.AuthMethods[1]
	assert.Equal(t, entities.ProviderLocal, githubMethod.Provider)
	assert.Equal(t, "github789", githubMethod.ProviderID)
}

func TestAuthRepository_UpsertData_CreateAuthWithEmptyRequiredFieldsShouldFail(t *testing.T) {
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		UserID: "", // Empty required field
		Email:  "invalid@example.com",
		Role:   "user",
		Active: true,
		Name:   "Invalid User",
	}

	err := repo.UpsertData(auth)
	assert.Error(t, err)
}
