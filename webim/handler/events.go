package handler

import (
	"sync"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var (
	AgentStatusUpdate = "agent_update"

	QuickReplyGroupCreate = "quick_reply_group_create"
	QuickReplyGroupUpdate = "quick_reply_group_update"
	QuickReplyGroupDelete = "quick_reply_group_delete"

	QuickReplyCreate = "quick_reply_create"
	QuickReplyUpdate = "quick_reply_update"
	QuickReplyDelete = "quick_reply_delete"

	ClientTagCreate = "client_tag_create"
	ClientTagUpdate = "client_tag_update"
	ClientTagDelete = "client_tag_delete"

	EndConvAgent   = "end_conv_agent"
	EndConvTimeout = "end_conv_timeout"
	EndConvOffline = "end_conv_offline"

	ClientAddTag    = "client_add_tag"
	ClientRemoveTag = "client_remove_tag"

	ClientEvaluation = "client_evaluation"

	VisitBlackAdd    = "visit_black_add"
	VisitBlackRemove = "visit_black_del"

	VisitorOnline  = "online"
	VisitorOffline = "offline"
)

type Event struct {
	Action        string      `json:"action"`
	AgentID       string      `json:"agent_id"`
	AgentNickname string      `json:"agent_nickname"`
	Body          interface{} `json:"body"`
	CreatedOn     string      `json:"created_on"`
	EnterpriseID  string      `json:"enterprise_id"`
	ID            string      `json:"id"`
	RealName      string      `json:"realname"`
	Source        string      `json:"source"`
	TargetID      string      `json:"target_id"`
	TargetKind    string      `json:"target_kind"`
	TraceStart    float64     `json:"trace_start"`
	TrackID       string      `json:"track_id"`
}

type EndConversationEvent struct {
	AgentMsgNum  int    `json:"agent_msg_num"`
	ClientMsgNum int    `json:"client_msg_num"`
	ConvID       string `json:"conv_id"`
	EndedBy      string `json:"ended_by"`
	Evaluation   bool   `json:"evaluation"`
	MsgNum       int    `json:"msg_num"`
}

type tagEventBody struct {
	Color   string `json:"color"`
	Name    string `json:"name"`
	TagID   string `json:"tag_id"`
	TrackID string `json:"track_id"`
}

type evaluationEventBody struct {
	AgentID string `json:"agent_id"`
	Content string `json:"content"`
	EvaType string `json:"eva_type"`
	Level   int    `json:"level"`
}

type visitorStatusUpdateEventBody struct {
	ResidenceTimeSec int `json:"residence_time_sec"`
}

type VisitPageEventBody struct {
	VisitPage *models.VisitPage `json:"visit_page"`
}

type AgentRedirectEventBody struct {
	ConversationBody *adapter.Conversation `json:"conversation_body"`
	ConversationID   string                `json:"conversation_id"`
	From             *adapter.Agent        `json:"from"`
	RedirectedBy     *adapter.Agent        `json:"redirected_by"`
	To               *adapter.Agent        `json:"to"`
}

type agentKickedEventBody struct {
	Message string   `json:"message"`
	Token   []string `json:"token"`
}

// {\"action\":\"update_perms\",\"enterprise_id\":38523,\"id\":\"dc83d3652dae71190633553d9e526bfb\",\"trace_start\":1554534690.732995}"
type PermUpdateEvent struct {
	Action       string  `json:"action"`
	EnterpriseID string  `json:"enterprise_id"`
	ID           string  `json:"id"`
	TraceStart   float64 `json:"trace_start"`
}

var agentKickedEvent = func(agent *models.Agent, body *agentKickedEventBody) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "agent_kicked",
		AgentID:       agent.ID,
		AgentNickname: agent.NickName,
		RealName:      agent.RealName,
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  agent.EntID,
		TargetID:      agent.ID,
		TargetKind:    "agent",
		TraceStart:    -1,
		TrackID:       "",
		Body:          body,
	}
}

