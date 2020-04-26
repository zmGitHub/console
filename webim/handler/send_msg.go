package handler

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
	"bitbucket.org/forfd/custm-chat/webim/imclient"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var (
	visitorChan = `%s_%s`      // {track_id}_{ent_id}
	agentChan   = `%s_message` // {agent_id}_message

	sendMessageToAgent = func(imCli imclient.IMClient, agentID, message string) {
		ch := fmt.Sprintf(agentChan, agentID)
		if err := imCli.PublishMessage(context.Background(), ch, message); err != nil {
			log.Logger.Warnf("[sendMessageToAgent] error: %v", err)
		}
	}

	sendMessageToVisitor = func(imCli imclient.IMClient, trackID, entID string, message string) {
		ch := fmt.Sprintf(visitorChan, trackID, entID)
		if err := imCli.PublishMessage(context.Background(), ch, message); err != nil {
			log.Logger.Warnf("[sendMessageToVisitor] error: %v", err)
		}
	}
)

type ClientSendMsgReq struct {
	EntID          string     `json:"ent_id"`
	TrackID        string     `json:"track_id"`
	ConversationID string     `json:"conversation_id"`
	Type           string     `json:"type"`
	Content        string     `json:"content"`
	BrowserID      string     `json:"browser_id"`
	Extra          *fileExtra `json:"extra"`
}

type ClientMsg struct {
	CreatedOn string `json:"created_on"`
	ID        string `json:"id"`
}

type ClientSendMsgResp struct {
	Msg     *ClientMsg `json:"msg"`
	Success bool       `json:"success"`
}

type fileExtra struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
	ExpireAt string `json:"expire_at"`
}

type AgentSendMsgReq struct {
	Type           string     `json:"type"` // text
	Content        string     `json:"content"`
	ConversationID string     `json:"conversation_id"`
	TrackID        string     `json:"track_id"`
	BrowserID      string     `json:"browser_id"`
	Extra          *fileExtra `json:"extra"`
}

type InternalMsgReq struct {
	Type           string `json:"type"`
	Content        string `json:"content"`
	ConversationID string `json:"conversation_id"`
	TrackID        string `json:"track_id"`
	BrowserID      string `json:"browser_id"`
}

// POST /api/agent/send_internal_msg
func (s *IMService) SendInternalMsg(ctx echo.Context) error {
	req := &InternalMsgReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	agentID := ctx.Get(middleware.AgentIDKey).(string)
	agents, err := models.AgentsByConversationID(db.Mysql, req.ConversationID)
	if err != nil {
		agents = []string{agentID}
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	sendMsgReq := &sendMsgReq{
		EntID:       entID,
		TraceID:     req.TrackID,
		Creator:     agentID,
		Agents:      agents,
		ConvID:      req.ConversationID,
		FromType:    models.MessageFromAgentType,
		Content:     req.Content,
		ContentType: req.Type,
		Internal:    true,
	}

	chanMsg, err := s.sendMsg(sendMsgReq)
	if err != nil {
		return errResp(ctx, common.SendingMsgErr, err.Error())
	}

	monitor.MessagesSentCount.WithLabelValues("internal").Inc()
	return jsonResponse(ctx, chanMsg)
}

// ClientSendMsg ...
// POST /client/send_msg
func (s *IMService) ClientSendMsg(ctx echo.Context) (err error) {
	req := &ClientSendMsgReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.TrackID == "" || req.ConversationID == "" || req.EntID == "" {
		return invalidParameterResp(ctx, "invalid ent_id/track_id/conversation_id")
	}

	conv, err := models.ConversationByID(db.Mysql, req.ConversationID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if conv.EndedAt.Valid {
		return invalidParameterResp(ctx, "对话已结束")
	}

	agents, err := models.AgentsByConversationID(db.Mysql, req.ConversationID)
	if err != nil {
		agents = []string{conv.AgentID}
	}

	sendMsgReq := &sendMsgReq{
		EntID:       req.EntID,
		TraceID:     req.TrackID,
		Creator:     conv.AgentID,
		Agents:      agents,
		ConvID:      conv.ID,
		FromType:    models.MessageFromVisitorType,
		Content:     req.Content,
		ContentType: req.Type,
		Internal:    false,
		extra:       req.Extra,
	}

	chanMsg, err := s.sendMsg(sendMsgReq)
	if err != nil {
		return errResp(ctx, common.SendingMsgErr, err.Error())
	}

	monitor.MessagesSentCount.WithLabelValues("client").Inc()
	return jsonResponse(ctx, &ClientSendMsgResp{
		Msg:     &ClientMsg{CreatedOn: *chanMsg.CreatedOn, ID: chanMsg.ID},
		Success: true,
	})
}

// POST /api/agent/send_msg
func (s *IMService) AgentSendMsgV1(ctx echo.Context) (err error) {
	agentInfo := getAgentInfoFromJwtToken(ctx)
	req := &AgentSendMsgReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.TrackID == "" || req.ConversationID == "" || req.Type == "" || req.Content == "" {
		return invalidParameterResp(ctx, "invalid track_id/conversation_id/type/content")
	}

	end, err := models.IsConversationEnd(db.Mysql, req.ConversationID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if end {
		return invalidParameterResp(ctx, "对话已结束")
	}

	agents, err := models.AgentsByConversationID(db.Mysql, req.ConversationID)
	if err != nil {
		agents = []string{agentInfo.UserID}
	}

	sendMsgReq := &sendMsgReq{
		EntID:       agentInfo.EntID,
		TraceID:     req.TrackID,
		Creator:     agentInfo.UserID,
		Agents:      agents,
		ConvID:      req.ConversationID,
		FromType:    models.MessageFromAgentType,
		Content:     req.Content,
		ContentType: req.Type,
		Internal:    false,
		extra:       req.Extra,
	}

	chanMsg, err := s.sendMsg(sendMsgReq)
	if err != nil {
		return errResp(ctx, common.SendingMsgErr, err.Error())
	}

	monitor.MessagesSentCount.WithLabelValues("agent").Inc()
	return jsonResponse(ctx, &ClientSendMsgResp{
		Msg:     &ClientMsg{CreatedOn: *chanMsg.CreatedOn, ID: chanMsg.ID},
		Success: true,
	})
}
