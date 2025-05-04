//go:generate mockgen -source=cockroachUsecase.go -destination=../../mock/mock_cockroach_usecase.go -package=mock

package usecases

import "template-golang/modules/cockroach/models"

type CockroachUsecase interface {
	ProcessData(data *models.AddCockroachData) error
}
