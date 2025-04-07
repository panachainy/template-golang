package usecases

import "template-golang/modules/cockroach/models"

type CockroachUsecase interface {
	CockroachDataProcessing(in *models.AddCockroachData) error
}
