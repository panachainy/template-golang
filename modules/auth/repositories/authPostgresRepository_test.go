package repositories

import (
	"fmt"
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

func TestInsertData(t *testing.T) {
	tests := []struct {
		name        string
		auth        *entities.Auth
		expectError bool
		expectedMsg string
	}{
		{
			name: "Success - basic auth without methods",
			auth: &entities.Auth{
				ID:       "test-id-1",
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: false,
		},
		{
			name: "Success - auth with single auth method",
			auth: &entities.Auth{
				ID:       "test-id-2",
				Username: "testuser2",
				Email:    "test2@example.com",
				Role:     "user",
				Active:   true,
				AuthMethods: []entities.AuthMethod{
					{
						AuthID:     "test-id-2",
						Provider:   entities.ProviderLocal,
						ProviderID: "local-user-2",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Success - auth with multiple auth methods",
			auth: &entities.Auth{
				ID:       "test-id-3",
				Username: "testuser3",
				Email:    "test3@example.com",
				Role:     "admin",
				Active:   true,
				AuthMethods: []entities.AuthMethod{
					{
						AuthID:       "test-id-3",
						Provider:     entities.ProviderLocal,
						ProviderID:   "local-user-3",
						AccessToken:  "local-token",
						RefreshToken: "local-refresh",
					},
					{
						AuthID:       "test-id-3",
						Provider:     entities.ProviderFirebase,
						ProviderID:   "firebase-user-3",
						AccessToken:  "firebase-token",
						RefreshToken: "firebase-refresh",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Success - minimal auth data",
			auth: &entities.Auth{
				ID:     "test-id-4",
				Email:  "minimal@example.com",
				Role:   "user",
				Active: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := &authPostgresRepository{db: testDB}

			// Test the actual insertion
			err := repo.InsertData(tt.auth)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedMsg)
				}
			} else {
				assert.NoError(t, err)

				// Verify the data was inserted correctly
				var retrievedAuth entities.Auth
				result := testDB.GetDb().Preload("AuthMethods").First(&retrievedAuth, "id = ?", tt.auth.ID)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.auth.ID, retrievedAuth.ID)
				assert.Equal(t, tt.auth.Email, retrievedAuth.Email)
				assert.Equal(t, len(tt.auth.AuthMethods), len(retrievedAuth.AuthMethods))
			}
		})
	}
}

func TestInsertData_DatabaseInterface(t *testing.T) {
	// Test that the repository properly uses the Database interface
	testDB := setupTestDB(t)
	repo := ProvideAuthRepository(testDB)

	auth := &entities.Auth{
		ID:       "interface-test",
		Email:    "interface@example.com",
		Username: "interfaceuser",
		Role:     "user",
		Active:   true,
		AuthMethods: []entities.AuthMethod{
			{
				AuthID:     "interface-test",
				Provider:   entities.ProviderLocal,
				ProviderID: "local-interface-test",
			},
		},
	}

	// This tests that:
	// 1. The repository accepts Database interface
	// 2. It calls GetDb() method
	// 3. Repository structure is correct
	// 4. Real Auth entity with AuthMethods works properly
	err := repo.InsertData(auth)

	// Should succeed with proper GORM configuration
	assert.NoError(t, err)

	// Verify the data was inserted correctly
	var retrievedAuth entities.Auth
	result := testDB.GetDb().Preload("AuthMethods").First(&retrievedAuth, "id = ?", auth.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, auth.ID, retrievedAuth.ID)
	assert.Equal(t, auth.Email, retrievedAuth.Email)
	assert.Len(t, retrievedAuth.AuthMethods, 1)
	assert.Equal(t, entities.ProviderLocal, retrievedAuth.AuthMethods[0].Provider)
}

// TestRepositoryStructure tests the repository pattern implementation
func TestRepositoryStructure(t *testing.T) {
	testDB := setupTestDB(t)

	// Test that ProvideAuthRepository returns correct type
	repo := ProvideAuthRepository(testDB)
	assert.IsType(t, &authPostgresRepository{}, repo)

	// Test that repository has correct database dependency
	assert.Equal(t, testDB, repo.db)
}

// TestAuthMethodRelationship tests that the Auth and AuthMethod relationship works correctly
func TestAuthMethodRelationship(t *testing.T) {
	testDB := setupTestDB(t)
	repo := &authPostgresRepository{db: testDB}

	// Create auth with multiple methods
	auth := &entities.Auth{
		ID:       "relationship-test",
		Username: "relationuser",
		Email:    "relation@example.com",
		Role:     "user",
		Active:   true,
		AuthMethods: []entities.AuthMethod{
			{
				AuthID:       "relationship-test",
				Provider:     entities.ProviderLocal,
				ProviderID:   "local-rel-test",
				AccessToken:  "local-access",
				RefreshToken: "local-refresh",
			},
			{
				AuthID:       "relationship-test",
				Provider:     entities.ProviderFirebase,
				ProviderID:   "firebase-rel-test",
				AccessToken:  "firebase-access",
				RefreshToken: "firebase-refresh",
			},
			{
				AuthID:       "relationship-test",
				Provider:     entities.ProviderLine,
				ProviderID:   "line-rel-test",
				AccessToken:  "line-access",
				RefreshToken: "line-refresh",
			},
		},
	}

	// Insert the auth record
	err := repo.InsertData(auth)
	assert.NoError(t, err)

	// Retrieve and verify the relationship
	var retrievedAuth entities.Auth
	result := testDB.GetDb().Preload("AuthMethods").First(&retrievedAuth, "id = ?", auth.ID)
	assert.NoError(t, result.Error)

	// Verify basic auth data
	assert.Equal(t, auth.ID, retrievedAuth.ID)
	assert.Equal(t, auth.Username, retrievedAuth.Username)
	assert.Equal(t, auth.Email, retrievedAuth.Email)
	assert.Equal(t, auth.Role, retrievedAuth.Role)
	assert.Equal(t, auth.Active, retrievedAuth.Active)

	// Verify auth methods
	assert.Len(t, retrievedAuth.AuthMethods, 3)

	// Create a map for easier verification
	methodsByProvider := make(map[entities.Provider]entities.AuthMethod)
	for _, method := range retrievedAuth.AuthMethods {
		methodsByProvider[method.Provider] = method
	}

	// Verify each auth method
	localMethod := methodsByProvider[entities.ProviderLocal]
	assert.Equal(t, "relationship-test", localMethod.AuthID)
	assert.Equal(t, "local-rel-test", localMethod.ProviderID)
	assert.Equal(t, "local-access", localMethod.AccessToken)
	assert.Equal(t, "local-refresh", localMethod.RefreshToken)

	firebaseMethod := methodsByProvider[entities.ProviderFirebase]
	assert.Equal(t, "relationship-test", firebaseMethod.AuthID)
	assert.Equal(t, "firebase-rel-test", firebaseMethod.ProviderID)
	assert.Equal(t, "firebase-access", firebaseMethod.AccessToken)
	assert.Equal(t, "firebase-refresh", firebaseMethod.RefreshToken)

	lineMethod := methodsByProvider[entities.ProviderLine]
	assert.Equal(t, "relationship-test", lineMethod.AuthID)
	assert.Equal(t, "line-rel-test", lineMethod.ProviderID)
	assert.Equal(t, "line-access", lineMethod.AccessToken)
	assert.Equal(t, "line-refresh", lineMethod.RefreshToken)
}

func TestInsertData_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		auth        *entities.Auth
		expectError bool
		expectedMsg string
	}{
		{
			name: "Error - duplicate ID",
			auth: &entities.Auth{
				ID:       "duplicate-id",
				Username: "user1",
				Email:    "user1@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: true,
			expectedMsg: "UNIQUE constraint failed",
		},
		{
			name: "Error - duplicate email",
			auth: &entities.Auth{
				ID:       "test-id-unique-1",
				Username: "user2",
				Email:    "duplicate@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: true,
			expectedMsg: "UNIQUE constraint failed",
		},
		{
			name: "Error - duplicate username",
			auth: &entities.Auth{
				ID:       "test-id-unique-2",
				Username: "duplicateuser",
				Email:    "user3@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: true,
			expectedMsg: "UNIQUE constraint failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := &authPostgresRepository{db: testDB}

			// First, insert the initial record to create duplicates
			if tt.name == "Error - duplicate ID" {
				initialAuth := &entities.Auth{
					ID:       "duplicate-id",
					Username: "initialuser",
					Email:    "initial@example.com",
					Role:     "user",
					Active:   true,
				}
				err := repo.InsertData(initialAuth)
				assert.NoError(t, err)
			} else if tt.name == "Error - duplicate email" {
				initialAuth := &entities.Auth{
					ID:       "initial-id",
					Username: "initialuser",
					Email:    "duplicate@example.com",
					Role:     "user",
					Active:   true,
				}
				err := repo.InsertData(initialAuth)
				assert.NoError(t, err)
			} else if tt.name == "Error - duplicate username" {
				initialAuth := &entities.Auth{
					ID:       "initial-id-2",
					Username: "duplicateuser",
					Email:    "initial2@example.com",
					Role:     "user",
					Active:   true,
				}
				err := repo.InsertData(initialAuth)
				assert.NoError(t, err)
			}

			// Now test the duplicate insertion
			err := repo.InsertData(tt.auth)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInsertData_NilInput(t *testing.T) {
	testDB := setupTestDB(t)
	repo := &authPostgresRepository{db: testDB}

	err := repo.InsertData(nil)
	assert.Error(t, err)
}

func TestInsertData_EmptyFields(t *testing.T) {
	tests := []struct {
		name        string
		auth        *entities.Auth
		expectError bool
	}{
		{
			name: "Empty ID",
			auth: &entities.Auth{
				ID:       "",
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: false, // GORM will auto-generate ID if empty
		},
		{
			name: "Empty email",
			auth: &entities.Auth{
				ID:       "test-empty-email",
				Username: "testuser",
				Email:    "",
				Role:     "user",
				Active:   true,
			},
			expectError: false, // Depends on entity validation
		},
		{
			name: "Empty username",
			auth: &entities.Auth{
				ID:       "test-empty-username",
				Username: "",
				Email:    "empty-username@example.com",
				Role:     "user",
				Active:   true,
			},
			expectError: false, // Depends on entity validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := &authPostgresRepository{db: testDB}

			err := repo.InsertData(tt.auth)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInsertData_AuthMethodValidation(t *testing.T) {
	tests := []struct {
		name        string
		auth        *entities.Auth
		expectError bool
		expectedMsg string
	}{
		{
			name: "Success - auth method with minimal data",
			auth: &entities.Auth{
				ID:     "minimal-method-test",
				Email:  "minimal-method@example.com",
				Role:   "user",
				Active: true,
				AuthMethods: []entities.AuthMethod{
					{
						AuthID:     "minimal-method-test",
						Provider:   entities.ProviderLocal,
						ProviderID: "minimal-provider-id",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Success - auth method with all fields",
			auth: &entities.Auth{
				ID:     "full-method-test",
				Email:  "full-method@example.com",
				Role:   "user",
				Active: true,
				AuthMethods: []entities.AuthMethod{
					{
						AuthID:       "full-method-test",
						Provider:     entities.ProviderFirebase,
						ProviderID:   "firebase-full-id",
						AccessToken:  "full-access-token",
						RefreshToken: "full-refresh-token",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Auth method with empty provider ID",
			auth: &entities.Auth{
				ID:     "empty-provider-test",
				Email:  "empty-provider@example.com",
				Role:   "user",
				Active: true,
				AuthMethods: []entities.AuthMethod{
					{
						AuthID:     "empty-provider-test",
						Provider:   entities.ProviderLocal,
						ProviderID: "",
					},
				},
			},
			expectError: false, // Depends on entity validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := setupTestDB(t)
			repo := &authPostgresRepository{db: testDB}

			err := repo.InsertData(tt.auth)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedMsg)
				}
			} else {
				assert.NoError(t, err)

				// Verify the insertion
				var retrievedAuth entities.Auth
				result := testDB.GetDb().Preload("AuthMethods").First(&retrievedAuth, "id = ?", tt.auth.ID)
				assert.NoError(t, result.Error)
				assert.Equal(t, len(tt.auth.AuthMethods), len(retrievedAuth.AuthMethods))
			}
		})
	}
}
