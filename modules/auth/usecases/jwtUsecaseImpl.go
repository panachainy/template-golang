package usecases

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"template-golang/config"
	"template-golang/modules/auth/entities"
	"template-golang/modules/auth/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth"
)

type jwtUsecaseImpl struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	repo       repositories.AuthRepository
}

func Provide(conf *config.Config, repo repositories.AuthRepository) *jwtUsecaseImpl {
	privateKey := loadPrivateKey(conf.Auth.PrivateKeyPath)
	publicKey := &privateKey.PublicKey

	return &jwtUsecaseImpl{
		privateKey: privateKey,
		publicKey:  publicKey,
		repo:       repo,
	}
}

// we use panic because if not have private key, we cannot run the server
func loadPrivateKey(path string) *ecdsa.PrivateKey {
	// init key
	var keyByteArray []byte
	var key *ecdsa.PrivateKey

	// Load the private key from a file
	keyByteArray, err := os.ReadFile(path)
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

func (a *jwtUsecaseImpl) ValidateJWT(tokenString string) (*TokenValidationResult, error) {
	result := &TokenValidationResult{
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

func (a *jwtUsecaseImpl) UpsertUser(user goth.User) error {
	if err := a.repo.UpsertData(&entities.Auth{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
		// Username: , // FIXME: gen by system if empty

		AvatarURL: user.AvatarURL,
		Location:  user.Location,

		// RawData: user.RawData,

		FirstName: user.FirstName,
		LastName:  user.LastName,
		NickName:  user.NickName,

		AuthMethods: []entities.AuthMethod{
			// TODO: when we have multiple auth methods, we need to handle it
			{
				Provider:          entities.Provider(user.Provider),
				ProviderID:        "goth_" + user.Provider,
				AccessToken:       user.AccessToken,
				RefreshToken:      user.RefreshToken,
				IDToken:           user.IDToken,
				ExpiresAt:         user.ExpiresAt,
				AccessTokenSecret: user.AccessTokenSecret,
			},
		}},
	); err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	// Return nil to indicate success
	return nil
}
