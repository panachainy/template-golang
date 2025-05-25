package repositories

import "template-golang/modules/auth/entities"

type AuthRepository interface {
	InsertData(in *entities.Auth) error
}
