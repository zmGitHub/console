package handler

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/handler/allocator"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AllocateAgentReq struct {
	EntID     string `query:"ent_id"`
	MenuValue string `query:"menu_value"`
}

type AllocateAgentResp struct {
	Code    int    `json:"code"`
	AgentID string `json:"agent_id"`
	Rule    string `json:"rule"`
}

type newConv struct {
	Action        string                `json:"action"`
	AgentID       *string               `json:"agent_id"`
	AgentNickname *string               `json:"agent_nickname"`
	Agent         *adapter.Agent        `json:"agent"`
	Body          *adapter.Conversation `json:"body"`
	CreatedOn     time.Time             `json:"created_on"`
	EnterpriseID  string                `json:"enterprise_id"`
	ID            string                `json:"id"`
	Source        string                `json:"source"`
	TargetID      string                `json:"target_id"`
	TargetKind    string                `json:"target_kind"`
	TraceStart    float64               `json:"trace_start"`
	TrackID       string                `json:"track_id"`
}

// AllocateAgent allocate agent by the setting rules
// GET /api/v1/allocate_agent?trace_id=xxxx&ent_id=yyyyy&menu_value={}
func (s *IMService) AllocateAgent(ctx echo.Context) (err error) {
	req := &AllocateAgentReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.EntID == "" {
		return invalidParameterResp(ctx, "ent_id is invalid")
	}

	if req.MenuValue != "" {
		agentID, err := s.allocateAgentFromMenuValue(req.MenuValue, req.EntID)
		if err == nil {
			return jsonResponse(ctx, &AllocateAgentResp{Code: 0, AgentID: agentID, Rule: ""})
		}

		log.Logger.Errorf("allocateAgentFromMenuValue error: %v", err)
	}

	allocateRuleType, err := models.AllocationRuleByEntID(db.Mysql, req.EntID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	var alc allocator.Allocator

	switch allocateRuleType {
	case models.ConversationNumAllocation:
		alc = &allocator.LoadBalanced{EntID: req.EntID, AgentInfo: agentInfo}
	case models.OrderTakeTurnsAllocation:
		alc = &allocator.TakeTurns{EntID: req.EntID, AgentInfo: agentInfo}
	case models.OrderPriorityAllocation:
		alc = &allocator.Priority{EntID: req.EntID, AgentInfo: agentInfo}
	default:
		return invalidParameterResp(ctx, "unsupported rule_type")
	}

	agentID, err := alc.Allocate()
	if err != nil {
		return jsonResponse(ctx, &ErrMsg{Code: common.AgentAllocateErr, Message: err.Error()})
	}

	if err = agentInfo.SetLastAllocatedAgent(req.EntID, agentID); err != nil {
		log.Logger.Error("SetLastAllocatedAgent error: ", err)
	}

	return jsonResponse(ctx, &AllocateAgentResp{Code: 0, AgentID: agentID, Rule: allocateRuleType})
}

func (s *IMService) allocateAgentFromMenuValue(menuValue, entID string) (agentID string, err error) {
	menu := &menuField{}
	err = common.Unmarshal(menuValue, &menu)
	if err != nil {
		return
	}

	if menu.AgentType == "group" {
		return allocateAgentFromGroup(entID, menu.Value, "")
	}

	isAgentOnline, err := s.imCli.IsOnline(context.Background(), menu.Value)
	if err != nil {
		return
	}

	if !isAgentOnline {
		return "", fmt.Errorf("agent not online")
	}

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	convNum, err := agentInfo.AgentActiveConvNum(menu.Value)
	if err != nil {
		return
	}

	serveLimit, err := models.AgentServeLimitByID(db.Mysql, menu.Value)
	if err != nil {
		return
	}

	if convNum < serveLimit {
		return menu.Value, nil
	}

	return "", fmt.Errorf("agent conversation num exceed")
}

func (s *IMService) AllocateAgentByEntID(entID string) (agentID string, errMsg *ErrMsg) {
	allocateRuleType, err := models.AllocationRuleByEntID(db.Mysql, entID)
	if err != nil {
		return "", &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	var alc allocator.Allocator

	switch allocateRuleType {
	case models.ConversationNumAllocation:
		alc = &allocator.LoadBalanced{EntID: entID, AgentInfo: agentInfo}
	case models.OrderTakeTurnsAllocation:
		alc = &allocator.TakeTurns{EntID: entID, AgentInfo: agentInfo}
	case models.OrderPriorityAllocation:
		alc = &allocator.Priority{EntID: entID, AgentInfo: agentInfo}
	default:
		return "", &ErrMsg{Code: common.InvalidParameterErr, Message: "unsupported rule_type"}
	}

	agentID, err = alc.Allocate()
	if err != nil {
		return agentID, &ErrMsg{Code: common.AgentAllocateErr, Message: err.Error()}
	}

	if err = agentInfo.SetLastAllocatedAgent(entID, agentID); err != nil {
		log.Logger.Warnf("SetLastAllocatedAgent error: %v", err)
	}

	return agentID, nil
}

func allocateAgentFromGroup(entID, groupID string, excludeAgentID string) (agentID string, err error) {
	agents, err := models.AgentsByGroupID(db.Mysql, groupID)
	if err != nil {
		return "", err
	}

	var rks []*models.AgentRanking
	var ids []string
	for _, agent := range agents {
		if agent.ID == excludeAgentID {
			continue
		}

		rks = append(rks, &models.AgentRanking{
			AgentID:    agent.ID,
			Ranking:    agent.Ranking,
			ServeLimit: agent.ServeLimit,
		})
		ids = append(ids, agent.ID)
	}

	if len(rks) == 0 {
		return "", common.NoOnlineAgents
	}

	onlineAgents := models.FilterOnline(ids)
	if len(onlineAgents) == 0 {
		return "", common.NoOnlineAgents
	}

	var onlineRKs []*models.AgentRanking
	for _, rk := range rks {
		for _, id := range onlineAgents {
			if id == rk.AgentID {
				onlineRKs = append(onlineRKs, rk)
				break
			}
		}
	}

	sort.SliceStable(onlineRKs, func(i, j int) bool {
		return onlineRKs[i].Ranking <= onlineRKs[j].Ranking
	})

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	return allocator.TakeTurnsByAgentIDs(agentInfo, entID, onlineRKs)
}

// CreateOrUpdateAllocationRule
// POST /admin/api/v1/allocation_rules
// PUT /api/enterprise/allocation_rule
func (s *IMService) CreateOrUpdateAllocationRule(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_agent_allocation"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	rule := &struct {
		RuleType string `json:"rule"`
	}{}
	if err = ctx.Bind(rule); err != nil {
		return
	}

	if _, ok := models.AllocationRuleTypeMap[rule.RuleType]; !ok {
		return invalidParameterResp(ctx, "unsupported rule_type")
	}

	if err := models.UpdateEntInfo(db.Mysql, entID, map[string]interface{}{"allocation_rule": rule.RuleType}); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, adapter.ConvertEntToAdapterEnt(ent))
}

