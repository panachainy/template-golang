package repositories

import (
	"template-golang/database"
	"template-golang/modules/cockroach/entities"

	"github.com/labstack/gommon/log"
)

type cockroachPostgresRepository struct {
	db database.Database
}

func NewCockroachPostgresRepository(db database.Database) CockroachRepository {
	return &cockroachPostgresRepository{db: db}
}

func (r *cockroachPostgresRepository) InsertCockroachData(in *entities.InsertCockroachDto) error {
	data := &entities.Cockroach{
		Amount: in.Amount,
	}

	result := r.db.GetDb().Create(data)

	if result.Error != nil {
		log.Errorf("InsertCockroachData: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertCockroachData: %v", result.RowsAffected)
	return nil
}
