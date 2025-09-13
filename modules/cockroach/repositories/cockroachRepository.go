package repositories

import (
	"context"
	"template-golang/modules/cockroach/entities"
)

type CockroachRepository interface {
	InsertCockroachData(ctx context.Context, in *entities.InsertCockroachDto) (*entities.Cockroach, error)
	GetCockroachByID(ctx context.Context, id uint32) (*entities.Cockroach, error)
	ListCockroaches(ctx context.Context) ([]*entities.Cockroach, error)
}
