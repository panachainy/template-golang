package entities

import (
	"template-golang/modules/auth/models"
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
	ID        string `gorm:"primaryKey;type:varchar(36);default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	AuthID     string   `gorm:"index;type:varchar(36)" json:"auth_id"` // Foreign key to Auth
	Provider   Provider `json:"provider"`
	ProviderID string   `json:"provider_id"` // ID from the SSO provider

	Email string
	// The UserID field is populated by the specific provider’s implementation of the FetchUser method **goth**
	UserID      string `json:"user_id,omitempty"` // populated by the provider
	Name        string `json:"name,omitempty"`    // Optional for SSO
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	NickName    string `json:"nick_name,omitempty"`
	Description string `json:"description,omitempty"` // Optional for SSO
	AvatarURL   string `json:"avatar_url,omitempty"`  // Optional for SSO
	Location    string `json:"location,omitempty"`    // Optional for SSO

	// FIXME: split to new table for support multiple login. **maybe
	AccessToken       string    `json:"access_token"`                  // SSO provider access token
	RefreshToken      string    `json:"refresh_token"`                 // SSO provider refresh token
	IDToken           string    `json:"id_token"`                      // SSO provider ID token
	ExpiresAt         time.Time `json:"expires_at,omitempty"`          // Optional for SSO
	AccessTokenSecret string    `json:"access_token_secret,omitempty"` // Optional for SSO

	// RawData datatypes.JSON `json:"raw_data,omitempty"`
}

// Auth represents user authentication data
type Auth struct {
	ID        string `gorm:"primaryKey;type:varchar(36);default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username    string       `gorm:"uniqueIndex" json:"username,omitempty"`
	Password    string       `json:"password,omitempty"` // Optional for SSO
	Email       string       `gorm:"uniqueIndex" json:"email,omitempty"`
	Role        models.Role  `gorm:"not null" json:"role" validate:"required"`
	Active      bool         `gorm:"not null;default:true" json:"active"`
	AuthMethods []AuthMethod `gorm:"foreignKey:AuthID" json:"auth_methods"` // Support multiple login methods
}