var VisitComeEvent = func(body *adapter.VisitInfo) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "visitor_come",
		AgentID:       "",
		AgentNickname: "",
		RealName:      "",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  body.EnterpriseID,
		TargetID:      body.ID,
		TargetKind:    "visit",
		TraceStart:    -1,
		TrackID:       body.TrackID,
		Body:          body,
	}
}

var VisitPageEvent = func(trackID string, body *VisitPageEventBody) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "visit_page",
		AgentID:       "",
		AgentNickname: "",
		RealName:      "",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  body.VisitPage.EntID,
		TargetID:      body.VisitPage.VisitID,
		TargetKind:    "visit",
		TraceStart:    -1,
		TrackID:       trackID,
		Body:          body,
	}
}

var AgentRedirectEvent = func(body *AgentRedirectEventBody) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "agent_redirect",
		AgentID:       body.From.ID,
		AgentNickname: body.From.Nickname,
		RealName:      body.From.Realname,
		Source:        "web",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  body.From.EnterpriseID,
		TargetID:      body.ConversationID,
		TargetKind:    "conv",
		TraceStart:    float64(body.ConversationBody.CreatedOn.UnixNano()),
		TrackID:       body.ConversationBody.TrackID,
		Body:          body,
	}
}

var CreateEvalEvent = func(action, convID, traceID string, agent *models.Agent, e *evaluationEventBody) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agent.ID,
		AgentNickname: agent.NickName,
		RealName:      agent.RealName,
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  agent.EntID,
		TargetID:      convID,
		TargetKind:    "conv",
		TraceStart:    -1,
		TrackID:       traceID,
		Body:          e,
	}
}

var UseClientTagEvent = func(action string, agent *models.Agent, conv *models.Conversation, tag *models.VisitorTag) *Event {
	tagEvent := &tagEventBody{
		Color:   tag.Color,
		Name:    tag.Name,
		TagID:   tag.ID,
		TrackID: conv.TraceID,
	}

	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agent.ID,
		AgentNickname: agent.NickName,
		RealName:      agent.RealName,
		Source:        "web",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  conv.EntID,
		TargetID:      conv.ID,
		TargetKind:    "conv",
		TraceStart:    -1,
		TrackID:       conv.TraceID,
		Body:          tagEvent,
	}
}

func (s *IMService) EndConvEvent(action string, agent *models.Agent, conv *models.Conversation) *Event {
	var autoInvitation bool
	configs, errMsg := s.getEnterpriseConfigs(agent.EntID)
	if errMsg == nil && configs != nil {
		if configs.ServiceEvaluationConfig != nil {
			autoInvitation = configs.ServiceEvaluationConfig.AutoInvitation == "open"
		}
	}
	event := &EndConversationEvent{
		AgentMsgNum:  int(conv.AgentMsgCount),
		ClientMsgNum: int(conv.ClientMsgCount),
		ConvID:       conv.ID,
		EndedBy:      conv.EndedBy.String,
		Evaluation:   autoInvitation,
		MsgNum:       int(conv.MsgCount),
	}
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agent.ID,
		AgentNickname: agent.NickName,
		RealName:      agent.RealName,
		Source:        "web",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  conv.EntID,
		TargetID:      conv.ID,
		TargetKind:    "conv",
		TraceStart:    -1,
		TrackID:       conv.TraceID,
		Body:          event,
	}
}

var AgentStatusUpdateEvent = func(agent *adapter.Agent) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        AgentStatusUpdate,
		AgentID:       agent.ID,
		AgentNickname: agent.Nickname,
		RealName:      agent.Realname,
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  agent.EnterpriseID,
		TargetID:      agent.ID,
		TargetKind:    "agent",
		TraceStart:    -1,
		TrackID:       "",
		Body:          agent,
	}
}

var QkReplyGroupEvent = func(action, agentID, agentNickName, agentRealName string, group *adapter.QkReplyGroup) *Event {
	var t string
	if !group.CreatedOn.IsZero() {
		t = *common.ConvertUTCToTimeString(group.CreatedOn)
	}

	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agentID,
		AgentNickname: agentNickName,
		RealName:      agentRealName,
		CreatedOn:     t,
		EnterpriseID:  group.EnterpriseID,
		TargetID:      group.ID,
		TargetKind:    "quick_reply_group",
		TraceStart:    -1,
		TrackID:       "",
		Body:          group,
	}
}

