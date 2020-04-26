package events

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AttrsChangeBody struct {
	Attrs map[string]interface{} `json:"attrs"`
}

type AttrsChange struct {
	*Event
	Body *AttrsChangeBody `json:"body"`
}

func NewAttrsChange(agent *models.Agent, trackID string, attrs map[string]interface{}) *AttrsChange {
	return &AttrsChange{
		Event: &Event{
			ID:            common.GenUniqueID(),
			Action:        AgentChangeAtrr,
			AgentID:       "",
			AgentNickname: "",
			RealName:      "",
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  agent.EntID,
			TargetID:      "None",
			TargetKind:    "client",
			TraceStart:    float64(time.Now().Unix()),
			TrackID:       trackID,
		},
		Body: &AttrsChangeBody{Attrs: attrs},
	}
}
