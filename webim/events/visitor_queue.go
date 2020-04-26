package events

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/dto"
)

type QueueingAddBody struct {
	AgentID []string `json:"agent_id"`
	*dto.VisitInfo
}

type QueueingAdd struct {
	*Event
	Body *QueueingAddBody `json:"body"`
}

type QueueRemoveBody struct {
	TrackIDs []string `json:"track_ids"`
}

type QueueingRemove struct {
	*Event
	Body *QueueRemoveBody `json:"body"`
}

func NewQueueingAdd(agentIDs []string, visitorID string, visitInfo *dto.VisitInfo) *QueueingAdd {
	return &QueueingAdd{
		Event: &Event{
			ID:            common.GenUniqueID(),
			Action:        QueueingAddAction,
			AgentID:       "",
			AgentNickname: "",
			RealName:      "",
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  visitInfo.EnterpriseID,
			TargetID:      visitorID,
			TargetKind:    "visitor",
			TraceStart:    float64(time.Now().Unix()),
			TrackID:       visitInfo.TrackID,
		},
		Body: &QueueingAddBody{
			AgentID:   agentIDs,
			VisitInfo: visitInfo,
		},
	}
}

func NewQueueingRemove(trackID, entID string) *QueueingRemove {
	return &QueueingRemove{
		Event: &Event{
			ID:            common.GenUniqueID(),
			Action:        QueueingRemoveAction,
			AgentID:       "",
			AgentNickname: "",
			RealName:      "",
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  entID,
			TargetID:      trackID,
			TargetKind:    "visitor",
			TraceStart:    float64(time.Now().Unix()),
			TrackID:       trackID,
		},
		Body: &QueueRemoveBody{TrackIDs: []string{trackID}},
	}
}
