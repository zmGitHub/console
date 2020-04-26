package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type ConvConfigs struct {
	EndingConversationConf   *models.EndingConversation   `json:"ending_conversation_conf"`
	EndingMessageConf        *models.EndingMessage        `json:"ending_message_conf"`
	VisitorQueueConf         *models.QueueConfig          `json:"visitor_queue_conf"`
	ConversationTransferConf *models.ConversationTransfer `json:"conversation_transfer_conf"`
	ConversationQualityConf  *models.ConversationQuality  `json:"conversation_quality_conf"`
}

type GetEntConfigsResp struct {
	ConvConfigs        *ConvConfigs               `json:"conversation_configs"`
	SecurityConfig     *SecurityConfig            `json:"security_config,omitempty"`
	AutomaticMessages  []*models.AutomaticMessage `json:"automatic_messages,omitempty"`
	LeaveMessageConfig *models.LeaveMessageConfig `json:"leave_message_config,omitempty"`
}

type UpdateConvConfigsReq struct {
	ConvConfigs *ConvConfigs `json:"conv_configs"`
}

type SecurityConfig struct {
	LoginLimit *models.EntLoginLimit `json:"login_limit"`
	SendFile   *models.SendFile      `json:"send_file"`
}

type ChangeMsgStatusReq struct {
	MsgType string `json:"msg_type"` // msg_type
	Source  string `json:"source"`   // source
	Status  string `json:"status"`   // status open close
}

// GetEntConfigs ...
// GET /admin/api/v1/enterprise/configs
func (s *IMService) GetEntConfigs(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	configs, err := s.getConvConfigs(entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var errMsg *ErrMsg
	configs.SecurityConfig, errMsg = s.getSecurityConfig(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	configs.AutomaticMessages, err = models.AutomaticMessagesByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	configs.LeaveMessageConfig, err = models.LeaveMessageConfigByEntID(db.Mysql, entID)
	if err != nil {
		if err != sql.ErrNoRows {
			return dbErrResp(ctx, err.Error())
		}

		configs.LeaveMessageConfig = &models.LeaveMessageConfig{}
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: configs})
}

// GET /api/agent/enterprise_config
func (s *IMService) GetEntConfigsV1(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)

	configs, errMsg := s.getEnterpriseConfigs(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}
	return jsonResponse(ctx, configs)
}

func (s *IMService) getEnterpriseConfigs(entID string) (*adapter.EnterpriseConfigs, *ErrMsg) {
	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrMsg{Message: "租户不存在"}
		}

		return nil, &ErrMsg{Message: err.Error()}
	}

	var configs = &adapter.EnterpriseConfigs{}
	configsContent, err := models.GetConfigsFromCache(db.Mysql, entID)
	if err != nil || configsContent == "" {
		return configs, nil
	}

	if err = common.Unmarshal(configsContent, &configs); err != nil {
		return nil, &ErrMsg{Code: common.DecodeJSONErr, Message: err.Error()}
	}
	configs.IsActivated = ent.IsActivated

	return configs, nil
}

