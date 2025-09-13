package repositories

import (
	"context"
	"math"
	db "template-golang/db/sqlc"
	"template-golang/modules/cockroach/entities"
	"template-golang/pkg/errors"
	"template-golang/pkg/logger"
)

type cockroachPostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) CockroachRepository {
	return &cockroachPostgresRepository{queries: queries}
}

func (r *cockroachPostgresRepository) InsertCockroachData(ctx context.Context, in *entities.InsertCockroachDto) (*entities.Cockroach, error) {
	if in.Amount > math.MaxInt32 {
		return nil, errors.BadRequest("amount exceeds maximum allowed value")
	}

	cockroach, err := r.queries.CreateCockroach(ctx, int32(in.Amount))
	if err != nil {
		logger.Errorf("InsertCockroachData: %v", err)
		return nil, err
	}

	if cockroach.ID < 0 || cockroach.Amount < 0 {
		return nil, errors.Internal("invalid negative values returned from database")
	}

	result := &entities.Cockroach{
		Id:        uint32(cockroach.ID),
		Amount:    uint32(cockroach.Amount),
		CreatedAt: cockroach.CreatedAt.Time,
	}

	logger.Debugf("InsertCockroachData: created cockroach with ID %d", cockroach.ID)
	return result, nil
}

func (r *cockroachPostgresRepository) GetCockroachByID(ctx context.Context, id uint32) (*entities.Cockroach, error) {
	if id > math.MaxInt32 {
		return nil, errors.BadRequest("id exceeds maximum allowed value")
	}

	cockroach, err := r.queries.GetCockroachByID(ctx, int32(id))
	if err != nil {
		logger.Errorf("GetCockroachByID: %v", err)
		return nil, err
	}

	if cockroach.ID < 0 || cockroach.Amount < 0 {
		return nil, errors.Internal("invalid negative values returned from database")
	}

	result := &entities.Cockroach{
		Id:        uint32(cockroach.ID),
		Amount:    uint32(cockroach.Amount),
		CreatedAt: cockroach.CreatedAt.Time,
	}

	return result, nil
}

func (r *cockroachPostgresRepository) ListCockroaches(ctx context.Context) ([]*entities.Cockroach, error) {
	cockroaches, err := r.queries.ListCockroaches(ctx)
	if err != nil {
		logger.Errorf("ListCockroaches: %v", err)
		return nil, err
	}

	var result []*entities.Cockroach
	for _, c := range cockroaches {
		if c.ID < 0 || c.Amount < 0 {
			return nil, errors.Internal("invalid negative values returned from database")
		}

		result = append(result, &entities.Cockroach{
			Id:        uint32(c.ID),
			Amount:    uint32(c.Amount),
			CreatedAt: c.CreatedAt.Time,
		})
	}

	return result, nil
}
