package repositories

import (
	"template-golang/modules/cockroach/entities"

	"github.com/labstack/gommon/log"
)

type cockroachFCMMessaging struct{}

func NewCockroachFCMMessaging() *cockroachFCMMessaging {
	return &cockroachFCMMessaging{}
}

func (c *cockroachFCMMessaging) PushNotification(m *entities.CockroachPushNotificationDto) error {
	// ... handle logic to push FCM notification here ...
	log.Debugf("Pushed FCM notification with data: %v", m)
	return nil
}
