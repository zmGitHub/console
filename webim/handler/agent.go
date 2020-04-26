package handler

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"bitbucket.org/forfd/gooxml/spreadsheet"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AgentsResp struct {
	*models.Agent
	Status string `json:"status"`
}

// {"new_password":"abc123@#","password":"bob@chat186.com","repeat_password":"abc123@#"}
type UpdatePasswordReq struct {
	Password       *string `json:"password"`
	NewPassword    *string `json:"new_password"`
	RepeatPassword *string `json:"repeat_password"`
}

type UpdateAgentInfoV1Req struct {
	Nickname        string      `json:"nickname"`
	Email           string      `json:"email"`
	Avatar          string      `json:"avatar"`
	Signature       string      `json:"signature"`
	Telephone       string      `json:"telephone"`
	PublicCellphone string      `json:"public_cellphone"`
	PublicEmail     string      `json:"public_email"`
	Qq              string      `json:"qq"`
	Weixin          string      `json:"weixin"`
	Realname        string      `json:"realname"`
	GroupID         string      `json:"group_id"`
	Privilege       string      `json:"privilege"`
	PrivilegeRange  interface{} `json:"privilege_range"`
	ServingLimit    *int        `json:"serving_limit"`
	WorkNum         string      `json:"work_num"`
	Position        *int        `json:"position"`
	*UpdatePasswordReq
}

type UpdateAgentStatusReq struct {
	Status string `json:"status"`
}

type UpdateAgentEmailReq struct {
	Email string `json:"email"`
}

type agentRanking struct {
	AgentID string `json:"agent_id"`
	Ranking int    `json:"ranking"`
}

type UpdateAgentsRankingReq struct {
	Rankings []*agentRanking `json:"rankings"`
}

type GetInvitationsResp struct {
	Invitations []*dto.Invitation `json:"invitations"`
}

// // '{"token":"44c107855793a65de11047e92f75b4a749a512e17d9833a24b93f09f40489d19","password":"hanxu317","nickname":"hanxu317"}
type AcceptInvitationReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type UpdateInvitationReq struct {
	Status string `json:"status"` // cancelled
}

type ConfirmAgentEmailReq struct {
	Token string `json:"token"`
}

