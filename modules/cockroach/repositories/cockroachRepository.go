package repositories

import "template-golang/modules/cockroach/entities"

type CockroachRepository interface {
	InsertCockroachData(in *entities.InsertCockroachDto) error
}
