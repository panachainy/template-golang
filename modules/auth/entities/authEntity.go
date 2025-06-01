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
	AuthID     string   `gorm:"index" json:"auth_id"` // Foreign key to Auth
	Provider   Provider `json:"provider"`
	ProviderID string   `json:"provider_id"` // ID from the SSO provider

	// FIXME: split to new table for support multiple login.
	AccessToken       string    `json:"access_token"`                  // SSO provider access token
	RefreshToken      string    `json:"refresh_token"`                 // SSO provider refresh token
	IDToken           string    `json:"id_token"`                      // SSO provider ID token
	ExpiresAt         time.Time `json:"expires_at,omitempty"`          // Optional for SSO
	AccessTokenSecret string    `json:"access_token_secret,omitempty"` // Optional for SSO
}

// Auth represents user authentication data
type Auth struct {
	gorm.Model
	UserID string `gorm:"index;not null" json:"user_id" validate:"required"`
	// username in this system & must be unique if provided
	// For SSO, this field is optional and can be empty
	// If empty, the system will generate a unique username based on email or provider ID
	Username    string       `gorm:"uniqueIndex" json:"username,omitempty"`
	Password    string       `json:"password,omitempty"` // Optional for SSO
	Email       string       `gorm:"uniqueIndex" json:"email" validate:"required,email"`
	Role        string       `gorm:"not null" json:"role" validate:"required"`
	Active      bool         `gorm:"not null;default:true" json:"active"`
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
