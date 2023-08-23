package webpush

type tMessage struct {
	Type tMessageType `json:"type"`
	Data interface{}  `json:"data"`
}

type tMessageType string

const (
	planeMessage          tMessageType = "palne"
	exchSessionKeyMessage tMessageType = "exchSessionKey"
)

type tPlaneData struct {
	Text string `json:"text"`
}

type tExchSessionKeyData struct {
	SessionID int    `json:"id"`
	Key       string `json:"key"`
	Text      string `json:"text"`
}
