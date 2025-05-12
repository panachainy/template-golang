package usecases

import (
	"template-golang/config"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/line"
)

type authUsecaseImpl struct {
}

func Provide(conf *config.Config) *authUsecaseImpl {

	goth.UseProviders(
		// line.New(os.Getenv("LINE_KEY"), os.Getenv("LINE_SECRET"), "http://localhost:3000/auth/line/callback", "profile", "openid", "email"),
		line.New(conf.Auth.Line.ClientID, conf.Auth.Line.ClientSecret, "http://localhost:3000/auth/line/callback", "profile", "openid", "email"),

		// facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		// google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),

		// instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
		// apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), "http://localhost:3000/auth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail),
	)
	return &authUsecaseImpl{}
}
