package usecases

import (
	"context"
	"template-golang/modules/cockroach/entities"
	"template-golang/modules/cockroach/models"
	"template-golang/modules/cockroach/repositories"
)

type cockroachUsecaseImpl struct {
	cockroachRepository repositories.CockroachRepository
	cockroachMessaging  repositories.CockroachMessaging
}

func NewCockroachUsecaseImpl(
	cockroachRepository repositories.CockroachRepository,
	cockroachMessaging repositories.CockroachMessaging,
) CockroachUsecase {
	return &cockroachUsecaseImpl{
		cockroachRepository: cockroachRepository,
		cockroachMessaging:  cockroachMessaging,
	}
}

func (u *cockroachUsecaseImpl) ProcessData(in *models.AddCockroachData) error {
	ctx := context.Background()

	insertCockroachData := &entities.InsertCockroachDto{
		Amount: in.Amount,
	}

	cockroach, err := u.cockroachRepository.InsertCockroachData(ctx, insertCockroachData)
	if err != nil {
		return err
	}

	pushCockroachData := &entities.CockroachPushNotificationDto{
		Title:        "Cockroach Detected 🪳 !!!",
		Amount:       cockroach.Amount,
		ReportedTime: cockroach.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := u.cockroachMessaging.PushNotification(pushCockroachData); err != nil {
		return err
	}

	return nil
}