// GetEntAgents
// GET /admin/api/v1/enterprise/agents
// GET /api/agent/agents
func (s *IMService) GetEntAgents(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	currentAgentID := ctx.Get(middleware.AgentIDKey).(string)

	agents, err := models.AgentsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var currentAgent *models.Agent
	var allAgentIDs []string
	for _, agent := range agents {
		if agent.ID == currentAgentID {
			currentAgent = agent
		}

		allAgentIDs = append(allAgentIDs, agent.ID)
	}

	switch currentAgent.PermsRangeType {
	case models.AgentPermsRangePersonalType:
		allAgentIDs = []string{currentAgentID}
	case models.AgentPermsRangeAllType:
	default:
		allAgentIDs, err = s.getGroupAgents(currentAgentID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	onlineAgents := models.GetEntOnlineAgentsV1(entID)
	sort.SliceStable(agents, func(i, j int) bool {
		return agents[i].Ranking <= agents[j].Ranking
	})

	resultAgents := adapter.ConvertAgentsToAdapterAgents(agents)
	groups, err := models.PermsRangeGroupIDsByAgentIDs(db.Mysql, allAgentIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentMap := map[string]struct{}{}
	for _, id := range allAgentIDs {
		agentMap[id] = struct{}{}
	}

	for _, agent := range resultAgents {
		_, ok := agentMap[agent.ID]
		if !ok {
			continue
		}

		if v, ok := groups[agent.ID]; ok {
			agent.PrivilegeRange = v
		}

		for _, id := range onlineAgents {
			if agent.ID == id {
				agent.IsOnline = true
				break
			}
		}
	}

	return jsonResponse(ctx, resultAgents)
}

// PUT /api/agent/agents/:agent_id
func (s *IMService) UpdateAgentInfoV1(ctx echo.Context) (err error) {
	currentAgent := ctx.Get(middleware.AgentIDKey).(string)
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	if !conf.IMConf.Debug {
		if errMsg := hasPerm(entID, currentAgent, "ent_info", "create_change_delete_agent"); errMsg != nil {
			return noPermResp(ctx, errMsg)
		}
	}

	req := &UpdateAgentInfoV1Req{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	agentID := ctx.Param("agent_id")
	if agentID == "" {
		return invalidParameterResp(ctx, "agent_id invalid")
	}

	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if req.Email != "" {
		header := ctx.Request().Header

		// Captcha-Token: bseOta3OME0oDPVWzOAy
		// Captcha-Value: 341131
		token := header.Get("Captcha-Token")
		if token == "" {
			return invalidParameterResp(ctx, "invalid captcha token")
		}

		// captcha-value
		value := header.Get("Captcha-Value")
		if value == "" {
			return invalidParameterResp(ctx, "invalid captcha value")
		}

		if !captcha.VerifyString(token, value) {
			return invalidParameterResp(ctx, "图形验证码错误,请重新输入")
		}

		if currentAgent != agentID {
			return invalidParameterResp(ctx, "")
		}

		if err := s.UpdateAgentEmail(agentID, agent.NickName, req.Email); err != nil {
			return internalServerErr(ctx, err.Error())
		}

		return jsonResponse(ctx, adapter.ConvertAgentToAgentInfo(agent, true))
	}

	if req.Avatar != "" {
		agent.Avatar = req.Avatar
	}

	if req.Nickname != "" {
		agent.NickName = req.Nickname
	}

	if req.Signature != "" {
		agent.Signature = req.Signature
	}

	if req.Telephone != "" {
		agent.Mobile = req.Telephone
	}

	if req.PublicCellphone != "" {
		agent.PublicTelephone = req.PublicCellphone
	}

	if req.PublicEmail != "" {
		agent.PublicEmail = req.PublicEmail
	}

	if req.Weixin != "" {
		agent.Wechat = req.Weixin
	}

	if req.Qq != "" {
		agent.QqNum = req.Qq
	}

	//if req.Realname != "" {
	//	agent.RealName = req.Realname
	//}

	if req.GroupID != "" {
		agent.GroupID = req.GroupID
	}

	if req.Privilege != "" {
		agent.RoleID = req.Privilege
	}

	var newGroups []string
	var updateGroups bool
	if req.PrivilegeRange != nil {
		switch v := req.PrivilegeRange.(type) {
		case string:
			if !(v == models.AgentPermsRangeAllType || v == models.AgentPermsRangePersonalType) {
				return invalidParameterResp(ctx, "unsupported privilege_range")
			}

			agent.PermsRangeType = v
			updateGroups = true
		case []interface{}:
			if len(v) == 0 {
				return invalidParameterResp(ctx, "invalid privilege_range")
			}

			for _, id := range v {
				if id != nil {
					if groupID, ok := id.(string); ok {
						newGroups = append(newGroups, groupID)
					}
				}
			}
		default:
			return invalidParameterResp(ctx, "invalid privilege_range")
		}
	}

	if req.ServingLimit != nil {
		lmt := *req.ServingLimit
		serveLimit, err := models.GetAgentServeLimitByEntID(db.Mysql, entID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		if lmt > serveLimit {
			return invalidParameterResp(ctx, "坐席数超过上限")
		}

		agent.ServeLimit = lmt
	}

	//if req.WorkNum != "" {
	//	agent.JobNumber = req.WorkNum
	//}

	if req.Position != nil {
		agent.Ranking = *req.Position
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			log.Logger.Warnf("[UpdateAgentInfoV1] panic: %v", e)
		}

		rollBackOrCommit(tx, dbErr)
	}()

	agent.UpdateAt = time.Now().UTC()
	if dbErr = agent.Update(tx); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	if len(newGroups) > 0 || updateGroups {
		if dbErr = models.UpdateAgentPermGroups(tx, agent.ID, newGroups); dbErr != nil {
			return dbErrResp(ctx, dbErr.Error())
		}
	}

	if req.UpdatePasswordReq != nil {
		updatePwdReq := req.UpdatePasswordReq
		if updatePwdReq.Password != nil && updatePwdReq.NewPassword != nil && updatePwdReq.RepeatPassword != nil {
			if *req.NewPassword != *req.RepeatPassword {
				dbErr = fmt.Errorf("密码不一致")
				return invalidParameterResp(ctx, "密码不一致")
			}

			if err = bcrypt.CompareHashAndPassword([]byte(agent.HashedPassword), []byte(*req.Password)); err != nil {
				return invalidParameterResp(ctx, "原密码错误")
			}

			if dbErr = s.updatePassword(tx, agentID, *req.NewPassword); dbErr != nil {
				return internalServerErr(ctx, dbErr.Error())
			}
		}
	}

	return jsonResponse(ctx, adapter.ConvertAgentToAgentInfo(agent, true))
}

func (s *IMService) updatePassword(db models.XODB, agentID string, newPassword string) error {
	v, err := common.GenHashedPassword([]byte(newPassword))
	if err != nil {
		return err
	}

	if err := models.UpdateAgentPassword(db, agentID, string(v)); err != nil {
		return err
	}

	return nil
}

// GetAgentByID ...
// GET /admin/api/v1/agents/:agent_id
func (s *IMService) GetAgentByID(ctx echo.Context) (err error) {
	agentID := ctx.Param("agent_id")
	if agentID == "" {
		return invalidParameterResp(ctx, "agent_id is invalid")
	}

	viewAgent, errMsg := s.getAgentInfo(agentID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: viewAgent})
}

// GET /admin/api/v1/agents/info
// GET /api/agent/info
func (s *IMService) GetCurrentAgentInfo(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	viewAgent, errMsg := s.getAgentInfo(agentID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, adapter.ConvertAgentToAgentInfo(viewAgent.Agent, true))
}

func (s *IMService) getAgentInfo(agentID string) (viewAgent *viewAgent, errMsg *ErrMsg) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrMsg{common.UserNotExistErr, "agent not exists/deleted"}
		}
		return nil, &ErrMsg{common.DBErr, err.Error()}
	}
	agent.HashedPassword = ""

	status, err := s.imCli.IsOnline(context.Background(), agentID)
	if err != nil {
		log.Logger.Warnf("get agent online status error: %v", err)
		agent.IsOnline = false
	}

	if status {
		agent.IsOnline = true
	} else {
		agent.IsOnline = false
	}

	return modelAgentToViewAgent(agent), nil
}

