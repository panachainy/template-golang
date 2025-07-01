package repositories

import (
	"template-golang/database"
	"template-golang/modules/cockroach/entities"
	"template-golang/pkg/logger"
)

type cockroachPostgresRepository struct {
	db database.Database
}

func ProvidePostgresRepository(db database.Database) *cockroachPostgresRepository {
	return &cockroachPostgresRepository{db: db}
}

func (r *cockroachPostgresRepository) InsertCockroachData(in *entities.InsertCockroachDto) error {
	data := &entities.Cockroach{
		Amount: in.Amount,
	}

	result := r.db.GetDb().Create(data)

	if result.Error != nil {
		logger.Errorf("InsertCockroachData: %v", result.Error)
		return result.Error
	}

	logger.Debugf("InsertCockroachData: %v", result.RowsAffected)
	return nil
}
