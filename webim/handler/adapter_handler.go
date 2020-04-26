package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type EnterprisePlan struct {
	EntID           string    `json:"ent_id"`
	CreatedAt       time.Time `json:"created_at"`
	ExpirationTime  time.Time `json:"expiration_time"`
	LoginAgentLimit int       `json:"login_agent_limit"`
	AgentNum        int       `json:"agent_num"`
	Plan            int       `json:"plan"`
	TrialStatus     int       `json:"trial_status"`
}

type AgentStats struct {
	AvgDurationTime  int `json:"avg_duration_time"`
	AvgWaitTime      int `json:"avg_wait_time"`
	ConvCnt          int `json:"conv_cnt"`
	DurationTime     int `json:"duration_time"`
	MsgCnt           int `json:"msg_cnt"`
	TicketResolveCnt int `json:"ticket_resolve_cnt"`
	WaitTime         int `json:"wait_time"`
}

type AgentStatsResp struct {
	// summary: map[agent_id]-> AgentStats
	Summary map[string]*AgentStats `json:"summary"`
}

type GetPromotionsResp struct {
	RichTexts []*adapter.PromotionMessage `json:"rich_texts"`
}

type SetPromotionMessagesReq struct {
	Content   string `json:"content"`
	Thumbnail string `json:"thumbnail"`
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Enabled   bool   `json:"enabled"`
	Countdown int    `json:"countdown"`
	BrowserID string `json:"browser_id"`
}

type SetPromotionMessagesResp = adapter.PromotionMessage

type UpdatePromotionMessagesReq struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Thumbnail string `json:"thumbnail"`
	Summary   string `json:"summary"`
	Countdown *int   `json:"countdown"`
	Enabled   *bool  `json:"enabled"`
}

type SelectingRules struct {
	Rules []*dto.SelectingRule `json:"rules"`
}

// {"agent_welcome_msg_settings":{"message":"","status":"open"}}
type UpdateAgentPersonalConfigsReq struct {
	AgentWelcomeMsgSettings *adapter.WelcomeMsg      `json:"agent_welcome_msg_settings"`
	ConvGroupConfig         *adapter.ConvGroupConfig `json:"conv_group_config"`
	ConvOrderConfig         *adapter.ConvOrderConfig `json:"conv_order_config"`
	MessageTone             *adapter.MessageTone     `json:"message_tone"`
	QuickReplyRule          string                   `json:"quick_reply_rule"`
}

func (s *IMService) GetEnterprisePlan(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	plan, err := models.EntPlanByEntID(db.Mysql, entID)
	if err != nil {
		log.Logger.Errorf("get EntPlanByEntID error: %v", err)
		return ctx.JSON(http.StatusOK, &adapter.Plan{})
	}

	adapterPlan := adapter.Plan{
		ID:               plan.ID,
		CreatedOn:        plan.CreateAt,
		ExpirationTime:   plan.ExpirationTime,
		LoginAgentLimit:  plan.LoginAgentLimit,
		PayAgentQuantity: plan.AgentNum,
		PayAmount:        plan.PayAmount,
		Plan:             int(plan.PlanType),
		TrialStatus:      plan.TrialStatus,
		Valid:            true,
		VisitorLimit:     -1,
	}

	return ctx.JSON(http.StatusOK, adapterPlan)
}

// GET /api/enterprise/features
func (s *IMService) GetEnterpriseFeatures(ctx echo.Context) (err error) {
	return jsonResponse(ctx, &Resp{
		Code: 0,
		Body: struct {
			Features []string `json:"features"`
		}{
			Features: dto.Features,
		},
	})
}

// GET /api/enterprise/promotion_msgs
func (s *IMService) GetPromotionMessages(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_auto_message"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	promotions, err := models.PromotionMsgsByEnterpriseID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &GetPromotionsResp{
		RichTexts: []*adapter.PromotionMessage{},
	}
	for _, p := range promotions {
		resp.RichTexts = append(resp.RichTexts, modelPromotionMsgToView(p))
	}

	return jsonResponse(ctx, resp)
}