// GetEntAgentGroups
// GET /admin/api/v1/enterprise/agent_groups
// GET /api/agent/agent_groups
func (s *IMService) GetEntAgentGroups(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	//agentID := ctx.Get(middleware.AgentIDKey).(string)
	//if !conf.IMConf.Debug {
	//	if msg := hasPerm(entID, agentID, "ent_info", "see_groups"); msg != nil {
	//		return noPermResp(ctx, msg)
	//	}
	//}

	groups, err := models.AgentGroupsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	return jsonResponse(ctx, adapter.ConvertAgentGroupsToAdapterGroups(groups))
}

// UpdateAgentStatus
// PUT /admin/api/v1/agents/status
// PUT /api/agent/agents/:agent_id/status
func (s *IMService) UpdateAgentStatus(ctx echo.Context) (err error) {
	agentID := ctx.Param("agent_id")
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	currentAgent := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug && currentAgent != agentID {
		if msg := hasPerm(entID, currentAgent, "visitor_and_conv", "change_agent_online_status"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &UpdateAgentStatusReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Status != models.AgentAvailableStatus && req.Status != models.AgentUnavailableStatus {
		return invalidParameterResp(ctx, "unsupported status")
	}

	if err = models.UpdateAgentStatus(db.Mysql, agentID, req.Status); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	go s.sendAgentStatusUpdateEvent(agentID, req.Status, true)

	return jsonResponse(ctx, &Resp{Code: 0})
}

// ExportAgents
// GET /admin/api/v1/agents/export
func (s *IMService) ExportAgents(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	if entID == "" {
		return invalidParameterResp(ctx, "ent_id is invalid")
	}

	agents, err := models.AgentsByEntID(db.Mysql, entID)
	if err != nil {
		return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: err.Error()})
	}

	name, err := s.saveAgents(entID, agents)
	if err != nil {
		return errResp(ctx, common.ExportFileErr, err.Error())
	}

	defer func() {
		if err := os.Remove(name); err != nil {
			log.Logger.Errorf("Remove file: %s, error: %v", name, err)
		}
	}()

	bs, err := ioutil.ReadFile(name)
	if err != nil {
		return errResp(ctx, common.ExportFileErr, err.Error())
	}

	fileName := fmt.Sprintf("%s/files/%s", entID, name)
	url, err := s.uploader.Upload(fileName, bytes.NewBuffer(bs), time.Now().Add(defaultExpires))
	if err != nil {
		return errResp(ctx, common.UploadFileErr, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: map[string]interface{}{"location": url}})
}

