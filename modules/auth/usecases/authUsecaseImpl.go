package usecases

import (
	"template-golang/config"
)

type authUsecaseImpl struct {
}

func Provide(conf *config.Config) *authUsecaseImpl {
	return &authUsecaseImpl{}
}

// var (
// 	key *ecdsa.PrivateKey
// 	t   *jwt.Token
// 	s   string
//   )

//   key = /* Load key from somewhere, for example a file */
//   t = jwt.NewWithClaims(jwt.SigningMethodES256,
// 	jwt.MapClaims{
// 	  "iss": "my-auth-server",
// 	  "sub": "john",
// 	  "foo": 2,
// 	})
//   s = t.SignedString(key)
