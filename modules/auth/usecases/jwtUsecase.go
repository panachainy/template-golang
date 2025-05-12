//go:generate mockgen -source=jwtUsecase.go -destination=../../../mock/mock_jwt_usecase.go -package=mock

package usecases

type JWTUsecase interface {
}
