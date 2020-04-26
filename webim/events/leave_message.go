package events

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type LeaveMessageEvent struct {
	*Event
	Body *dto.LeaveMessage `json:"body"`
}

func NewLeaveMessageEvent(agent *models.Agent, msg *models.LeaveMessage, visitInfo *dto.VisitInfo) *LeaveMessageEvent {
	return &LeaveMessageEvent{
		Event: &Event{
			ID:            common.GenUniqueID(),
			Action:        LeaveMessageUpdate,
			AgentID:       agent.ID,
			AgentNickname: agent.NickName,
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  agent.EntID,
			RealName:      agent.RealName,
			TargetID:      msg.ID,
			TargetKind:    "leave_message",
			TraceStart:    -1,
			TrackID:       msg.TrackID,
		},
		Body: dto.ModelToLeaveMessage(msg, visitInfo),
	}
}