// POST /api/enterprise/promotion_msgs
func (s *IMService) SetPromotionMessages(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_auto_message"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &SetPromotionMessagesReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	promotion := &models.PromotionMsg{
		ID:           common.GenUniqueID(),
		EnterpriseID: entID,
		Source:       req.Source,
		Content:      sql.NullString{String: req.Content, Valid: true},
		ContentSdk:   "",
		Countdown:    req.Countdown,
		Enabled:      req.Enabled,
		Summary:      req.Summary,
		Thumbnail:    req.Thumbnail,
		CreatedOn:    time.Now().UTC(),
		UpdatedOn:    time.Now().UTC(),
	}

	if err = promotion.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, modelPromotionMsgToView(promotion))
}

// PUT /api/enterprise/promotion_msgs/:msg_id
func (s *IMService) UpdatePromotionMessages(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_auto_message"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	msgID := ctx.Param("msg_id")
	if msgID == "" {
		return invalidParameterResp(ctx, "msg_id invalid")
	}

	req := &UpdatePromotionMessagesReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	msg, err := models.PromotionMsgByID(db.Mysql, msgID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "promotion_msg not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	if req.Content != "" {
		msg.Content.String = req.Content
		msg.Content.Valid = true
	}

	if req.Thumbnail != "" {
		msg.Thumbnail = req.Thumbnail
	}

	if req.Summary != "" {
		msg.Summary = req.Summary
	}

	var countDown = req.Countdown
	if countDown != nil {
		msg.Countdown = *countDown
	}

	if req.Enabled != nil {
		msg.Enabled = *req.Enabled
	}

	msg.CreatedOn = msg.CreatedOn.UTC()
	msg.UpdatedOn = time.Now().UTC()
	if err = msg.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, modelPromotionMsgToView(msg))
}

// DELETE /api/enterprise/promotion_msgs/:msg_id
func (s *IMService) DeletePromotionMessage(ctx echo.Context) (err error) {
	msgID := ctx.Param("msg_id")
	if msgID == "" {
		return invalidParameterResp(ctx, "msg_id invalid")
	}

	msg := &models.PromotionMsg{ID: msgID}
	if err := msg.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GET /api/agent/sales_config
func (s *IMService) GetEnterpriseSalesConfig(ctx echo.Context) (err error) {
	return jsonResponse(ctx, struct {
		Code    int      `json:"code"`
		Body    []string `json:"body"`
		Message string   `json:"message"`
	}{
		Code:    0,
		Message: "",
		Body:    []string{},
	})
}

func (s *IMService) GetMentionedConvs(ctx echo.Context) (err error) {
	return jsonResponse(ctx, struct {
		Convs []interface{} `json:"convs"`
	}{
		Convs: []interface{}{},
	})
}

func (s *IMService) GPUL(ctx echo.Context) (err error) {
	return jsonResponse(ctx, struct {
		ErrorCode int    `json:"errorcode"`
		ErrorMsg  string `json:"errormsg"`
	}{
		ErrorCode: 5,
		ErrorMsg:  "all is well",
	})
}

// GET /api/agent/personal_config
func (s *IMService) GetAgentPersonalConfigs(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	personalConfig, err := models.PersonalConfigByAgentID(db.Mysql, agentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return jsonResponse(ctx, adapter.DefaultPersonalConfig)
		}

		return dbErrResp(ctx, err.Error())
	}

	var config = &adapter.PersonalConfig{}
	if personalConfig.ConfigContent.String != "" {
		if err = common.Unmarshal(personalConfig.ConfigContent.String, &config); err != nil {
			return errResp(ctx, common.DecodeJSONErr, err.Error())
		}
	}

	return jsonResponse(ctx, config)
}