func (s *IMService) saveAgents(entID string, agents []*models.Agent) (name string, err error) {
	if len(agents) == 0 {
		return
	}
	var roleIDs []string
	for _, agent := range agents {
		roleIDs = append(roleIDs, agent.RoleID)
	}

	names, err := models.RoleNamesByRoleIDs(db.Mysql, roleIDs)
	if err != nil {
		return
	}

	ss := spreadsheet.New()
	sheet := ss.AddSheet()

	cells := []string{
		"客服ID",
		"工号",
		"名称",
		"账号",
		"角色",
		"权限范围",
		"昵称",
		"服务上限",
	}
	row := sheet.AddRow()
	for _, cellName := range cells {
		cell := row.AddCell()
		cell.SetString(cellName)
	}

	for _, agent := range agents {
		row := sheet.AddRow()
		roleName := names[agent.RoleID]
		permsRange := models.AgentPermsRangeType[agent.PermsRangeType]
		values := []string{
			agent.ID,
			agent.JobNumber,
			agent.RealName,
			agent.Email,
			roleName,
			permsRange,
			agent.Username,
			fmt.Sprintf("%d", agent.ServeLimit),
		}
		for _, v := range values {
			cell := row.AddCell()
			cell.SetString(v)
		}
	}

	if err = ss.Validate(); err != nil {
		return
	}

	name = fmt.Sprintf("/tmp/%s-%d-agent-list.xlsx", entID, time.Now().Unix())
	if err = ss.SaveToFile(name); err != nil {
		return
	}

	return
}

// DeleteAgent
// DELETE /admin/api/v1/agents/:agent_id
func (s *IMService) DeleteAgent(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if errMsg := hasPerm(entID, agentID, "", "create_change_delete_agent"); errMsg != nil {
			return jsonResponse(ctx, errMsg)
		}
	}

	deleteAgentID := ctx.Param("agent_id")
	if deleteAgentID == "" {
		return invalidParameterResp(ctx, "agent_id invalid")
	}

	if err = models.DeleteAgent(db.Mysql, deleteAgentID); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	db.RedisClient.Del(fmt.Sprintf(common.AgentPerms, deleteAgentID))
	db.RedisClient.Del(fmt.Sprintf(common.AgentConversationNum, deleteAgentID))
	db.RedisClient.Del(fmt.Sprintf(common.AgentServeLimit, deleteAgentID))

	agentTokensKey := fmt.Sprintf(common.AgentTokenList, agentID)
	tokens, err := db.RedisClient.SMembers(agentTokensKey).Result()
	if err != nil {
		log.Logger.Warnf("get SMembers from %s, error: %v", agentTokensKey, err)
	}

	var userTokens []interface{}
	for _, token := range tokens {
		userTokens = append(userTokens, token)
	}
	db.RedisClient.ZRem(fmt.Sprintf(common.EntOnlineAgentList, entID), userTokens...)

	db.RedisClient.Del(fmt.Sprintf(common.AgentTokenList, deleteAgentID))
	db.RedisClient.Del(fmt.Sprintf(common.AgentLoginCount, deleteAgentID))
	db.RedisClient.ZRem(fmt.Sprintf(common.EntAgents, entID), deleteAgentID)

	return jsonResponse(ctx, &Resp{Code: 0})
}

// UpdateAgentEmail
// // %s/change-email?token=%s&nickname=%s&email=%s
func (s *IMService) UpdateAgentEmail(agentID, agentNickName, email string) (err error) {
	randomStr := common.RandStringBytesMask(32)
	duration := conf.IMConf.AgentConf.ActivateCodeEffectiveDuration.Duration
	_, err = db.RedisClient.HMSet(randomStr, map[string]interface{}{
		"agent_id": agentID,
		"email":    email,
	}).Result()
	if err != nil {
		log.Logger.Errorf("HMSet activate code key(%s) in redis error: %v", randomStr, err)
		return
	}

	if err = db.RedisClient.Expire(randomStr, duration).Err(); err != nil {
		log.Logger.Error("Expire activate code key in redis error: ", err)
		return
	}

	title := `更新邮件地址`
	content := fmt.Sprintf("%s/change-email?token=%s&nickname=%s&email=%s", conf.IMConf.Host, randomStr, agentNickName, email)
	if err = s.mailClient.SendEmail(email, title, content); err != nil {
		log.Logger.Warnf("send email: %s,  error: %v", email, err)
		return
	}

	return nil
}

