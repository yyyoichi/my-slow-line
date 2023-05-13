package webpush_services

import "github.com/SherClockHolmes/webpush-go"

type WebpushService struct{}

func (wp *WebpushService) SendNotification(suscription *webpush.Subscription, message []byte) error {
	return nil
}
