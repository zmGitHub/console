package handler

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AddPermsToRoleReq struct {
	PermIDs []string `json:"perm_ids"`
}

type GetAgentPermsResp struct {
	PermsRange string         `json:"perms_range"`
	Perms      []*models.Perm `json:"perms"`
}

type updateRoleReq struct {
	Name string `json:"name"`
}

// AddRole
// POST /admin/api/v1/roles
// POST /api/perm/roles
func (s *IMService) AddRole(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "ent_info", "config_role"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &models.Role{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Name == "" {
		return invalidParameterResp(ctx, "name is invalid")
	}

	req.ID = common.GenUniqueID()
	req.EntID = entID
	if err = req.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &adapter.Role{Name: req.Name, Token: req.ID})
}

// PUT /api/perm/roles/:role_id
func (s *IMService) UpdateRole(ctx echo.Context) error {
	roleID := ctx.Param("role_id")
	req := &updateRoleReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Name != "" {
		entID := ctx.Get(middleware.AgentEntIDKey).(string)
		role := models.Role{ID: roleID, EntID: entID, Name: req.Name}
		if err := role.Update(db.Mysql); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	return jsonResponse(ctx, adapter.Role{Name: req.Name, Token: roleID})
}

// DELETE /api/perm/roles/:role_id
func (s *IMService) DeleteRole(ctx echo.Context) error {
	roleID := ctx.Param("role_id")

	role := models.Role{ID: roleID}
	if err := role.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, nil)
}

// AddOrUpdatePermsToRole
// POST /admin/api/v1/roles/:role_id/perms
// PUT /api/perm/roles/:role_id/perms
func (s *IMService) AddOrUpdatePermsToRole(ctx echo.Context) (err error) {
	roleID := ctx.Param("role_id")
	if roleID == "" {
		return invalidParameterResp(ctx, "role_id is invalid")
	}

	req := &AddPermsToRoleReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	updatePerms := func() error {
		tx, err := db.Mysql.Begin()
		if err != nil {
			return err
		}

		var dbErr error
		defer func() {
			if dbErr = rollBackOrCommit(tx, dbErr); dbErr != nil {
				log.Logger.Warnf("rollBackOrCommit error: %v", err)
			}
		}()

		if dbErr = models.UpdateRolePerms(tx, roleID, req.PermIDs); dbErr != nil {
			return dbErr
		}

		return nil
	}
	if err := updatePerms(); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	perms, errMsg := s.getRolePerms(roleID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	agents, err := models.AgentIDsByRoleID(db.Mysql, roleID)
	if err == nil {
		var keys []string
		for _, id := range agents {
			keys = append(keys, fmt.Sprintf(common.AgentPerms, id))
		}
		db.RedisClient.Del(keys...)
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	go s.sendPermsUpdateEvent(entID, roleID)

	return jsonResponse(ctx, perms)
}

// GetAgentPerms
// GET /admin/api/v1/agents/:agent_id/perms
// GET /api/agent/agents/:agent_id/perms_ranges
func (s *IMService) GetAgentPerms(ctx echo.Context) (err error) {
	agentID := ctx.Param("agent_id")
	if agentID == "" {
		return invalidParameterResp(ctx, "agent_id is invalid")
	}

	perms, errMsg := s.getAgentPerms(agentID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, perms)
}

func (s *IMService) getAgentPerms(agentID string) (*adapter.PermsRangeResp, *ErrMsg) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrMsg{common.UserNotExistErr, "agent not exists"}
		}
		return nil, &ErrMsg{common.DBErr, err.Error()}
	}

	var permsRange interface{} = agent.PermsRangeType
	var groups []string
	if agent.PermsRangeType != models.AgentPermsRangePersonalType && agent.PermsRangeType != models.AgentPermsRangeAllType {
		if strings.HasPrefix(agent.PermsRangeType, "[") {
			err := common.Unmarshal(agent.PermsRangeType, &groups)
			if err != nil {
				permsRange = []string{}
			} else {
				permsRange = groups
			}
		}
	}

	perms, err := models.GetAgentPerms(db.Mysql, agentID)
	if err != nil {
		return nil, &ErrMsg{common.DBErr, err.Error()}
	}

	return adapter.ConvertAgentPermsToPermsRangesResp(permsRange, perms), nil
}

// GetRolePerms
// GET /admin/api/v1/roles/:role_id/perms
// GET /api/perm/roles/:role_id/perms
func (s *IMService) GetRolePerms(ctx echo.Context) (err error) {
	roleID := ctx.Param("role_id")
	if roleID == "" {
		return invalidParameterResp(ctx, "role_id is invalid")
	}

	perms, errMsg := s.getRolePerms(roleID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, perms)
}

func (s *IMService) getRolePerms(roleID string) (adapter.Perms, *ErrMsg) {
	rolePerms, err := models.RolePermsByRoleID(db.Mysql, roleID)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	var permIDs []string
	for _, p := range rolePerms {
		permIDs = append(permIDs, p.PermID)
	}

	perms, err := models.PermsByPermIDs(db.Mysql, permIDs)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	var permsMap = make(adapter.Perms, len(perms))
	for _, p := range perms {
		perm := &adapter.Perm{ID: p.ID, Key: p.Name}
		if v, ok := permsMap[p.AppName]; ok {
			v = append(v, perm)
			permsMap[p.AppName] = v
		} else {
			permsMap[p.AppName] = []*adapter.Perm{perm}
		}
	}

	return permsMap, nil
}

// GET /api/perm/roles
func (s *IMService) GetEntRoles(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	roles, err := models.RolesByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, adapter.ConvertRolesToAdapterRoles(roles))
}

func (s *IMService) sendPermsUpdateEvent(entID, roleID string) {
	agents, err := models.AgentsByRoleID(db.Mysql, roleID)
	if err != nil {
		return
	}

	event := &PermUpdateEvent{
		Action:       "update_perms",
		EnterpriseID: entID,
		ID:           roleID,
		TraceStart:   -1,
	}
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	var wg sync.WaitGroup

	for _, agent := range agents {
		wg.Add(1)
		go func(agent *models.Agent) {
			sendMessageToAgent(s.imCli, agent.ID, eventContent)
			wg.Done()
		}(agent)
	}
	wg.Wait()
}
