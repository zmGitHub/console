package handler

import (
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

// CreateOrUpdateMessageBeep
// POST /admin/api/v1/agents/message_beep
func (s *IMService) CreateOrUpdateMessageBeep(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	msgBeep := &models.MessageBeep{}
	if err = ctx.Bind(msgBeep); err != nil {
		return
	}
	msgBeep.AgentID = agentID

	if msgBeep.ClientType != models.MessageBeepWebClientType {
		return invalidParameterResp(ctx, "unsupported client_type")
	}

	if msgBeep.BeepType != models.MessageDesktopBeepType && msgBeep.BeepType != models.MessagePopupBeepType {
		return invalidParameterResp(ctx, "unsupported beep_type")
	}

	if err = models.CreateOrUpdateMessageBeep(db.Mysql, msgBeep); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: msgBeep})
}

// GetMessageBeepConfig
// GET /admin/api/v1/agents/message_beep
func (s *IMService) GetMessageBeepConfig(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	beep, err := models.MessageBeepByAgentID(db.Mysql, agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: beep})
}