// POST /api/agent/personal_config
func (s *IMService) UpdateAgentPersonalConfigs(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	req := &UpdateAgentPersonalConfigsReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	personalConfig, err := models.PersonalConfigByAgentID(db.Mysql, agentID)
	if err != nil && err != sql.ErrNoRows {
		return dbErrResp(ctx, err.Error())
	}

	if personalConfig == nil {
		personalConfig = &models.PersonalConfig{}
	}

	configContent := personalConfig.ConfigContent.String
	if configContent == "" {
		configContent = `{}`
	}

	var config = &adapter.PersonalConfig{}
	if err = common.Unmarshal(configContent, &config); err != nil {
		log.Logger.Warnf("[UpdateAgentPersonalConfigs] decode json: %s error: %+v", configContent, err)
		return errResp(ctx, common.DecodeJSONErr, err.Error())
	}

	if req.MessageTone != nil {
		config.MessageTone = req.MessageTone
	}
	if req.AgentWelcomeMsgSettings != nil {
		config.WelcomeMsg = req.AgentWelcomeMsgSettings
	}
	if req.QuickReplyRule != "" {
		config.QuickReplyRule = req.QuickReplyRule
	}
	if req.ConvOrderConfig != nil {
		config.ConvOrderConfig = req.ConvOrderConfig
	}
	if req.ConvGroupConfig != nil {
		config.ConvGroupConfig = req.ConvGroupConfig
	}

	newConfigContent, err := common.Marshal(config)
	if err != nil {
		return errResp(ctx, common.EncodeJSONErr, err.Error())
	}
	personalConfig.ConfigContent = sql.NullString{String: newConfigContent, Valid: true}
	if err = personalConfig.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, config)
}

// GET /api/enterprise
func (s *IMService) GetEnterpriseInfo(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的租户")
		}

		return dbErrResp(ctx, err.Error())
	}

	return ctx.JSON(http.StatusOK, adapter.ConvertEntToAdapterEnt(ent))
}

// GetEnterprisePerms ...
// GET /api/perm/perms
func (s *IMService) GetEnterprisePerms(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	perms, err := models.PermsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var permsMap = make(map[string][]*models.Perm, len(perms))
	for _, p := range perms {
		if v, ok := permsMap[p.AppName]; ok {
			v = append(v, p)
			permsMap[p.AppName] = v
		} else {
			permsMap[p.AppName] = []*models.Perm{p}
		}
	}

	return jsonResponse(ctx, adapter.ConvertEntPermsToPerms(permsMap))
}

// GetAgentsStats ...
// GET /api/stats/agents?begin=2019-01-24&end=2019-01-25
func (s *IMService) GetAgentsStats(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	startTime := ctx.QueryParam("begin")
	endTime := ctx.QueryParam("end")
	if startTime == "" || endTime == "" {
		return invalidParameterResp(ctx, "begin/end invalid")
	}

	agents, err := models.AgentsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	start, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	end, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	agentsStats, err := models.AgentStatisticsByEntID(db.Mysql, entID, start, end)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &AgentStatsResp{
		Summary: make(map[string]*AgentStats, len(agents)),
	}

	var defaultStats = &AgentStats{}

	for _, agent := range agents {
		var stat *models.AgentStatistic
		for _, stats := range agentsStats {
			if stats.AgentID == agent.ID {
				stat = stats
				break
			}
		}

		resp.Summary[agent.ID] = defaultStats
		if stat != nil {
			agentStat := &AgentStats{
				DurationTime:     0,
				ConvCnt:          int(stat.ConversationCount),
				MsgCnt:           int(stat.MessageCount),
				TicketResolveCnt: 0,
				WaitTime:         0,
			}
			agentStat.AvgDurationTime = int(stat.Duration / int(stat.ConversationCount))
			agentStat.AvgWaitTime = int(stat.FirstRespDuration / int(stat.ConversationCount))
			resp.Summary[agent.ID] = agentStat
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

//GET /api/client/custom_field
func (s *IMService) GetClientCustomFields(ctx echo.Context) (err error) {
	return
}

// GET /api/client/attr_order
func (s *IMService) GetClientAttrOrder(ctx echo.Context) (err error) {
	return nil
}

// GET /api/agent/selecting_rules
func (s *IMService) GetSelectingRules(ctx echo.Context) error {
	return jsonResponse(ctx, &SelectingRules{Rules: []*dto.SelectingRule{}})
}
