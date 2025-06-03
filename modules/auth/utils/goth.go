package utils

import (
	"template-golang/modules/auth/entities"

	"github.com/markbates/goth"
)

// TransformGothUser transforms goth.User to auth entities
func GothUserTo(gothUser goth.User) *entities.Auth {
	auth := &entities.Auth{}

	authMethod := &entities.AuthMethod{
		Provider:   entities.Provider(gothUser.Provider),
		ProviderID: "goth_" + gothUser.Provider,

		Email:       gothUser.Email,
		UserID:      gothUser.UserID,
		Name:        gothUser.Name,
		FirstName:   gothUser.FirstName,
		LastName:    gothUser.LastName,
		NickName:    gothUser.NickName,
		Description: gothUser.Description,
		AvatarURL:   gothUser.AvatarURL,
		Location:    gothUser.Location,

		AccessToken:       gothUser.AccessToken,
		RefreshToken:      gothUser.RefreshToken,
		IDToken:           gothUser.IDToken,
		ExpiresAt:         gothUser.ExpiresAt,
		AccessTokenSecret: gothUser.AccessTokenSecret,
	}

	auth.AuthMethods = []entities.AuthMethod{*authMethod}

	return auth
}
