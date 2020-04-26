package handler

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/external/elasticsearch"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type NewConversationReq struct {
	TrackID string `json:"track_id"`
}

type conversationResp struct {
	*models.Conversation
	AgentType             string `json:"agent_type"`
	ClientFirstSendTime   string `json:"client_first_send_time"`
	FirstMsgCreatedAt     string `json:"first_msg_created_at"`
	FirstResponseWaitTime int64  `json:"first_response_wait_time"`
	LastMsgContent        string `json:"last_msg_content"`
	LastMsgCreatedAt      string `json:"last_msg_created_at"`
	QualityGrade          string `json:"quality_grade"`
	EndedAt               string `json:"ended_at"`
	EndedBy               string `json:"ended_by"`
}

type GetHistoryMessagesReq struct {
	TraceID string `query:"trace_id"`
	Offset  int    `query:"offset"`
	Limit   int    `query:"limit"`
}

type GetHistoryMessagesResp struct {
	Conversations []*conversationResp `json:"conversations"`
	VisitInfo     *models.Visit       `json:"visit_info"`
}

type TransferConversationReq struct {
	TraceID     string `json:"trace_id"`
	TargetAgent string `json:"target_agent"`
}

// {"group_token":"72f9de7add336f218a574ab98b28b569","browser_id":"agent1554247101187"}
type TransferConversationReqV1 struct {
	GroupToken string `json:"group_token"`
	BrowserID  string `json:"browser_id"`
}

type TransferToAgentReq struct {
	To string `json:"to"`
}

type EndConversationReq struct {
	EndBy string `json:"end_by"` // agent_id
}

type EndConversationV1Req struct {
	TrackID        string `json:"track_id"`
	ConversationID string `json:"conversation_id"`
	BrowserID      string `json:"browser_id"`
}

type UpdateSummaryReq struct {
	Content string `json:"content"`
}

type timeRange struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

type region struct {
	Province string `json:"province"`
	City     string `json:"city"`
}

type SearchConversationsResp struct {
	TotalCount    int64                         `json:"total_count"`
	Conversations []*adapter.SearchConversation `json:"conversations"`
}

type msg struct {
	CreatedOn string `json:"created_on"`
	ID        string `json:"id"`
}

type InviteEvalResp struct {
	Msg     *msg `json:"msg"`
	Success bool `json:"success"`
}