// POST /api/agent/enterprise_config
func (s *IMService) CreateOrUpdateEntConfigsV1(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	isDebug := conf.IMConf.Debug

	checkPerm := func(permName string) *ErrMsg {
		if !isDebug {
			if msg := hasPerm(entID, agentID, "online_agent_config", permName); msg != nil {
				return msg
			}
		}

		return nil
	}

	req := &adapter.EnterpriseConfigs{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	config, msg := s.getEnterpriseConfigs(entID)
	if msg != nil {
		return jsonResponse(ctx, msg)
	}

	if req.AgentsPermissionsConfig != nil {
		config.AgentsPermissionsConfig = req.AgentsPermissionsConfig
	}

	if req.AutoReplyMsgSettings != nil {
		config.AutoReplyMsgSettings = req.AutoReplyMsgSettings
	}

	if req.ChatLinkAutoMsgConfig != nil {
		config.ChatLinkAutoMsgConfig = req.ChatLinkAutoMsgConfig
	}

	if req.ClientWalkingAutoMsg != nil {
		config.ClientWalkingAutoMsg = req.ClientWalkingAutoMsg
	}

	if req.ConvGradeConfig != nil {
		config.ConvGradeConfig = req.ConvGradeConfig
	}

	if req.EndConvExpireConfig != nil {
		// config_chat_rule
		if msg := checkPerm("config_chat_rule"); msg != nil {
			return noPermResp(ctx, msg)
		}
		config.EndConvExpireConfig = req.EndConvExpireConfig
	}

	if req.EndingMsgSettings != nil {
		config.EndingMsgSettings = req.EndingMsgSettings
	}

	if req.InvitationConfig != nil {
		// config_invitation
		if msg := checkPerm("config_invitation"); msg != nil {
			return noPermResp(ctx, msg)
		}
		config.InvitationConfig = req.InvitationConfig
	}

	if req.OAuthSettings != nil {
		config.OAuthSettings = req.OAuthSettings
	}

	if req.PromotionMsgSettings != nil {
		config.PromotionMsgSettings = req.PromotionMsgSettings
	}

	if req.ReserveCluesConfig != nil {
		config.ReserveCluesConfig = req.ReserveCluesConfig
	}

	if req.QueueSettings != nil {
		// config_queuing
		if msg := checkPerm("config_queuing"); msg != nil {
			return noPermResp(ctx, msg)
		}
		config.QueueSettings = req.QueueSettings
	}

	if req.SendFileSettings != nil {
		config.SendFileSettings = req.SendFileSettings
	}

	if req.ServiceEvaluationConfig != nil {
		// config_evaluation
		if msg := checkPerm("config_evaluation"); msg != nil {
			return noPermResp(ctx, msg)
		}
		config.ServiceEvaluationConfig = req.ServiceEvaluationConfig
	}

	if req.StandaloneWindowConfig != nil {
		config.StandaloneWindowConfig = req.StandaloneWindowConfig
	}

	if req.Survey != nil {
		config.Survey = req.Survey
	}

	if req.TicketConfig != nil {
		config.TicketConfig = req.TicketConfig
	}

	if req.TimeoutRedirectConfig != nil {
		config.TimeoutRedirectConfig = req.TimeoutRedirectConfig
	}

	if req.VisitorVisible != nil {
		config.VisitorVisible = req.VisitorVisible
	}

	if req.WebCallbackSettings != nil {
		config.WebCallbackSettings = req.WebCallbackSettings
	}

	if req.WelcomeMsgSettings != nil {
		config.WelcomeMsgSettings = req.WelcomeMsgSettings
	}

	if req.WidgetSettings != nil {
		config.WidgetSettings = req.WidgetSettings
	}

	if errMsg := s.updateEntConfigs(entID, config); errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, config)
}

// POST /api/agent/enterprise_config/change_msg_status
func (s *IMService) ChangeMsgStatus(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req := &ChangeMsgStatusReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}
	if req.Status != "open" && req.Status != "close" {
		return invalidParameterResp(ctx, "status is not supported")
	}

	configs, errMsg := s.getConfigsFromCache(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	switch req.MsgType {
	case "auto_reply_msg_settings":
		settings := configs.AutoReplyMsgSettings
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	case "chat_link_auto_msg_config":
		settings := configs.ChatLinkAutoMsgConfig
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	case "client_waking_auto_msg":
		settings := configs.ClientWalkingAutoMsg
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	case "ending_msg_settings":
		settings := configs.EndingMsgSettings
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	case "promotion_msg_settings":
		settings := configs.PromotionMsgSettings
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	case "welcome_msg_settings":
		settings := configs.WelcomeMsgSettings
		switch req.Source {
		case "web":
			settings.Web.Status = req.Status
		}
	default:
		return invalidParameterResp(ctx, "not supported msg_type")
	}

	if errMsg := s.updateEntConfigs(entID, configs); errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, adapter.GetAutoMsgConfigsFromEntConfigs(configs))
}

func (s *IMService) updateEntConfigs(entID string, configs *adapter.EnterpriseConfigs) *ErrMsg {
	content, err := common.Marshal(configs)
	if err != nil {
		return &ErrMsg{Code: common.EncodeJSONErr, Message: err.Error()}
	}

	now := time.Now().UTC()
	config := &models.EntAllConfig{
		ID:            common.GenUniqueID(),
		EntID:         entID,
		ConfigContent: sql.NullString{Valid: true, String: content},
		CreateAt:      now,
		UpdateAt:      now,
	}
	if err = models.CreateOrUpdateConfigs(db.Mysql, config); err != nil {
		return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	db.RedisClient.HSet(fmt.Sprintf(common.EntConfigs, entID), common.EntConfigsContent, content)
	return nil
}

// UpdateEntConvConfigs ...
// PUT /admin/api/v1/enterprise/conv_configs
func (s *IMService) UpdateEntConvConfigs(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "agent_settings", "check_update_conversation_rule"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &UpdateConvConfigsReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	log.Logger.Println(req.ConvConfigs.VisitorQueueConf)

	configs := req.ConvConfigs
	if configs != nil {
		if configs.EndingMessageConf != nil {
			if err := models.InsertOrUpdateEndMessage(db.Mysql, entID, configs.EndingMessageConf); err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}

		if configs.EndingConversationConf != nil {
			if err := models.InsertOrUpdateEndConversation(db.Mysql, entID, configs.EndingConversationConf); err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}

		if configs.VisitorQueueConf != nil {
			if err := models.InsertOrUpdateQueueConfig(db.Mysql, entID, configs.VisitorQueueConf); err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}

		if configs.ConversationQualityConf != nil {
			if err := models.InsertOrUpdateConversationQuality(db.Mysql, entID, configs.ConversationQualityConf); err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}

		if configs.ConversationTransferConf != nil {
			if err := models.InsertOrUpdateConversationTransfer(db.Mysql, entID, configs.ConversationTransferConf); err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}
	}

	resp, err := s.getConvConfigs(entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: resp})
}

