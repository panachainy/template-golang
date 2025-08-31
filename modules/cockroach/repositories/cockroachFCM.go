package repositories

import (
	"template-golang/modules/cockroach/entities"
	"template-golang/pkg/logger"
)

type cockroachFCMMessaging struct{}

func ProvideFCMMessaging() CockroachMessaging {
	return &cockroachFCMMessaging{}
}

func (c *cockroachFCMMessaging) PushNotification(m *entities.CockroachPushNotificationDto) error {
	// ... handle logic to push FCM notification here ...
	logger.Debugf("Pushed FCM notification with data: %v", m)
	return nil
}
