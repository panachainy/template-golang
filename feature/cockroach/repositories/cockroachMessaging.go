package repositories

import "template-golang/feature/cockroach/entities"

type CockroachMessaging interface {
	PushNotification(m *entities.CockroachPushNotificationDto) error
}
