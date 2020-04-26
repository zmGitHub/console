package dto

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type LeaveMessage struct {
	*models.LeaveMessage
	LastOptionAgent *string    `json:"last_option_agent"`
	CreatedAt       *string    `json:"created_at"` // created_at
	UpdatedAt       *string    `json:"updated_at"` // updated_at
	VisitInfo       *VisitInfo `json:"visit_info"`
}

func ModelToLeaveMessage(msg *models.LeaveMessage, visitInfo *VisitInfo) *LeaveMessage {
	var opAgent *string
	if msg.LastOptionAgent.Valid {
		opAgent = &msg.LastOptionAgent.String
	}
	return &LeaveMessage{
		LeaveMessage:    msg,
		LastOptionAgent: opAgent,
		CreatedAt:       common.ConvertUTCToTimeString(msg.CreatedAt),
		UpdatedAt:       common.ConvertUTCToTimeString(msg.UpdatedAt),
		VisitInfo:       visitInfo,
	}
}
