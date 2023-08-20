package webpush

import (
	"encoding/json"
	"os"

	gwebpush "github.com/SherClockHolmes/webpush-go"
)

type notificationInterface interface {
	SendNotification(subscription *gwebpush.Subscription, m tMessage) error
}

type notification struct{}

func (ntf *notification) SendNotification(subscription *gwebpush.Subscription, m tMessage) error {
	message, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := gwebpush.SendNotification(message, subscription, &gwebpush.Options{
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
