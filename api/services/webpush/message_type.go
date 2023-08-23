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

type TExchSessionKeyData struct {
	SessionID         int    `json:"id"`
	SeesionName       string `json:"sessionName"`
	NumOfParticipants int    `json:"numOfParticipants"`
	Invetee           string `json:"inviteUserName"`
	Key               string `json:"key"`
}
