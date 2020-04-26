package events

import (
	"fmt"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type KickOfflineEventBody struct {
	Message string `json:"message"`
}

type KickOfflineEvent struct {
	*Event
	Body *KickOfflineEventBody `json:"body"`
}

func NewKickOfflineEvent(agent *models.Agent, offlineAgentID, offlineAgentName string) *KickOfflineEvent {
	return &KickOfflineEvent{
		Event: &Event{
			ID:            common.GenUniqueID(),
			Action:        "agent_kicked",
			AgentID:       agent.ID,
			AgentNickname: agent.NickName,
			RealName:      agent.RealName,
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  agent.EntID,
			TargetID:      offlineAgentID,
			TargetKind:    "agent",
			TraceStart:    float64(time.Now().Unix()),
			TrackID:       "",
		},
		Body: &KickOfflineEventBody{
			Message: fmt.Sprintf("%s客服已经被踢下线", offlineAgentName),
		},
	}
}