var QuickReplyEvent = func(action, agentID, agentNickName, agentRealName string, reply *adapter.QkReply) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agentID,
		AgentNickname: agentNickName,
		RealName:      agentRealName,
		CreatedOn:     *common.ConvertUTCToTimeString(reply.CreatedOn),
		EnterpriseID:  reply.EnterpriseID,
		TargetID:      reply.ID,
		TargetKind:    "quick_reply",
		TraceStart:    -1,
		TrackID:       "",
		Body:          reply,
	}
}

var ClientTagEvent = func(action, agentID, agentNickName, agentRealName string, tag *adapter.ClientTag) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agentID,
		AgentNickname: agentNickName,
		RealName:      agentRealName,
		CreatedOn:     *common.ConvertUTCToTimeString(tag.CreatedOn),
		EnterpriseID:  tag.EnterpriseID,
		TargetID:      tag.ID,
		TargetKind:    "client_tag",
		TraceStart:    -1,
		TrackID:       "",
		Body:          tag,
	}
}

var InviteClientEvent = func(agentID, entID, nickName, realName, visitID, trackID string, traceStart float64) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "agent_visit_inviting",
		AgentID:       agentID,
		AgentNickname: nickName,
		RealName:      realName,
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  entID,
		TargetID:      visitID,
		TargetKind:    "visit",
		TraceStart:    traceStart,
		TrackID:       trackID,
		Body: struct {
			Token string `json:"token"`
		}{
			Token: agentID,
		},
	}
}

// {
//  "action": "visit_reject_invite",
//  "agent_id": null,
//  "agent_nickname": null,
//  "body": {},
//  "created_on": "2019-03-31T09:44:38.211025",
//  "enterprise_id": 5869,
//  "id": 2537885169,
//  "target_id": "1JDMdCgJJwtvaEps2L7gCyGBMEi",
//  "target_kind": "visit",
//  "trace_start": 1554025478.236056,
//  "track_id": "1IP6L3kfOGgSDP5XbiIeonmTsCw"
//}
var InviteRejectEvent = func(entID, visitID, trackID string, traceStart float64) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        "visit_reject_invite",
		AgentID:       "",
		AgentNickname: "",
		RealName:      "",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  entID,
		TargetID:      visitID,
		TargetKind:    "visit",
		TraceStart:    traceStart,
		TrackID:       trackID,
		Body:          nil,
	}
}

var VisitorStatusUpdateEvent = func(action, entID, traceID, visitID string, body *visitorStatusUpdateEventBody) *Event {
	return &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       "",
		AgentNickname: "",
		RealName:      "",
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  entID,
		TargetID:      visitID,
		TargetKind:    "visit",
		TraceStart:    -1,
		TrackID:       traceID,
		Body:          body,
	}
}

func (s *IMService) sendEventToAllAgents(entID string, event string) {
	agentIDs, err := models.AgentIDsByEntID(db.Mysql, entID)
	if err != nil {
		log.Logger.Warnf("AgentIDsByEntID error: %v", err)
	}

	if len(agentIDs) > 0 && event != "" {
		var wg sync.WaitGroup
		for _, id := range agentIDs {
			wg.Add(1)
			go func(agentID string) {
				sendMessageToAgent(s.imCli, agentID, event)
				wg.Done()
			}(id)
		}
		wg.Wait()
	}
}

func (s *IMService) sendQkReplyGroupEvent(agentID, action string, group *adapter.QkReplyGroup) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		log.Logger.Warnf("get agent error: %v", err)
		return
	}

	event := QkReplyGroupEvent(action, agentID, agent.NickName, agent.RealName, group)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(group.EnterpriseID, eventContent)
}

func (s *IMService) sendQuickReplyEvent(agentID, action string, reply *adapter.QkReply) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		log.Logger.Warnf("get agent error: %v", err)
		return
	}

	event := QuickReplyEvent(action, agentID, agent.NickName, agent.RealName, reply)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(reply.EnterpriseID, eventContent)
}

