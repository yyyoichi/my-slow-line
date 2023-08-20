package webpush

type tMessage struct {
	Type tMessageType `json:"type"`
	Data interface{}  `json:"data"`
}

type tMessageType string

const (
	planeMessage tMessageType = "palne"
)

type tPlaneData struct {
	Text string `json:"text"`
}