func (s *IMService) initConversation(entID, traceID, agentID, title string) (conv *models.Conversation, isNew bool, errMsg *ErrMsg) {
	convs, err := models.ConversationsByTraceID(db.Mysql, traceID, 0, 1, true)
	if err != nil {
		return nil, false, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	if len(convs) > 0 {
		return convs[0], false, nil
	}

	now := time.Now().UTC()
	conversation := &models.Conversation{
		ID:             common.GenUniqueID(),
		EntID:          entID,
		TraceID:        traceID,
		AgentID:        agentID,
		CreatedAt:      now,
		UpdateAt:       now,
		AgentMsgCount:  uint(0),
		AgentType:      sql.NullString{String: models.ConversationHumanAgentType, Valid: true},
		MsgCount:       uint(0),
		Title:          title,
		ClientMsgCount: uint(0),
		Duration:       uint(0),
		EvalLevel:      -1,
	}

	if err := conversation.Insert(db.Mysql); err != nil {
		return nil, false, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	if err := db.RedisClient.Incr(fmt.Sprintf(common.AgentConversationNum, conversation.AgentID)).Err(); err != nil {
		log.Logger.Warnf("incr agent: %s, conversation number error: %v", conversation.AgentID, err)
	}

	monitor.ConversationsCreationCount.WithLabelValues("new_conversation").Inc()
	return conversation, true, nil
}

// POST /api/agent/end_conversation
func (s *IMService) EndConversationV1(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor_and_conv", "end_others_conv"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &EndConversationV1Req{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.ConversationID == "" {
		return invalidParameterResp(ctx, "conversation_id is invalid")
	}

	errMsg := s.endConversation(agentID, req.ConversationID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	go s.sendEndConvEvent(EndConvAgent, agentID, req.ConversationID)

	return jsonResponse(ctx, &Resp{Code: 0})
}

func (s *IMService) endConversation(endBy, conversationID string) *ErrMsg {
	dbErrResp := func(dbErr error) *ErrMsg {
		return &ErrMsg{Code: common.DBErr, Message: dbErr.Error()}
	}
	var lastMsgContent string
	var firstMsgCreateTime, firstAgentMsgTime, firstClientMsgSendTime, lastClientMsgSendTime, lastMsgTime time.Time
	lastMsg, dbErr := models.LastMessageByConversationID(db.Mysql, conversationID)
	if dbErr != nil {
		return dbErrResp(dbErr)
	}

	if lastMsg != nil {
		lastMsgContent, lastMsgTime = lastMsg.Content.String, lastMsg.CreatedAt
	}

	firstMsgCreateTime, dbErr = models.FirstConversationMsgCreateTime(db.Mysql, conversationID)
	if dbErr != nil && dbErr != sql.ErrNoRows {
		return dbErrResp(dbErr)
	}

	firstAgentMsgTime, dbErr = models.FirstAgentMsgCreateTime(db.Mysql, conversationID)
	if dbErr != nil && dbErr != sql.ErrNoRows {
		return dbErrResp(dbErr)
	}

	firstClientMsgSendTime, dbErr = models.FirstClientMessageSendTime(db.Mysql, conversationID)
	if dbErr != nil && dbErr != sql.ErrNoRows {
		return dbErrResp(dbErr)
	}

	lastClientMsgSendTime, dbErr = models.LastClientMessageSendTime(db.Mysql, conversationID)
	if dbErr != nil && dbErr != sql.ErrNoRows {
		return dbErrResp(dbErr)
	}

	conv, dbErr := models.ConversationByID(db.Mysql, conversationID)
	if dbErr != nil {
		return dbErrResp(dbErr)
	}
	if conv.EndedAt.Valid {
		return &ErrMsg{Code: common.ConversationEndedErr, Message: "conversation already ended"}
	}

	qualityGrade := s.getQualityLevel(conv.EntID, conv.ID)
	now := time.Now().UTC()
	duration := now.Sub(conv.CreatedAt.UTC()).Seconds()
	var responseWaitTimeInSec int64 = 0
	if !firstAgentMsgTime.IsZero() {
		responseWaitTimeInSec = int64(firstAgentMsgTime.UTC().Sub(conv.CreatedAt.UTC()).Seconds())
	}

	if dbErr = models.EndConversation(
		db.Mysql, conversationID, endBy, lastMsgContent,
		lastMsgTime, int64(duration), responseWaitTimeInSec, qualityGrade,
		firstMsgCreateTime, firstClientMsgSendTime, lastClientMsgSendTime); dbErr != nil {
		return dbErrResp(dbErr)
	}

	convNumKey := fmt.Sprintf(common.AgentConversationNum, conv.AgentID)
	if _, dbErr = db.RedisClient.Decr(convNumKey).Result(); dbErr != nil {
		log.Logger.Warnf("decr: %s, error: %v", convNumKey, dbErr)
		return &ErrMsg{Code: common.RedisErr, Message: dbErr.Error()}
	}

	go s.sendEndMessage(endBy == "system", conv.EntID, conv.TraceID, conversationID, conv.AgentID)

	conv.Duration = uint(duration)
	conv.EndedAt = mysql.NullTime{Time: now, Valid: true}
	conv.EndedBy = sql.NullString{String: endBy, Valid: true}
	conv.QualityGrade.String = qualityGrade
	conv.ClientFirstSendTime = mysql.NullTime{Time: firstClientMsgSendTime, Valid: true}
	conv.ClientLastSendTime = mysql.NullTime{Time: lastClientMsgSendTime, Valid: true}
	conv.LastMsgContent = sql.NullString{String: lastMsgContent, Valid: true}
	conv.LastMsgCreatedAt = mysql.NullTime{Time: lastMsgTime, Valid: true}
	conv.FirstResponseWaitTime = sql.NullInt64{Int64: responseWaitTimeInSec, Valid: true}

	updateConversationStats(conv)
	updateAgentStats(conv)

	go s.DeleteConversationTasks(conv.ID)
	go CreateConversationDoc(conv)
	return nil
}

func (s *IMService) getQualityLevel(entID, convID string) string {
	configs, errMsg := s.getEnterpriseConfigs(entID)
	if errMsg != nil {
		log.Logger.Warnf("getEnterpriseConfigs error: %+v", errMsg)
	}

	if configs == nil {
		return ""
	}

	gradeConf := configs.ConvGradeConfig
	if gradeConf != nil && gradeConf.Enable {
		counts, err := models.MessageCountsByConversationID(db.Mysql, convID)
		if err != nil {
			log.Logger.Warnf("MessageCountsByConversationID error: %v", err)
		} else {
			agentCount, clientCount := counts[models.MessageFromAgentType], counts[models.MessageFromVisitorType]
			if agentCount >= gradeConf.FirstLevel.AgentMsgCnt && clientCount >= gradeConf.FirstLevel.ClientMsgCnt {
				return models.FirstLevel
			}

			if agentCount >= gradeConf.SecondLevel.AgentMsgCnt && clientCount >= gradeConf.SecondLevel.ClientMsgCnt {
				return models.SecondLevel
			}

			if agentCount >= gradeConf.ThirdLevel.AgentMsgCnt && clientCount >= gradeConf.ThirdLevel.ClientMsgCnt {
				return models.ThirdLevel
			}
		}
	}

	return ""
}

// checkQueueVisitor 对话结束以后查看一下队列是不是有等待的访客，有的话发起一个新对话
func (s *IMService) checkQueueVisitor(entID, agentID string) {
	online, err := models.IsAgentOnline(db.Mysql, entID, agentID)
	if err != nil {
		return
	}

	if !online {
		return
	}

	onDuty, err := models.IsAgentOnDuty(db.Mysql, agentID)
	if err != nil {
		return
	}

	if !onDuty {
		return
	}

	queues, err := models.VisitorQueuesByEntID(db.Mysql, entID)
	if err != nil {
		log.Logger.Warnf("[checkQueueVisitor] error: %v", err)
		return
	}

	if len(queues) <= 0 {
		return
	}

	trackID := queues[0].TrackID
	visit, visitor, err := s.getVisitInfoByTrackID(entID, trackID)
	if err != nil {
		log.Logger.Warnf("[checkQueueVisitor] error: %v", err)
		return
	}

	title := fmt.Sprintf("#%d-new-conv", time.Now().Unix())
	_, _, _, errMsg := s.newConversation(entID, trackID, agentID, title, visit, visitor)
	if errMsg != nil {
		log.Logger.Warnf("[checkQueueVisitor] new converation error msg: %v", errMsg)
		return
	}

	queue := &models.VisitorQueue{TrackID: trackID}
	if err := queue.Delete(db.Mysql); err == nil {
		go s.sendVisitorQueueingRemove(entID, trackID)
	}
}

// GetColleagueConversations get colleague conversations
// GET /admin/api/v1/colleague_conversations?offset=0&limit=10
// GET /api/conversation/agent/colleagues
func (s *IMService) GetColleagueConversations(ctx echo.Context) (err error) {
	agentInfo := getAgentInfoFromJwtToken(ctx)
	offset, limit, err := getOffsetLimitFromCtx(ctx)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	permsRange, err := models.GetAgentPermsRange(db.Mysql, agentInfo.UserID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var allConvs bool
	var groupAgents []string
	var resp = &adapter.ColleageConvs{
		Conversations: []*adapter.Conversation{},
		ConvNums:      0,
		FromCache:     false,
		HasNext:       false,
	}
	switch permsRange {
	case models.AgentPermsRangePersonalType:
		return jsonResponse(ctx, resp)
	case models.AgentPermsRangeAllType:
		allConvs = true
	default:
		groupAgents, err = s.getGroupAgents(agentInfo.UserID)
		if err != nil {
			return invalidParameterResp(ctx, err.Error())
		}
	}

	convs, hasNext, err := models.ColleagueConversationsByEntIDAgentID(db.Mysql, agentInfo.EntID, agentInfo.UserID, offset, limit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var traceIDs []string
	if allConvs {
		for _, conv := range convs {
			traceIDs = append(traceIDs, conv.TraceID)
		}
	} else {
	loop:
		for _, conv := range convs {
			for _, id := range groupAgents {
				if id == conv.AgentID {
					traceIDs = append(traceIDs, conv.TraceID)
					continue loop
				}
			}
		}
	}

	if len(traceIDs) == 0 {
		return jsonResponse(ctx, resp)
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	tags, err := models.VisitorTagRelationsByVisitors(db.Mysql, visitors)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var adapterConvs = make([]*adapter.Conversation, len(convs))
	for i, conversation := range convs {
		var visit *models.Visit
		for _, vst := range visits {
			if vst.TraceID == conversation.TraceID {
				visit = vst
				break
			}
		}

		var visitor *models.Visitor
		for _, vstor := range visitors {
			if vstor.TraceID == conversation.TraceID {
				visitor = vstor
				break
			}
		}

		adapterConvs[i] = adapter.ModelConversationToConversation(conversation, visit.ID, visit, visitor, tags)
	}

	resp = &adapter.ColleageConvs{
		Conversations: adapterConvs,
		ConvNums:      len(adapterConvs),
		FromCache:     false,
		HasNext:       hasNext,
	}
	if hasNext {
		var nextPos = offset + limit
		resp.NextCursor = &nextPos
	}

	return jsonResponse(ctx, resp)
}

// GetActiveConversations ...
// GET /admin/api/v1/active_conversations
// GET /api/conversation/agent/active
func (s *IMService) GetActiveConversations(ctx echo.Context) (err error) {
	userID := ctx.Get(middleware.AgentIDKey).(string)
	activeConversations, err := models.ActiveConversationsByAgentID(db.Mysql, userID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var traceIDs []string
	for _, conv := range activeConversations {
		traceIDs = append(traceIDs, conv.TraceID)
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	tagRelations, err := models.VisitorTagRelationsByVisitors(db.Mysql, visitors)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	onlineClients, err := models.OnlineVisitorsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var adapterConvs = make([]*adapter.Conversation, len(activeConversations))
	for i, conversation := range activeConversations {
		var visit *models.Visit
		for _, vst := range visits {
			if vst.TraceID == conversation.TraceID {
				visit = vst
				break
			}
		}

		var visitor *models.Visitor
		for _, vstor := range visitors {
			if vstor.TraceID == conversation.TraceID {
				visitor = vstor
				break
			}
		}

		var visitID string
		if visit != nil {
			visitID = visit.ID
		}

		adapterConv := adapter.ModelConversationToConversation(conversation, visitID, visit, visitor, tagRelations)
		if _, ok := onlineClients[adapterConv.TrackID]; ok {
			adapterConv.IsClientOnline = true
		}
		adapterConvs[i] = adapterConv
	}

	return jsonResponse(ctx, &adapter.ActiveConversationsResp{
		Conversations: adapterConvs,
	})
}

// GetHistoryConversations
// GET /api/v1/enterprises/:ent_id/conversations/history
func (s *IMService) GetHistoryConvs(ctx echo.Context) (err error) {
	entID := ctx.Param("ent_id")
	if entID == "" {
		return invalidParameterResp(ctx, "ent_id is invalid")
	}

	req := &GetHistoryMessagesReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}
	if req.TraceID == "" {
		return invalidParameterResp(ctx, "trace_id is invalid")
	}

	var offset, limit = 0, 5
	if req.Offset > 0 {
		offset = req.Offset
	}

	if req.Limit > 0 {
		limit = req.Limit
	}

	var resp = &GetHistoryMessagesResp{}
	conversations, err := models.ConversationsByEntIDTraceID(db.Mysql, entID, req.TraceID, offset, limit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	for _, conv := range conversations {
		resp.Conversations = append(resp.Conversations, convertModelConversationToResp(conv))
	}

	visits, err := models.VisitsByEntIDTraceID(db.Mysql, entID, req.TraceID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	if len(visits) > 0 {
		resp.VisitInfo = visits[0]
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: resp})
}

// TransferConversation transfer conversation from current agent to another
// POST /api/agent/conversation/:conversation_id/group_redirect
func (s *IMService) TransferConversation(ctx echo.Context) (err error) {
	convID := ctx.Param("conversation_id")
	agentInfo := getAgentInfoFromJwtToken(ctx)
	if !conf.IMConf.Debug {
		if msg := hasPerm(agentInfo.EntID, agentInfo.UserID, "visitor_and_conv", "redirect_others_conv"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &TransferConversationReqV1{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.GroupToken == "" {
		return invalidParameterResp(ctx, "group_token is invalid")
	}

	agent, err := models.GetAgentIDByConversationID(db.Mysql, convID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	targetAgent, err := allocateAgentFromGroup(agentInfo.EntID, req.GroupToken, agent)
	if err != nil {
		return errResp(ctx, common.AgentAllocateErr, err.Error())
	}

	if dbErr := models.UpdateConversationAgentID(db.Mysql, convID, targetAgent); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	conv, err := s.sendAgentTransferEvent(agent, targetAgent, agentInfo.UserID, agentInfo.EntID, convID)
	if err != nil {
		return errResp(ctx, common.DBErr, err.Error())
	}

	return jsonResponse(ctx, conv)
}

// POST /api/agent/conversation/:conversation_id/redirect
func (s *IMService) TransferToAgent(ctx echo.Context) (err error) {
	req := &TransferToAgentReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}
	if req.To == "" {
		return invalidParameterResp(ctx, "to agent id invalid")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	isOnline := models.IsAgentOnlineV1(req.To)
	if !isOnline {
		return errResp(ctx, common.AgentAllocateErr, "agent not online")
	}

	convID := ctx.Param("conversation_id")
	from, err := models.GetAgentIDByConversationID(db.Mysql, convID)
	if err != nil {
		return invalidParameterResp(ctx, "get from agent error")
	}

	userID := ctx.Get(middleware.AgentIDKey).(string)
	if dbErr := models.UpdateConversationAgentID(db.Mysql, convID, req.To); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	conv, err := s.sendAgentTransferEvent(from, req.To, userID, entID, convID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, conv)
}

func (s *IMService) sendAgentTransferEvent(from, to, by, entID, convID string) (conversation *adapter.Conversation, err error) {
	var dbErr error
	defer func() {
		if e := recover(); e != nil {
			log.Logger.Warnf("[sendAgentTransferEvent] panic: %v", e)
		}

		if dbErr != nil {
			log.Logger.Warnf("[sendAgentTransferEvent] error: %v", dbErr)
		}
	}()

	conv, dbErr := models.ConversationByID(db.Mysql, convID)
	if dbErr != nil {
		return nil, dbErr
	}

	key := fmt.Sprintf(common.AgentConversationNum, conv.AgentID)
	if err := db.RedisClient.Decr(key).Err(); err != nil {
		log.Logger.Warnf("decr conv num error: %v", err)
	}

	visits, dbErr := models.VisitsByEntIDTraceID(db.Mysql, entID, conv.TraceID)
	if dbErr != nil {
		return nil, dbErr
	}

	visitor, dbErr := models.VisitorByEntIDTraceID(db.Mysql, entID, conv.TraceID)
	if dbErr != nil {
		return nil, dbErr
	}
	visit := visits[0]

	fromAgent, dbErr := models.AgentByID(db.Mysql, from)
	if dbErr != nil {
		return nil, dbErr
	}

	var redirectBy *models.Agent
	if by == from {
		redirectBy = fromAgent
	} else {
		redirectBy, dbErr = models.AgentByID(db.Mysql, by)
		if dbErr != nil {
			return nil, dbErr
		}
	}

	toAgent, dbErr := models.AgentByID(db.Mysql, to)
	if dbErr != nil {
		return nil, dbErr
	}

	tags, dbErr := models.VisitorTagRelationsByVisitors(db.Mysql, []*models.Visitor{visitor})
	if dbErr != nil {
		return nil, dbErr
	}

	conversation = adapter.ModelConversationToConversation(conv, visit.ID, visit, visitor, tags)
	body := &AgentRedirectEventBody{
		ConversationBody: conversation,
		ConversationID:   convID,
		From:             adapter.ConvertAgentToAgentInfo(fromAgent, true),
		To:               adapter.ConvertAgentToAgentInfo(toAgent, true),
		RedirectedBy:     adapter.ConvertAgentToAgentInfo(redirectBy, true),
	}

	go s.sendAgentRedirectEvent(body)
	return conversation, nil
}

// AddSummary ...
// PUT /admin/api/v1/conversations/:conversation_id/summary
// PUT /api/conversation/:conversation_id/summary
func (s *IMService) AddSummary(ctx echo.Context) (err error) {
	convID := ctx.Param("conversation_id")
	req := &UpdateSummaryReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if err = models.UpdateConversationSummary(db.Mysql, convID, req.Content); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// POST /api/conversation/new
func (s *IMService) NewConversation(ctx echo.Context) error {
	req := &NewConversationReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.TrackID == "" {
		return invalidParameterResp(ctx, "track_id is empty")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	convNum, err := agentInfo.AgentActiveConvNum(agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	serveLimit, err := models.AgentServeLimitByID(db.Mysql, agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if convNum >= serveLimit {
		return invalidParameterResp(ctx, "已达服务上限")
	}

	visit, visitor, err := s.getVisitInfoByTrackID(entID, req.TrackID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	title := fmt.Sprintf("#%d-new-conv", time.Now().Unix())
	_, _, _, errMsg := s.newConversation(entID, req.TrackID, agentID, title, visit, visitor)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	queue := &models.VisitorQueue{TrackID: req.TrackID}
	if err := queue.Delete(db.Mysql); err == nil {
		go s.sendVisitorQueueingRemove(entID, req.TrackID)
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GET /api/conversation/search/v2
func (s *IMService) SearchConversations(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	var canCheckOthers bool
	errMsg := hasPerm(entID, agentID, "history_conv", "see_others_history_conv")
	if errMsg == nil {
		canCheckOthers = true
	}

	qs := ctx.QueryParams()
	start := qs["created_on"]
	var createStart, createEnd interface{}
	if len(start) > 0 {
		createTimeRange := &timeRange{}
		if err := common.Unmarshal(start[0], &createTimeRange); err != nil {
			return invalidParameterResp(ctx, err.Error())
		}

		createStart, createEnd, err = parseTime(createTimeRange)
		if err != nil {
			return invalidParameterResp(ctx, err.Error())
		}
	}

	end := qs["ended_on"]
	if len(end) == 0 {
		return invalidParameterResp(ctx, "ended_on invalid")
	}
	endTimeRange := &timeRange{}
	if err := common.Unmarshal(end[0], &endTimeRange); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	endStart, endEnd, err := parseTime(endTimeRange)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("enterprise_id", entID))
	if len(start) > 0 {
		queries = append(queries, elastic.NewRangeQuery("created_on").From(createStart).To(createEnd).Relation("within"))
	}

	queries = append(queries, elastic.NewRangeQuery("ended_on").From(endStart).To(endEnd).Relation("within"))

	if !canCheckOthers {
		queries = append(queries, elastic.NewTermQuery("agent_id", agentID))
	}

	if len(qs["visitor_name"]) > 0 {
		if v := qs["visitor_name"][0]; v != "" {
			queries = append(queries, elastic.NewMatchQuery("visitor.name", v))
		}
	}

	if len(qs["eva_level"]) > 0 {
		levels := qs["eva_level"]
		var boolQueries = elastic.NewBoolQuery()
		var evaQueries []elastic.Query

		for _, level := range levels {
			levelInt, err := strconv.Atoi(level)
			if err == nil {
				evaQueries = append(evaQueries, elastic.NewTermQuery("eva_level", levelInt))
			}
		}

		boolQueries.Should(evaQueries...)
		queries = append(queries, boolQueries)
	}

	// converse_duration {"condition":"gt","second":10}
	if len(qs["converse_duration"]) > 0 {
		type convDuration struct {
			Condition string `json:"condition"`
			Second    int    `json:"second"`
		}
		duration := qs["converse_duration"][0]

		d := &convDuration{}
		err := common.Unmarshal(duration, &d)
		if err == nil {
			if d.Condition == "lt" {
				queries = append(queries, elastic.NewRangeQuery("converse_duration").Lt(d.Second))
			} else {
				queries = append(queries, elastic.NewRangeQuery("converse_duration").Gt(d.Second))
			}
		}
	}

	if len(qs["agents"]) > 0 {
		agents := qs["agents"]
		var agentsI []interface{}
		for _, agent := range agents {
			agentsI = append(agentsI, agent)
		}

		q := elastic.NewBoolQuery().Filter(elastic.NewTermsQuery("agent_id", agentsI...))
		queries = append(queries, q)
	}

	if len(qs["tags"]) > 0 {
		tags := qs["tags"]
		var tagI []interface{}
		for _, tag := range tags {
			tagI = append(tagI, tag)
		}
		q := elastic.NewBoolQuery().Filter(elastic.NewTermsQuery("tags", tagI...))
		queries = append(queries, q)
	}

	if len(qs["region"]) > 0 {
		regionStr := qs["region"][0]
		if regionStr != "" {
			rg := &region{}
			if err := common.Unmarshal(regionStr, &rg); err == nil {
				if rg.Province != "" {
					if common.IsMunicipality(rg.Province) {
						queries = append(queries, elastic.NewMatchQuery("visit_info.city", rg.Province))
					} else {
						queries = append(queries, elastic.NewMatchQuery("visit_info.province", rg.Province))
					}
				}

				if rg.City != "" {
					queries = append(queries, elastic.NewMatchQuery("visit_info.city", rg.City))
				}
			}
		}
	}

	if len(qs["ip"]) > 0 {
		if v := qs["ip"][0]; v != "" {
			queries = append(queries, elastic.NewMatchQuery("visit_info.ip", v))
		}
	}

	// quality_grade: first_level, second_level, third_level, not_exists
	if len(qs["quality_grade"]) > 0 {
		grades := qs["quality_grade"]
		var boolQueries = elastic.NewBoolQuery()
		var gradeQueries []elastic.Query

		for _, grade := range grades {
			quality := grade
			if grade == "not_exists" {
				quality = ""
			}

			gradeQueries = append(gradeQueries, elastic.NewTermQuery("quality_grade", quality))
		}

		boolQueries.Should(gradeQueries...)
		queries = append(queries, boolQueries)
	}

	if len(qs["tel"]) > 0 {
		tel := qs["tel"][0]
		if tel != "" {
			queries = append(queries, elastic.NewMatchQuery("visitor.telephone", tel))
		}
	}
	if len(qs["comment"]) > 0 {
		if v := qs["comment"][0]; v != "" {
			queries = append(queries, elastic.NewMatchQuery("visitor.remark", v))
		}
	}

	if len(qs["content"]) > 0 {
		content := qs["content"][0]
		if content != "" {
			bq := elastic.NewBoolQuery()
			bq = bq.Must(elastic.NewTermQuery("messages.content", content))
			queries = append(queries, elastic.NewNestedQuery("messages", bq))
		}
	}

	var offset, limit = defaultOffset, defaultLimit
	if len(qs["page_count"]) > 0 {
		pageCount := qs["page_count"][0]
		if v, err := strconv.Atoi(pageCount); err == nil {
			limit = v
		}
	}

	if len(qs["page"]) > 0 {
		if v, err := strconv.Atoi(qs["page"][0]); err == nil {
			offset = v * limit
		}
	}

	query := elastic.NewBoolQuery()
	query = query.Must(queries...)
	res, total, err := elasticsearch.SearchConversations(elasticsearch.ESClient, query, offset, limit)
	if err != nil {
		log.Logger.Warnf("[SearchConversations] elasticsearch.SearchConversations error: %v\n", err)
	}

	if len(res) == 0 {
		return jsonResponse(ctx, &SearchConversationsResp{TotalCount: total, Conversations: []*adapter.SearchConversation{}})
	}

	var convIDs []string
	for _, conv := range res {
		convIDs = append(convIDs, conv.ID)
	}

	messages, err := models.MessagesByConversationIDs(db.Mysql, convIDs)
	if err != nil {
		log.Logger.Warnf("[SearchConversations] MessagesByConversationIDs error: %v\n", err)
	}

	agents, err := models.AgentsByMsgs(db.Mysql, messages)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentGroups, err := models.PermsRangeGroupIDsByAgents(db.Mysql, agents)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	adapterAgents := adapter.ConvertAgentsToAdapterAgentsV1(agents, agentGroups)

	var adapterMessages []*adapter.Message
	if messages != nil && len(messages) > 0 {
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].CreatedAt.Before(messages[j].CreatedAt)
		})

		adapterMessages = convertModelMessagesToAdapterMessages(messages, adapterAgents)
	}

	var convs = make([]*adapter.SearchConversation, 0, len(res))
	for _, esConv := range res {
		conv := esConv.Conversation
		conv.CreatedOn = common.ConvertUTCToLocal(conv.CreatedOn)
		esConv.VisitInfo.CreatedOn = common.ConvertUTCToLocal(esConv.VisitInfo.CreatedOn)
		esConv.VisitInfo.VisitedOn = common.ConvertUTCToLocal(esConv.VisitInfo.VisitedOn)

		if conv.LastUpdated != nil {
			t := common.ConvertUTCToLocal(*conv.LastUpdated)
			conv.LastUpdated = &t
		}

		if conv.LastMsgCreatedOn != nil {
			t := common.ConvertUTCToLocal(*conv.LastMsgCreatedOn)
			conv.LastMsgCreatedOn = &t
		}

		if conv.EndedOn != nil {
			t := common.ConvertUTCToLocal(*conv.EndedOn)
			conv.EndedOn = &t
		}

		if conv.LastUpdated != nil {
			t := common.ConvertUTCToLocal(*conv.LastUpdated)
			conv.LastUpdated = &t
		}

		for _, msg := range adapterMessages {
			if msg.ConversationID == conv.ID {
				conv.Messages = append(conv.Messages, msg)
			}
		}

		if conv.EvaLevel != nil {
			if *conv.EvaLevel < 0 {
				conv.EvaLevel = nil
			}
		}

		createdOn := *common.ConvertUTCToTimeString(conv.CreatedOn)

		var endedOn *string
		if conv.EndedOn != nil {
			endedOn = common.ConvertUTCToTimeString(*conv.EndedOn)
		}
		convs = append(convs, &adapter.SearchConversation{Conversation: conv, CreatedOn: createdOn, EndedOn: endedOn})
	}

	return jsonResponse(ctx, &SearchConversationsResp{TotalCount: total, Conversations: convs})
}

func parseTime(t *timeRange) (begin, end time.Time, err error) {
	layout := "2006-01-02T15:04:05.000000"
	begin, err = time.Parse(layout, t.Begin)
	if err != nil {
		return
	}

	end, err = time.Parse(layout, t.End)
	return
}

func (s *IMService) sendEndMessage(isSys bool, entID, traceID, convID, endBy string) {
	agents, err := models.AgentsByConversationID(db.Mysql, convID)
	if err != nil {
		agents = []string{endBy}
	}

	sendMsgReq := &sendMsgReq{
		EntID:       entID,
		TraceID:     traceID,
		Creator:     endBy,
		Agents:      agents,
		ConvID:      convID,
		FromType:    models.MessageFromAgentType,
		Content:     "",
		ContentType: models.MessageTextContentType,
		Internal:    false,
		System:      true,
	}

	configs, errMsg := s.getEnterpriseConfigs(entID)
	if errMsg != nil {
		log.Logger.Warnf("getEnterpriseConfigs error: %+v", errMsg)
	}

	var content string
	if configs != nil {
		settings := configs.EndingMsgSettings
		if settings != nil && settings.Web != nil {
			web := settings.Web
			if web.Status == "open" {
				if isSys {
					content = web.AutoEndingMessage
				} else {
					content = web.AgentEndingMessage
				}
			}
		}
	}

	if content != "" {
		sendMsgReq.Content = content
		if _, err := s.sendMsg(sendMsgReq); err != nil {
			log.Logger.Warnf("sendMsg error: %v", err)
		}
	}
}

func (s *IMService) getAgentConvNum(agentID string) int {
	key := fmt.Sprintf(common.AgentConversationNum, agentID)
	v, err := db.RedisClient.Get(key).Result()
	if err != nil {
		log.Logger.Warnf("get agent conversation number from redis err: %v", err)
		count, err := models.ConvNumByAgentID(db.Mysql, agentID)
		if err != nil {
			return 0
		}

		db.RedisClient.Set(key, count, 0)
		return count
	}

	convNum, err := strconv.Atoi(v)
	if err != nil {
		log.Logger.Warnf("get agent conversation number err: %v", err)
	}
	return convNum
}

func convertModelConversationToResp(conv *models.Conversation) *conversationResp {
	resp := &conversationResp{Conversation: conv}
	if conv.AgentType.Valid {
		resp.AgentType = conv.AgentType.String
	}

	if conv.ClientFirstSendTime.Valid {
		resp.ClientFirstSendTime = conv.ClientFirstSendTime.Time.Format(time.RFC3339)
	}

	if conv.FirstMsgCreatedAt.Valid {
		resp.FirstMsgCreatedAt = conv.FirstMsgCreatedAt.Time.Format(time.RFC3339)
	}

	if conv.FirstResponseWaitTime.Valid {
		resp.FirstResponseWaitTime = conv.FirstResponseWaitTime.Int64
	}

	if conv.LastMsgContent.Valid {
		resp.LastMsgContent = conv.LastMsgContent.String
	}

	if conv.LastMsgCreatedAt.Valid {
		resp.LastMsgCreatedAt = conv.LastMsgCreatedAt.Time.Format(time.RFC3339)
	}

	if conv.QualityGrade.Valid {
		resp.QualityGrade = conv.QualityGrade.String
	}

	if conv.EndedAt.Valid {
		resp.EndedAt = conv.EndedAt.Time.Format(time.RFC3339)
	}

	if conv.EndedBy.Valid {
		resp.EndedBy = conv.EndedBy.String
	}
	return resp
}

func CreateConversationDoc(modelConversation *models.Conversation) {
	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, modelConversation.EntID, modelConversation.TraceID)
	if err != nil {
		log.Logger.Warnf("VisitorByEntIDTraceID error: %v\n", err)
	}

	visits, err := models.VisitsByEntIDTraceID(db.Mysql, modelConversation.EntID, modelConversation.TraceID)
	if err != nil {
		log.Logger.Warnf("VisitsByEntIDTraceID, ent: %s, trace: %s, error: %v", modelConversation.EntID, modelConversation.TraceID, err)
	}

	var visit *models.Visit
	if len(visits) > 0 {
		visit = visits[0]
	}

	var visitID string
	if visit != nil {
		visitID = visit.ID
	}

	tags, err := models.VisitorTagRelationsByVisitors(db.Mysql, []*models.Visitor{visitor})
	if err != nil {
		log.Logger.Warnf("VisitorTagRelationsByVisitors, error: %v", err)
	}

	doc := &elasticsearch.ConversationV1{
		Conversation: adapter.ModelConversationToConversation(modelConversation, visitID, visit, visitor, tags),
	}

	if visitor != nil {
		doc.Visitor = &elasticsearch.SimpleVisitor{
			Name:      visitor.Name,
			Telephone: visitor.Mobile,
			Remark:    visitor.Remark,
		}
	}

	msgs, err := models.MessagesByConversationID(db.Mysql, modelConversation.ID, 0, 200)
	if err != nil {
		log.Logger.Warnf("MessagesByConversationID error: %v\n", err)
	}
	for _, msg := range msgs {
		doc.Messages = append(doc.Messages, &elasticsearch.Message{ID: msg.ID, Content: msg.Content.String})
	}

	retryFunc := func() error {
		return elasticsearch.CreateConversationDoc(elasticsearch.ESClient, doc)
	}

	err = retry.Do(
		retryFunc,
		retry.Attempts(2),
		retry.Delay(10*time.Millisecond),
	)
	if err != nil {
		log.Logger.Warnf("elastic search CreateConversationDoc error: %v\n", err)
	}
}

func updateConversationStats(conv *models.Conversation) {
	var values = []string{
		"total_count",
		"effective_count",
		"message_count",
		"duration_in_sec",
		"wait_time_in_sec",
	}

	convStats := &models.ConversationStat{
		EntID:          conv.EntID,
		TotalCount:     1,
		EffectiveCount: 1,
		MessageCount:   conv.MsgCount,
		DurationInSec:  conv.Duration,
		WaitTimeInSec:  uint(conv.FirstResponseWaitTime.Int64),
		CreatedAt:      getDate(),
	}
	switch conv.EvalLevel {
	case models.GoodConversation:
		values = append(values, "good_count")
		convStats.GoodCount = 1
	case models.MediumConversation:
		values = append(values, "medium_count")
		convStats.MediumCount = 1
	case models.BadConversation:
		values = append(values, "bad_count")
		convStats.BadCount = 1
	}
	switch conv.QualityGrade.String {
	case models.FirstLevel:
		values = append(values, "gold_count")
		convStats.GoldCount = 1
	case models.SecondLevel:
		values = append(values, "silver_count")
		convStats.SilverCount = 1
	case models.ThirdLevel:
		values = append(values, "bronze_count")
		convStats.BronzeCount = 1
	default:
		values = append(values, "no_grade_count")
		convStats.NoGradeCount = 1
	}

	if err := convStats.Insert(db.Mysql, values); err != nil {
		log.Logger.WithField("func_name", "updateConversationStats").
			WithField("conversation_id", conv.ID).
			WithField("ent_id", conv.EntID).
			Warnf("error: %v\n", err)
	}
}

func updateAgentStats(conv *models.Conversation) {
	var err error
	var values = []string{
		"conversation_count",
		"message_count",
		"duration",
		"first_resp_duration",
	}

	agentStats := &models.AgentStatistic{
		EntID:             conv.EntID,
		AgentID:           conv.AgentID,
		ConversationCount: 1,
		MessageCount:      conv.AgentMsgCount,
		Duration:          int(conv.Duration),
		FirstRespDuration: int(conv.FirstResponseWaitTime.Int64),
		CreatedAt:         getDate(),
	}

	switch conv.EvalLevel {
	case models.GoodConversation:
		values = append(values, "good_count")
		agentStats.GoodCount = 1
	case models.MediumConversation:
		values = append(values, "medium_count")
		agentStats.MediumCount = 1
	case models.BadConversation:
		values = append(values, "bad_count")
		agentStats.BadCount = 1
	}

	if err = agentStats.Insert(db.Mysql, values); err != nil {
		log.Logger.WithField("func_name", "updateAgentStats").
			WithField("conversation_id", conv.ID).
			WithField("ent_id", conv.EntID).
			Warnf("models.CreateOrUpdateAgentStats error: %v\n", err)
	}
}

func getDate() time.Time {
	now := time.Now().UTC()
	year, month, day := now.Date()
	hour := now.Hour()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}
