package utils

import (
	db "template-golang/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth"
)

// GothUserToAuthMethod transforms goth.User to AuthMethod
func GothUserToAuthMethod(gothUser goth.User, authID string) *db.AuthMethod {
	var expiresAt pgtype.Timestamptz
	if !gothUser.ExpiresAt.IsZero() {
		expiresAt = pgtype.Timestamptz{
			Time:  gothUser.ExpiresAt,
			Valid: true,
		}
	}

	return &db.AuthMethod{
		AuthID:     &authID,
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,

		Email:       StringToPtr(gothUser.Email),
		UserID:      StringToPtr(gothUser.UserID),
		Name:        StringToPtr(gothUser.Name),
		FirstName:   StringToPtr(gothUser.FirstName),
		LastName:    StringToPtr(gothUser.LastName),
		NickName:    StringToPtr(gothUser.NickName),
		Description: StringToPtr(gothUser.Description),
		AvatarUrl:   StringToPtr(gothUser.AvatarURL),
		Location:    StringToPtr(gothUser.Location),

		AccessToken:       StringToPtr(gothUser.AccessToken),
		RefreshToken:      StringToPtr(gothUser.RefreshToken),
		IDToken:           StringToPtr(gothUser.IDToken),
		ExpiresAt:         expiresAt,
		AccessTokenSecret: StringToPtr(gothUser.AccessTokenSecret),
	}
}

// StringToPtr converts string to *string, returns nil for empty strings
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
