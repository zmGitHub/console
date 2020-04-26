package events

type InviteEval struct {
	*Event
	Body interface{} `json:"body"`
}