func (s *IMService) getConvConfigs(entID string) (configs *GetEntConfigsResp, err error) {
	configs = &GetEntConfigsResp{
		ConvConfigs: &ConvConfigs{},
	}

	configs.ConvConfigs.EndingConversationConf, err = models.EndingConversationByEntID(db.Mysql, entID)
	if err != nil && err != sql.ErrNoRows {
		return configs, err
	}

	configs.ConvConfigs.EndingMessageConf, err = models.EndingMessageByEntIDPlatform(db.Mysql, entID, models.EndingMessageWebPlatform)
	if err != nil && err != sql.ErrNoRows {
		return configs, err
	}

	configs.ConvConfigs.VisitorQueueConf, err = models.QueueConfigByEntID(db.Mysql, entID)
	if err != nil && err != sql.ErrNoRows {
		return configs, err
	}

	configs.ConvConfigs.ConversationTransferConf, err = models.ConversationTransferByEntID(db.Mysql, entID)
	if err != nil && err != sql.ErrNoRows {
		return configs, err
	}

	configs.ConvConfigs.ConversationQualityConf, err = models.ConversationQualityByEntID(db.Mysql, entID)
	if err != nil && err != sql.ErrNoRows {
		return configs, err
	}

	return configs, nil
}

// AddSecurityConfig
// POST /admin/api/v1/enterprise/security_configs
func (s *IMService) AddSecurityConfig(ctx echo.Context) (err error) {
	req := &SecurityConfig{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req.LoginLimit.EntID = entID
	req.SendFile.EntID = entID

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		dbErr = rollBackOrCommit(tx, dbErr)
		if dbErr != nil {
			log.Logger.Error("MySQL rollback/commit error: ", dbErr)
		}
	}()

	if dbErr = models.CreateLoginLimit(tx, req.LoginLimit); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	if dbErr = req.SendFile.Insert(tx); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: req})
}

// UpdateSecurityConfig
// PUT /admin/api/v1/enterprise/security_configs
func (s *IMService) UpdateSecurityConfig(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req := &SecurityConfig{}
	if err = ctx.Bind(req); err != nil {
		return
	}
	req.LoginLimit.EntID = entID
	req.SendFile.EntID = entID

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		dbErr = rollBackOrCommit(tx, dbErr)
		if dbErr != nil {
			log.Logger.Error("MySQL rollback/commit error: ", dbErr)
		}
	}()

	if dbErr = models.UpdateEntLoginLimit(tx, req.LoginLimit); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	if dbErr = req.SendFile.Update(tx); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: req})
}

// SecurityConfigByEntID
// GET /admin/api/v1/enterprise/security_configs
func (s *IMService) SecurityConfigByEntID(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	securityConfig, errMsg := s.getSecurityConfig(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: securityConfig})
}

func (s *IMService) getSecurityConfig(entID string) (*SecurityConfig, *ErrMsg) {
	limit, err := models.EntLoginLimitByEntID(entID)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	sendFile, err := models.SendFileByEntID(db.Mysql, entID)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	return &SecurityConfig{LoginLimit: limit, SendFile: sendFile}, nil
}

func (s *IMService) getConfigsFromCache(entID string) (configs *adapter.EnterpriseConfigs, errMsg *ErrMsg) {
	configs = &adapter.EnterpriseConfigs{}
	configsContent, err := models.GetConfigsFromCache(db.Mysql, entID)
	if err != nil || configsContent == "" {
		return configs, nil
	}

	if err = common.Unmarshal(configsContent, &configs); err != nil {
		return nil, &ErrMsg{common.DecodeJSONErr, err.Error()}
	}

	return configs, nil
}
