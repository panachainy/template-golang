package entities

// Provider represents supported SSO providers
type Provider string

const (
	ProviderLocal    Provider = "local"
	ProviderFirebase Provider = "firebase"
	ProviderLine     Provider = "line"
)

// AuthMethod represents a single authentication method
type AuthMethod struct {
	Provider     Provider `json:"provider"`
	ProviderID   string   `json:"provider_id"`   // ID from the SSO provider
	AccessToken  string   `json:"access_token"`  // SSO provider access token
	RefreshToken string   `json:"refresh_token"` // SSO provider refresh token
}

// UserAuth represents user authentication data
type UserAuth struct {
	ID          string       `json:"id"`
	Username    string       `json:"username"`
	Password    string       `json:"password,omitempty"` // Optional for SSO
	Email       string       `json:"email"`
	Role        string       `json:"role"`
	Active      bool         `json:"active"`
	AuthMethods []AuthMethod `json:"auth_methods"` // Support multiple login methods
}