// Scheduler ...
// POST /scheduler?ent_id=xxxxx
func (s *IMService) Scheduler(ctx echo.Context) (err error) {
	req := &adapter.SchedulerRequest{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	visit, visitor, err := s.getVisitInfoByTrackID(req.EntID, req.TrackID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	entModel, err := models.EnterpriseByID(db.Mysql, req.EntID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	entInfo := adapter.ConvertEntToAdapterEnt(entModel)

	activeConversations, err := models.ConversationsByTraceID(db.Mysql, req.TrackID, 0, 1, true)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if len(activeConversations) > 0 {
		agentID := activeConversations[0].AgentID
		agentInfo, convInfo, errMsg := s.newAdapterConversation(agentID, activeConversations[0], visit, visitor)
		if errMsg != nil {
			return invalidParameterResp(ctx, errMsg.Message)
		}

		sc := &adapter.SchedulerResponse{
			Agent:         agentInfo,
			Conv:          convInfo,
			ConvNewCreate: false,
			Ent:           entInfo,
			Messages:      nil,
			Result:        "existing",
			Success:       true,
		}
		setLatestMessage(sc, req.TrackID, agentInfo)
		return jsonResponse(ctx, sc)
	}

	configs, errMsg := s.getEnterpriseConfigs(req.EntID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	queueOpen := configs.QueueSettings != nil && configs.QueueSettings.Status == "open"
	queueLimit := 0
	if configs.QueueSettings != nil {
		queueLimit = configs.QueueSettings.QueueSize
	}

	queueCount, err := models.VisitorQueueCount(db.Mysql, req.EntID)
	if err != nil && err != sql.ErrNoRows {
		return dbErrResp(ctx, err.Error())
	}

	shouldEnqueue := queueOpen && queueCount != 0 && queueCount < queueLimit && req.Queueing && req.AgentToken != nil
	if shouldEnqueue {
		var resp *adapter.SchedulerResponse
		errMsg := s.EnqueueVisitor(req.EntID, req.TrackID, req.VisitID, *req.AgentToken, visit, visitor, resp)
		if errMsg != nil {
			return internalServerErr(ctx, errMsg.Message)
		}

		return jsonResponse(ctx, resp)
	}

	var agentID string
	var agentAllocated bool
	// 询前表单制定的坐席
	if req.AgentToken != nil {
		if s.shouldScheduler(*req.AgentToken) {
			agentID = *req.AgentToken
			agentAllocated = true
		}
	}

	if req.GroupToken != nil {
		agentID, err = allocateAgentFromGroup(req.EntID, *req.GroupToken, "")
		if err != nil {
			agentAllocated = false
			log.Logger.Warnf("[allocateAgentFromGroup] error: %v", err)
		} else {
			agentAllocated = true
		}
	}

	if !agentAllocated {
		agentID, errMsg = s.AllocateAgentByEntID(req.EntID)
		if errMsg != nil {
			log.Logger.Warnf("[Scheduler] Allocate Agent Error: %+v", errMsg)

			if errMsg.Message == common.AgentServeLimitExceed.Error() && queueOpen {
				if queueCount >= queueLimit {
					log.Logger.Warnf("[Scheduler] queueCount(%d) >= queueLimit(%d)", queueCount, queueLimit)
					return jsonResponse(ctx, &adapter.SchedulerResponse{
						ConvNewCreate: false,
						ReserveToken:  "",
						Result:        "fail",
						Success:       false,
					})
				}

				q := &models.VisitorQueue{EntID: req.EntID, TrackID: req.TrackID, VisitID: req.VisitID, EnqueueAt: time.Now().UTC()}
				if err := q.Insert(db.Mysql); err != nil {
					return dbErrResp(ctx, err.Error())
				}

				pos, err := models.GetVisitorPosition(db.Mysql, req.TrackID)
				if err != nil {
					pos = 1
				}

				go s.sendQueueingAddEvent(agentID, visit, visitor)
				return jsonResponse(ctx, &adapter.SchedulerResponse{
					ConvNewCreate: false,
					Position:      &pos,
					ReserveToken:  "",
					Result:        "queueing",
					Success:       false,
				})
			}

			return jsonResponse(ctx, &adapter.SchedulerResponse{
				ConvNewCreate: false,
				ReserveToken:  "",
				Result:        "fail",
				Success:       false,
			})
		}
	}

	agentInfo, convInfo, isNewCreated, errMsg := s.newConversation(req.EntID, req.TrackID, agentID, req.Title, visit, visitor)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	sc := &adapter.SchedulerResponse{
		Agent:         agentInfo,
		Conv:          convInfo,
		ConvNewCreate: isNewCreated,
		Ent:           entInfo,
		Messages:      []*adapter.Message{},
		ReserveToken:  "",
		Result:        "new_conv",
		Success:       true,
	}

	setLatestMessage(sc, req.TrackID, agentInfo)
	if !isNewCreated {
		sc.Result = "existing"
	}

	return jsonResponse(ctx, sc)
}

func setLatestMessage(sc *adapter.SchedulerResponse, trackID string, agentInfo *adapter.Agent) {
	lastMsg, err := models.LastMessagebyTraceID(db.Mysql, trackID)
	if err == nil {
		var agt *adapter.Agent
		if lastMsg.AgentID == agentInfo.ID {
			agt = agentInfo
		} else {
			agent, err := models.AgentByID(db.Mysql, lastMsg.AgentID)
			if err == nil {
				agt = adapter.ConvertAgentToAgentInfo(agent, false)
			}
		}
		if agt != nil {
			sc.Messages = append(sc.Messages, modelMsgToAdapterMessage(agt, lastMsg))
		}
	}
}

func (s *IMService) getVisitInfoByTrackID(entID, trackID string) (visit *models.Visit, visitor *models.Visitor, err error) {
	visits, err := models.VisitsByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		return
	}

	visitor, err = models.VisitorByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		return
	}

	return visits[0], visitor, nil
}

func (s *IMService) newConversation(entID, trackID, agentID, title string, visit *models.Visit, visitor *models.Visitor) (*adapter.Agent, *adapter.Conversation, bool, *ErrMsg) {
	conv, isNewCreated, errMsg := s.initConversation(entID, trackID, agentID, title)
	if errMsg != nil {
		return nil, nil, false, errMsg
	}

	agentInfo, convInfo, errMsg := s.newAdapterConversation(agentID, conv, visit, visitor)
	if errMsg != nil {
		return nil, nil, false, errMsg
	}

	if isNewCreated {
		if err := s.sendNewConversationToAgents(agentInfo, convInfo); err != nil {
			log.Logger.Warnf("[sendNewConversationToAgents] error: %v", err)
		}

		go s.AddTimedTasks(entID, convInfo.TrackID, agentInfo.ID, convInfo.ID, visitor.ID)
	}

	return agentInfo, convInfo, isNewCreated, nil
}

func (s *IMService) newAdapterConversation(agentID string, conv *models.Conversation, visit *models.Visit, visitor *models.Visitor) (agent *adapter.Agent, conv1 *adapter.Conversation, errMsg *ErrMsg) {
	modelAgent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return nil, nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	agent = adapter.ConvertAgentToAgentInfo(modelAgent, true)
	tags, err := models.VisitorTagRelationsByVisitors(db.Mysql, []*models.Visitor{visitor})
	if err != nil {
		return nil, nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	conv1 = adapter.ModelConversationToConversation(conv, visit.ID, visit, visitor, tags)
	return
}
func (s *IMService) sendNewConversationToAgents(agent *adapter.Agent, conv *adapter.Conversation) error {
	var agents []string
	var err error
	if agent.PrivilegeRange == models.AgentPermsRangeAllType {
		agents, err = models.AgentIDsByEntID(db.Mysql, agent.EnterpriseID)
	} else if agent.PrivilegeRange == models.AgentPermsRangePersonalType {
		agents = []string{agent.ID}
	} else {
		groups, ok := agent.PrivilegeRange.([]string)
		if !ok {
			agents = []string{agent.ID}
		} else {
			agents, err = models.AgentIDsByGroupIDs(db.Mysql, groups)
		}
	}

	if err != nil {
		return err
	}

	if len(agents) == 0 {
		agents = []string{agent.ID}
	}

	newConv := &newConv{
		Action:        "init_conv",
		AgentID:       &agent.ID,
		AgentNickname: &agent.Nickname,
		Body:          conv,
		CreatedOn:     conv.CreatedOn,
		EnterpriseID:  agent.EnterpriseID,
		ID:            conv.ID,
		Source:        conv.VisitInfo.FirstVisitPageSourceBySession,
		TargetID:      conv.ID,
		TargetKind:    "conv",
		TraceStart:    float64(conv.CreatedOn.UnixNano()),
		TrackID:       conv.TrackID,
	}

	content, err := common.Marshal(newConv)
	if err != nil {
		return err
	}

	var hasAgent bool
	for _, id := range agents {
		if id == agent.ID {
			hasAgent = true
			break
		}
	}

	if !hasAgent {
		agents = append(agents, agent.ID)
	}

	s.sendMessageToMultiAgents(agents, content)
	return nil
}

func (s *IMService) convertToNewConvAction(conv *adapter.Conversation) *newConv {
	return &newConv{
		Action:        "init_conv",
		AgentID:       nil,
		AgentNickname: nil,
		Body:          conv,
		CreatedOn:     conv.CreatedOn,
		EnterpriseID:  conv.EnterpriseID,
		ID:            conv.ID,
		Source:        conv.VisitInfo.FirstVisitPageSourceBySession,
		TargetID:      conv.ID,
		TargetKind:    "conv",
		TraceStart:    float64(conv.CreatedOn.UnixNano()),
		TrackID:       conv.TrackID,
	}
}

func (s *IMService) EnqueueVisitor(entID, trackID, visitID, agentID string, visit *models.Visit, visitor *models.Visitor, resp *adapter.SchedulerResponse) *ErrMsg {
	q := &models.VisitorQueue{EntID: entID, TrackID: trackID, VisitID: visitID, EnqueueAt: time.Now().UTC()}
	if err := q.Insert(db.Mysql); err != nil {
		return &ErrMsg{Message: err.Error()}
	}

	go s.sendQueueingAddEvent(agentID, visit, visitor)

	pos, err := models.GetVisitorPosition(db.Mysql, trackID)
	if err != nil {
		pos = 1
	}

	resp = &adapter.SchedulerResponse{
		ConvNewCreate: false,
		Position:      &pos,
		ReserveToken:  "",
		Result:        "queueing",
		Success:       false,
	}

	return nil
}

func (s *IMService) shouldScheduler(agentID string) bool {
	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	convNum, err := agentInfo.AgentActiveConvNum(agentID)
	if err != nil {
		log.Logger.Warnf("[shouldScheduler] error: %v", err)
		return false
	}

	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		log.Logger.Warnf("[shouldScheduler] error: %v", err)
		return false
	}

	return convNum < agent.ServeLimit && agent.Status == models.AgentAvailableStatus
}

func (s *IMService) sendQueueingAddEvent(agentID string, visit *models.Visit, visitor *models.Visitor) {
	visitInfo := dto.ModelVisitInfoToVisitInfo(visit, visitor)
	event := events.NewQueueingAdd([]string{agentID}, visitor.ID, visitInfo)

	bs, err := common.Marshal(event)
	if err == nil {
		s.sendMessageToMultiAgents([]string{agentID}, string(bs))
	}
}
