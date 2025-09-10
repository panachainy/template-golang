package usecases

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"template-golang/config"
	db "template-golang/db/sqlc"
	"template-golang/modules/auth/models"
	"template-golang/modules/auth/repositories"
	"template-golang/modules/auth/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth"
)

type jwtUsecaseImpl struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	authRepo   repositories.AuthRepository
}

func NewJWTUsecase(conf *config.Config, authRepo repositories.AuthRepository) JWTUsecase {
	privateKey := loadPrivateKey(conf.Auth.PrivateKeyPath)
	publicKey := &privateKey.PublicKey

	return &jwtUsecaseImpl{
		privateKey: privateKey,
		publicKey:  publicKey,
		authRepo:   authRepo,
	}
}

// we use panic because if not have private key, we cannot run the server
func loadPrivateKey(path string) *ecdsa.PrivateKey {
	// init key
	var keyByteArray []byte
	var key *ecdsa.PrivateKey

	// Validate and clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		panic(fmt.Errorf("failed to resolve absolute path: %w", err))
	}

	// Basic security check - prevent access to system directories
	if strings.Contains(absPath, "/etc/") || strings.Contains(absPath, "/usr/") ||
		strings.Contains(absPath, "/var/") || strings.Contains(absPath, "/root/") ||
		strings.Contains(absPath, "/home/") && !strings.Contains(absPath, "/home/"+os.Getenv("USER")) {
		panic(fmt.Errorf("invalid path: access to system directories not allowed"))
	}

	// Load the private key from a file
	keyByteArray, err = os.ReadFile(absPath) // #nosec G304 -- Path is validated to prevent directory traversal and system file access
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	// Parse the private key
	key, err = jwt.ParseECPrivateKeyFromPEM(keyByteArray)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	return key
}

func (a *jwtUsecaseImpl) GenerateJWT(userID string) (string, error) {
	// TODO: implement MapClaims
	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":       userID,
		"name":      "John",
		"last_name": "Doe",
		"iss":       "my-auth-server-issuer",
		"foo":       2,
	})

	// Set expiration time (e.g., 24 hours from now)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))

	// Sign the token with the private key
	signedString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedString, nil
}

func (a *jwtUsecaseImpl) ValidateJWT(tokenString string) (*models.TokenValidationResult, error) {
	result := &models.TokenValidationResult{
		Valid:    false,
		Expired:  false,
		NotExist: false,
		Claims:   nil,
		UserID:   "",
	}

	// Check if token string is empty
	if tokenString == "" {
		result.NotExist = true
		return result, nil
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// FIXME: check this
		return a.publicKey, nil
	})

	if err != nil {
		// Check if error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			result.Expired = true
			return result, nil
		}
		// Other validation errors (malformed token, invalid signature, etc.)
		return result, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Valid = true
		result.Claims = claims

		// Extract user ID from sub claim
		if sub, exists := claims["sub"]; exists {
			if userID, ok := sub.(string); ok {
				result.UserID = userID
			}
		}

		return result, nil
	}

	// Token is not valid
	return result, nil
}

func (a *jwtUsecaseImpl) UpsertUser(gothUser goth.User, role ...models.Role) error {
	ctx := context.Background()

	// Set default role if none provided
	userRole := models.RoleUser
	if len(role) > 0 {
		userRole = role[0]
	}

	// Check if auth method already exists
	existingAuthMethod, err := a.authRepo.GetAuthMethodByProviderAndID(ctx, gothUser.Provider, gothUser.UserID)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to check existing auth method: %w", err)
	}

	var auth *db.Auth

	if existingAuthMethod != nil {
		// User exists, get the auth record
		if existingAuthMethod.AuthID != nil {
			auth, err = a.authRepo.GetAuthByID(ctx, *existingAuthMethod.AuthID)
			if err != nil {
				return fmt.Errorf("failed to get existing auth: %w", err)
			}
		}

		// Update the auth method with new tokens
		var expiresAt pgtype.Timestamptz
		if !gothUser.ExpiresAt.IsZero() {
			expiresAt = pgtype.Timestamptz{
				Time:  gothUser.ExpiresAt,
				Valid: true,
			}
		}

		updateParams := db.UpdateAuthMethodParams{
			AuthID:       existingAuthMethod.AuthID,
			Provider:     gothUser.Provider,
			AccessToken:  utils.StringToPtr(gothUser.AccessToken),
			RefreshToken: utils.StringToPtr(gothUser.RefreshToken),
			IDToken:      utils.StringToPtr(gothUser.IDToken),
			ExpiresAt:    expiresAt,
		}

		_, err = a.authRepo.UpdateAuthMethod(ctx, updateParams)
		if err != nil {
			return fmt.Errorf("failed to update auth method: %w", err)
		}
	} else {
		// Create new auth record
		auth, err = a.authRepo.CreateAuth(ctx,
			utils.StringToPtr(gothUser.Email), // username
			nil,                               // password (nil for OAuth users)
			utils.StringToPtr(gothUser.Email), // email
			string(userRole),                  // role
			true,                              // active
		)
		if err != nil {
			return fmt.Errorf("failed to create auth: %w", err)
		}

		// Create auth method
		authMethod := utils.GothUserToAuthMethod(gothUser, auth.ID)

		createParams := db.CreateAuthMethodParams{
			AuthID:            authMethod.AuthID,
			Provider:          authMethod.Provider,
			ProviderID:        authMethod.ProviderID,
			Email:             authMethod.Email,
			UserID:            authMethod.UserID,
			Name:              authMethod.Name,
			FirstName:         authMethod.FirstName,
			LastName:          authMethod.LastName,
			NickName:          authMethod.NickName,
			Description:       authMethod.Description,
			AvatarUrl:         authMethod.AvatarUrl,
			Location:          authMethod.Location,
			AccessToken:       authMethod.AccessToken,
			RefreshToken:      authMethod.RefreshToken,
			IDToken:           authMethod.IDToken,
			ExpiresAt:         authMethod.ExpiresAt,
			AccessTokenSecret: authMethod.AccessTokenSecret,
		}

		_, err = a.authRepo.CreateAuthMethod(ctx, createParams)
		if err != nil {
			return fmt.Errorf("failed to create auth method: %w", err)
		}
	}

	return nil
}
