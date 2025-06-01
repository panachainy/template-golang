//go:generate mockgen -source=authRepository.go -destination=../../../mock/mock_auth_repository.go -package=mock
package repositories

import "template-golang/modules/auth/entities"

type AuthRepository interface {
	UpsertData(in *entities.Auth) error
	Gets(limit int) ([]*entities.Auth, error)
}
