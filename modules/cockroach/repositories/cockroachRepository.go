package repositories

import "template-golang/feature/cockroach/entities"

type CockroachRepository interface {
	InsertCockroachData(in *entities.InsertCockroachDto) error
}