func (s *IMService) sendAgentStatusUpdateEvent(agentID, status string, isOnline bool) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		log.Logger.Warnf("get agent error: %v", err)
		return
	}

	agt := adapter.ConvertAgentToAgentInfo(agent, true)
	agt.Status = agent.Status
	agt.IsOnline = isOnline

	event := AgentStatusUpdateEvent(agt)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(agent.EntID, eventContent)
}

func (s *IMService) sendClientTagEvent(agentID, action string, tag *adapter.ClientTag) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		log.Logger.Warnf("get agent error: %v", err)
		return
	}

	event := ClientTagEvent(action, agentID, agent.NickName, agent.RealName, tag)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(agent.EntID, eventContent)
}

func (s *IMService) sendInviteRejectEvent(entID, visitID, trackID string, start float64) {
	event := InviteRejectEvent(entID, visitID, trackID, start)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(entID, eventContent)
}

func (s *IMService) sendEndConvEvent(action, agentID, convID string) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	conv, err := models.ConversationByID(db.Mysql, convID)
	if err != nil {
		return
	}

	agents, err := models.AgentsByConversationID(db.Mysql, convID)
	if err != nil {
		agents = []string{agentID}
	}

	event := s.EndConvEvent(action, agent, conv)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendMessageToMultiAgents(agents, eventContent)
	sendMessageToVisitor(s.imCli, conv.TraceID, conv.EntID, eventContent)
}

func (s *IMService) sendUseClientTagEvent(action, agentID, convID, tagID string) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	conv, err := models.ConversationByID(db.Mysql, convID)
	if err != nil {
		return
	}

	tag, err := models.VisitorTagByID(db.Mysql, tagID)
	if err != nil {
		return
	}

	event := UseClientTagEvent(action, agent, conv, tag)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	sendMessageToAgent(s.imCli, agent.ID, eventContent)
}

func (s *IMService) sendCreateEvalEvent(action, convID string, e *evaluationEventBody) {
	conv, err := models.ConversationByID(db.Mysql, convID)
	if err != nil {
		return
	}

	agent, err := models.AgentByID(db.Mysql, conv.AgentID)
	if err != nil {
		return
	}

	e.AgentID = conv.AgentID
	event := CreateEvalEvent(action, convID, conv.TraceID, agent, e)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	sendMessageToAgent(s.imCli, agent.ID, eventContent)
	sendMessageToVisitor(s.imCli, conv.TraceID, conv.EntID, eventContent)
}

func (s *IMService) sendVisitBlackEvent(action, trackID, visitorID, agentID string, body visitBlack) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	event := &Event{
		ID:            common.GenUniqueID(),
		Action:        action,
		AgentID:       agent.ID,
		AgentNickname: agent.NickName,
		RealName:      agent.RealName,
		CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
		EnterpriseID:  body.EnterpriseID,
		TargetID:      visitorID,
		TargetKind:    "client",
		TraceStart:    -1,
		TrackID:       trackID,
		Body:          body,
	}

	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	sendMessageToAgent(s.imCli, agent.ID, eventContent)
	sendMessageToVisitor(s.imCli, trackID, agent.EntID, eventContent)
}

func (s *IMService) sendAgentRedirectEvent(body *AgentRedirectEventBody) {
	event := AgentRedirectEvent(body)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("[sendAgentRedirectEvent] marshal event error: %v", err)
	}

	s.sendMessageToMultiAgents([]string{body.From.ID, body.To.ID}, eventContent)
	sendMessageToVisitor(s.imCli, body.ConversationBody.TrackID, body.From.EnterpriseID, eventContent)
}

func (s *IMService) sendAgentKickedEvent(agentID string, tokens []string) {
	body := &agentKickedEventBody{Message: "你已经被管理员强制下线!", Token: tokens}
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	event := agentKickedEvent(agent, body)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("[sendAgentRedirectEvent] marshal event error: %v", err)
		return
	}

	s.sendEventToAllAgents(agent.EntID, eventContent)
}
