package usecases

import "template-golang/feature/cockroach/models"

type CockroachUsecase interface {
	CockroachDataProcessing(in *models.AddCockroachData) error
}
