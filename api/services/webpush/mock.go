package webpush

import (
	"encoding/json"
	"fmt"

	gwebpush "github.com/SherClockHolmes/webpush-go"
)

func NewWebpushServicesMock() UserWebpushServicesFunc {
	ws := &WebpushServices{notification: &notificationMock{}}
	return func(endpoint, auth, p256dh string) *WebpushServices {
		ws.subscription = &gwebpush.Subscription{
			Endpoint: endpoint,
			Keys: gwebpush.Keys{
				Auth:   auth,
				P256dh: p256dh,
			},
		}
		return ws
	}
}

type notificationMock struct{}

func (ntf *notificationMock) SendNotification(subscription *gwebpush.Subscription, m tMessage) error {
	message, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Println("///////////////////////////////////////")
	fmt.Println("////// SEND WEBPUSH NOTIFICATION //////")
	fmt.Println("///////////////////////////////////////")
	fmt.Printf("// Endpoint: %s\n", subscription.Endpoint)
	fmt.Printf("// Auth: %s\n", subscription.Keys.Auth)
	fmt.Printf("// P256dh: %s\n", subscription.Keys.P256dh)
	fmt.Println("//-----------------------------------//")
	fmt.Printf("// Message: %s\n", message)
	fmt.Println("//-----------------------------------//")
	fmt.Println("///////////////////////////////////////")
	return nil
}
