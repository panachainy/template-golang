//go:generate mockgen -source=authUsecase.go -destination=../../../mock/mock_auth_usecase.go -package=mock

package usecases

type AuthUsecase interface {
}
