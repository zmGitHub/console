package handler

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type viewAgent struct {
	*models.Agent
	DeletedAt string `json:"deleted_at"`
}

func modelAgentToViewAgent(agent *models.Agent) (v *viewAgent) {
	v = &viewAgent{
		Agent:     agent,
		DeletedAt: "",
	}

	if agent.DeletedAt.Valid {
		v.DeletedAt = agent.DeletedAt.Time.Format(time.RFC3339)
	}
	return
}

func modelPromotionMsgToView(msg *models.PromotionMsg) *adapter.PromotionMessage {
	return &adapter.PromotionMessage{
		PromotionMsg: msg,
		Content:      msg.Content.String,
		CreatedOn:    common.ConvertUTCToTimeString(msg.CreatedOn),
		UpdatedOn:    common.ConvertUTCToTimeString(msg.UpdatedOn),
	}
}
