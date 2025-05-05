package models

type AddCockroachData struct {
	Amount uint32 `json:"amount" validate:"required,gt=0"`
}
