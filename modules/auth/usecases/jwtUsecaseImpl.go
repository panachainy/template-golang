package usecases

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"template-golang/config"

	"github.com/golang-jwt/jwt/v5"
)

type jwtUsecaseImpl struct {
	key *ecdsa.PrivateKey
}

func Provide(conf *config.Config) *jwtUsecaseImpl {
	key := loadPrivateKey(os.Getenv("PRIVATE_KEY_PATH"))

	return &jwtUsecaseImpl{
		key: key,
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

func (a *jwtUsecaseImpl) generateJWT() (string, error) {
	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "my-auth-server",
		"sub": "john",
		"foo": 2,
	})

	// Sign the token with the private key
	signedString, err := token.SignedString(a.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedString, nil
}
