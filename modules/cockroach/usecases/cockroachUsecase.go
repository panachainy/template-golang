package usecases

import "template-golang/modules/cockroach/models"

type CockroachUsecase interface {
	ProcessData(data *models.AddCockroachData) error
}
