package repositories

import "template-golang/modules/cockroach/entities"

type CockroachMessaging interface {
	PushNotification(m *entities.CockroachPushNotificationDto) error
}
