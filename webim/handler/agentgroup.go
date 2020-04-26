package handler

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type CreateAgentGroupReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateAgentGroupResp struct {
	EnterpriseID string `json:"enterprise_id"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Token        string `json:"token"`
}

type UpdateAgentGroupReq struct {
	Name string `json:"name"`
}

// CreateAgentGroup ...
// POST /admin/api/v1/enterprise/agent_groups
// POST /api/agent/agent_groups
func (s *IMService) CreateAgentGroup(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "ent_info", "create_change_delete_group"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &CreateAgentGroupReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Name == "" {
		return invalidParameterResp(ctx, "name is invalid")
	}

	group := &models.AgentGroup{
		ID:          common.GenUniqueID(),
		EntID:       entID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := group.Insert(db.Mysql); err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return dbErrResp(ctx, err.Error())
		}

		if mysqlErr.Number == 1062 {
			return invalidParameterResp(ctx, "name is exists")
		}
	}

	return jsonResponse(ctx, &CreateAgentGroupResp{
		EnterpriseID: entID,
		ID:           group.ID,
		Name:         group.Name,
		Token:        group.ID,
	})
}

// DELETE /api/agent/agent_groups/:group_id
func (s *IMService) DeleteAgentGroup(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "ent_info", "create_change_delete_group"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	groupID := ctx.Param("group_id")
	if groupID == "" {
		return invalidParameterResp(ctx, "group_id invalid")
	}

	group, err := models.AgentGroupByID(db.Mysql, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return jsonResponse(ctx, "Group Not Exists!")
		}

		return dbErrResp(ctx, err.Error())
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		if err := rollBackOrCommit(tx, dbErr); err != nil {
			log.Logger.Warnf("[DeleteAgentGroup] error: %v", err)
		}
	}()

	if err := group.Delete(tx); err != nil {
		dbErr = err
		return dbErrResp(ctx, err.Error())
	}

	ag := models.AgentGroupRelation{GroupID: groupID}
	if err := ag.DeleteByGroupID(tx); err != nil {
		dbErr = err
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// PUT /api/agent/agent_groups/:group_id
func (s *IMService) UpdateAgentGroup(ctx echo.Context) error {
	req := &UpdateAgentGroupReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	groupID := ctx.Param("group_id")
	group, err := models.AgentGroupByID(db.Mysql, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "group not found")
		}

		return dbErrResp(ctx, err.Error())
	}

	if req.Name != "" {
		group.Name = req.Name
		if err := group.Update(db.Mysql); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}