// POST /api/agent/confirm_email_change
func (s *IMService) ConfirmAgentEmail(ctx echo.Context) (err error) {
	req := &ConfirmAgentEmailReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Token == "" {
		return invalidParameterResp(ctx, "invalid token")
	}

	fields, err := db.RedisClient.HMGet(req.Token, "agent_id", "email").Result()
	if err != nil {
		log.Logger.Warnf("get key: %s, error: %v", req.Token, err)
		return internalServerErr(ctx, "internal server error")
	}

	agentID, agentEmail := fields[0], fields[1]
	if agentID != nil && agentEmail != nil {
		email := agentEmail.(string)
		agent, err := models.AgentByEmail(db.Mysql, email)
		if err != nil && err != sql.ErrNoRows {
			return internalServerErr(ctx, "internal server error")
		}

		if agent != nil {
			return invalidParameterResp(ctx, fmt.Sprintf("Email %s 已存在", email))
		}

		if err := models.UpdateAgentEmail(db.Mysql, agentID.(string), email); err != nil {
			log.Logger.Warnf("[UpdateAgentEmail] error: %v", err)
			return internalServerErr(ctx, "email 更新失败")
		}
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})

}

// GET /api/v1/activate_email?activate_code=xxxxx
func (s *IMService) ActivateAgentEmail(ctx echo.Context) (err error) {
	activateCode := ctx.QueryParam("activate_code")
	if activateCode == "" {
		return invalidParameterResp(ctx, "activate_code invalid")
	}

	v, err := db.RedisClient.HMGet(activateCode, "agent_id", "email").Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "invalid activate_code")
		}

		log.Logger.Error("get email from redis err: ", err)
		return dbErrResp(ctx, err.Error())
	}

	if len(v) < 2 {
		return invalidParameterResp(ctx, "invalid activate_code")
	}

	if v[0] == nil || v[1] == nil {
		return invalidParameterResp(ctx, "invalid activate_code")
	}

	agentID, email := v[0].(string), v[1].(string)
	if dbErr := models.UpdateAgentEmail(db.Mysql, agentID, email); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

// UpdateAgentsRanking 调整坐席顺序
// PUT /admin/api/v1/agents/rankings
func (s *IMService) UpdateAgentsRanking(ctx echo.Context) (err error) {
	req := &UpdateAgentsRankingReq{}
	if err = ctx.Bind(req); err != nil {
		return errResp(ctx, common.InvalidParameterErr, err.Error())
	}

	if len(req.Rankings) == 0 {
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	var rankings = make(map[string]int, len(req.Rankings))
	for _, ranking := range req.Rankings {
		rankings[ranking.AgentID] = ranking.Ranking
	}

	if err = models.UpdateAgentsRanking(db.Mysql, rankings); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentRankings, err := models.AgentRankingsByEntID(db.Mysql, entID)
	if err == nil && len(agentRankings) > 0 {
		var members []redis.Z
		for _, ranking := range agentRankings {
			members = append(members, redis.Z{
				Score:  float64(ranking.Ranking),
				Member: ranking.AgentID,
			})
		}

		key := fmt.Sprintf(common.EntAgents, entID)
		err = db.RedisClient.ZAdd(key, members...).Err()
		if err != nil {
			log.Logger.Warn("ZAdd ent agents error: ", err)
		}
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

// POST /api/agent/kick_person_offline
func (s *IMService) KickOffline(ctx echo.Context) (err error) {
	type kickOffline struct {
		AgentID string `json:"agent_id"`
	}
	req := &kickOffline{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.AgentID == "" {
		return invalidParameterResp(ctx, "agent_id invalid")
	}

	tkListKey := fmt.Sprintf(common.AgentTokenList, req.AgentID)
	tokens, err := db.RedisClient.SMembers(tkListKey).Result()
	if err != nil {
		log.Logger.Warnf("Get Agent Token List(%s), error: %v", tkListKey, err)
	}

	s.sendAgentKickedEvent(req.AgentID, tokens)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func (s *IMService) getOnlineAgents(entID string) []string {
	agentIDs, err := models.AgentIDsByEntID(db.Mysql, entID)
	if err != nil {
		return []string{}
	}

	return models.FilterOnline(agentIDs)
}

func (s *IMService) getGroupAgents(agentID string) (agentIDs []string, err error) {
	groups, err := models.PermsRangeGroupIDsByAgentID(db.Mysql, agentID)
	if err != nil {
		return
	}

	if len(groups) > 0 {
		agentIDs, err = models.AgentIDsByPermGroupIDs(db.Mysql, groups)
		if err != nil {
			return
		}
	}

	return
}
