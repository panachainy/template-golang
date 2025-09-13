package repositories

import (
	"context"
	db "template-golang/db/sqlc"
)

type AuthRepository interface {
	GetAuthByID(ctx context.Context, id string) (*db.Auth, error)
	GetAuthByUsername(ctx context.Context, username string) (*db.Auth, error)
	GetAuthByEmail(ctx context.Context, email string) (*db.Auth, error)
	CreateAuth(ctx context.Context, username *string, password *string, email *string, role string, active bool) (*db.Auth, error)
	UpdateAuth(ctx context.Context, params db.UpdateAuthParams) (*db.Auth, error)
	SoftDeleteAuth(ctx context.Context, id string) error
	ListAllAuths(ctx context.Context) ([]*db.Auth, error)
	CreateAuthMethod(ctx context.Context, params db.CreateAuthMethodParams) (*db.AuthMethod, error)
	GetAuthMethodByProviderAndID(ctx context.Context, provider string, providerID string) (*db.AuthMethod, error)
	GetAuthMethodsByAuthID(ctx context.Context, authID string) ([]*db.AuthMethod, error)
	UpdateAuthMethod(ctx context.Context, params db.UpdateAuthMethodParams) (*db.AuthMethod, error)
	SoftDeleteAuthMethod(ctx context.Context, id string) error
}

type authRepository struct {
	queries *db.Queries
}

func NewAuthRepository(queries *db.Queries) AuthRepository {
	return &authRepository{
		queries: queries,
	}
}

func (r *authRepository) GetAuthByID(ctx context.Context, id string) (*db.Auth, error) {
	auth, err := r.queries.GetAuthByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetAuthByUsername(ctx context.Context, username string) (*db.Auth, error) {
	auth, err := r.queries.GetAuthByUsername(ctx, &username)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetAuthByEmail(ctx context.Context, email string) (*db.Auth, error) {
	auth, err := r.queries.GetAuthByEmail(ctx, &email)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) CreateAuth(ctx context.Context, username *string, password *string, email *string, role string, active bool) (*db.Auth, error) {
	auth, err := r.queries.CreateAuth(ctx, username, password, email, role, active)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) UpdateAuth(ctx context.Context, params db.UpdateAuthParams) (*db.Auth, error) {
	auth, err := r.queries.UpdateAuth(ctx, params)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) SoftDeleteAuth(ctx context.Context, id string) error {
	return r.queries.SoftDeleteAuth(ctx, id)
}

func (r *authRepository) ListAllAuths(ctx context.Context) ([]*db.Auth, error) {
	auths, err := r.queries.ListAllAuths(ctx)
	if err != nil {
		return nil, err
	}

	var result []*db.Auth
	for _, auth := range auths {
		authCopy := auth
		result = append(result, &authCopy)
	}

	return result, nil
}

func (r *authRepository) CreateAuthMethod(ctx context.Context, params db.CreateAuthMethodParams) (*db.AuthMethod, error) {
	authMethod, err := r.queries.CreateAuthMethod(ctx, params)
	if err != nil {
		return nil, err
	}
	return &authMethod, nil
}

func (r *authRepository) GetAuthMethodByProviderAndID(ctx context.Context, provider string, providerID string) (*db.AuthMethod, error) {
	authMethod, err := r.queries.GetAuthMethodByProviderAndID(ctx, provider, providerID)
	if err != nil {
		return nil, err
	}
	return &authMethod, nil
}

func (r *authRepository) GetAuthMethodsByAuthID(ctx context.Context, authID string) ([]*db.AuthMethod, error) {
	authMethods, err := r.queries.GetAuthMethodsByAuthID(ctx, &authID)
	if err != nil {
		return nil, err
	}

	var result []*db.AuthMethod
	for _, method := range authMethods {
		methodCopy := method
		result = append(result, &methodCopy)
	}

	return result, nil
}

func (r *authRepository) UpdateAuthMethod(ctx context.Context, params db.UpdateAuthMethodParams) (*db.AuthMethod, error) {
	authMethod, err := r.queries.UpdateAuthMethod(ctx, params)
	if err != nil {
		return nil, err
	}
	return &authMethod, nil
}

func (r *authRepository) SoftDeleteAuthMethod(ctx context.Context, id string) error {
	return r.queries.SoftDeleteAuthMethod(ctx, id)
}
