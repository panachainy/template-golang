package entities

import (
	"time"

	"gorm.io/gorm"
)

// Provider represents supported SSO providers
type Provider string

const (
	ProviderLocal    Provider = "local"
	ProviderFirebase Provider = "firebase"
	ProviderLine     Provider = "line"
)

// AuthMethod represents a single authentication method
type AuthMethod struct {
	gorm.Model
	// FIXME: check
	AuthID            string    `gorm:"index" json:"auth_id"` // Foreign key to Auth
	Provider          Provider  `json:"provider"`
	ProviderID        string    `json:"provider_id"`                   // ID from the SSO provider
	AccessToken       string    `json:"access_token"`                  // SSO provider access token
	RefreshToken      string    `json:"refresh_token"`                 // SSO provider refresh token
	IDToken           string    `json:"id_token"`                      // SSO provider ID token
	ExpiresAt         time.Time `json:"expires_at,omitempty"`          // Optional for SSO
	AccessTokenSecret string    `json:"access_token_secret,omitempty"` // Optional for SSO
}

// Auth represents user authentication data
type Auth struct {
	gorm.Model
	// FIXME: check =============== AuthMethod.AuthID

	UserID      string       `gorm:"index" json:"user_id"`
	Username    string       `gorm:"uniqueIndex" json:"username"`
	Password    string       `json:"password,omitempty"` // Optional for SSO
	Email       string       `gorm:"uniqueIndex" json:"email"`
	Role        string       `json:"role"`
	Active      bool         `json:"active"`
	AuthMethods []AuthMethod `gorm:"foreignKey:AuthID" json:"auth_methods"` // Support multiple login methods

	Name        string `json:"name,omitempty"` // Optional for SSO
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	NickName    string `json:"nick_name,omitempty"`
	Description string `json:"description,omitempty"` // Optional for SSO
	AvatarURL   string `json:"avatar_url,omitempty"`  // Optional for SSO
	Location    string `json:"location,omitempty"`    // Optional for SSO

	// RawData datatypes.JSON `json:"raw_data,omitempty"`
}
