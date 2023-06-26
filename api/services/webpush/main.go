package webpush_services

import (
	"himakiwa/services/database"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pkg/errors"
)

type WebpushService struct{}

func (wp *WebpushService) SendNotification(userId int, message string) error {
	w := database.WebpushRepository{}
	result, err := w.QueryByUserId(userId)
	if err != nil {
		return err
	}
	if len(result) < 1 {
		return errors.New("do not subscription")
	}
	s := &webpush.Subscription{Endpoint: result[0].Endpoint}
	s.Keys.Auth = result[0].Auth
	s.Keys.P256dh = result[0].P256dh

	// Send Notification
	resp, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
		Subscriber:      os.Getenv("EMAIL_ADDRESS"), // Do not include "mailto:"
		VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (wp *WebpushService) Subscription(userId int, endpoint, p256dh, auth, userAgent string, expTime *time.Time) error {
	w := database.WebpushRepository{}
	result, err := w.QueryByUserId(userId)
	if err != nil {
		return err
	}
	// delete old
	if len(result) > 0 {
		if err = w.DeleteAll((userId)); err != nil {
			return err
		}
	}

	// insert new webpush
	err = w.Create(userId, endpoint, p256dh, auth, userAgent, expTime)
	return err
}
