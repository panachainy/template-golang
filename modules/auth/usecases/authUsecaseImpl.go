package usecases

type authUsecaseImpl struct {
}

func Provide() *authUsecaseImpl {
	return &authUsecaseImpl{}
}
