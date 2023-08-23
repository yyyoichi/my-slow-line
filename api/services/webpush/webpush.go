package webpush

import (
	"fmt"

	gwebpush "github.com/SherClockHolmes/webpush-go"
)

type UserWebpushServices func(endpoint, auth, p256dh string) *WebpushServices
type WebpushServices struct {
	notification notificationInterface
	subscription *gwebpush.Subscription
}

func NewWebpushServices() UserWebpushServices {
	ws := &WebpushServices{notification: &notification{}}
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

// wrap SendNotification method of notification interface
func (ws *WebpushServices) sendNotification(m tMessage) error {
	return ws.notification.SendNotification(ws.subscription, m)
}

func (ws *WebpushServices) SendPlaneMessage(message string) error {
	return ws.sendNotification(tMessage{
		Type: planeMessage,
		Data: tPlaneData{
			Text: message,
		},
	})
}

type TExchSessionKeyArgs struct {
	SessionID         int
	SessionName       string
	NumOfParticipants int
	OfferUserName     string
	Key               string
}

func (ws *WebpushServices) SendExchSessionKeyMessage(data TExchSessionKeyArgs) error {
	text := fmt.Sprintf("%s invited【%s(%d)】", data.OfferUserName, data.SessionName, data.NumOfParticipants)

	return ws.sendNotification(tMessage{
		Type: exchSessionKeyMessage,
		Data: tExchSessionKeyData{
			data.SessionID,
			data.Key,
			text,
		},
	})
}
