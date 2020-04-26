package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/external/submail"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

// POST /api/agent/agent_invitations
func (s *IMService) SendInvitation(ctx echo.Context) error {
	req := &dto.InvitationReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if len(req.Agents) == 0 {
		return jsonResponse(ctx, &SuccessResp{Success: true})
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentInvitation := req.Agents[0]

	if agentInvitation.PrivilegeRange == nil {
		return invalidParameterResp(ctx, "privilege_range invalid")
	}

	agent, err := models.AgentByEmail(db.Mysql, agentInvitation.Email)
	if err != nil && err != sql.ErrNoRows {
		return invalidParameterResp(ctx, err.Error())
	}

	if agent != nil {
		return invalidParameterResp(ctx, fmt.Sprintf("邮箱 %s 已存在", agentInvitation.Email))
	}

	agentLimit, err := models.GetAgentLimitByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentNum, err := models.AgentNumOfEnt(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if agentNum >= agentLimit {
		return invalidParameterResp(ctx, "坐席数已达上限")
	}

	serveLimit, err := models.GetAgentServeLimitByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if agentInvitation.ServingLimit > serveLimit {
		return invalidParameterResp(ctx, "坐席接待数超过限制")
	}

	now := time.Now().UTC()
	invitation := &models.AgentInvitation{
		ID:           common.GenUniqueID(),
		EnterpriseID: entID,
		GroupID:      agentInvitation.GroupID,
		Email:        agentInvitation.Email,
		Realname:     agentInvitation.Realname,
		ServingLimit: agentInvitation.ServingLimit,
		Status:       models.AgentInvitationPendingStatus,
		WorkNum:      agentInvitation.WorkNum,
		Privilege:    agentInvitation.Privilege,
		CreatedOn:    now,
		LastUpdated:  now,
		ExpiredOn:    now.Add(7 * 24 * time.Hour),
		AcceptedOn:   mysql.NullTime{},
	}

	switch v := agentInvitation.PrivilegeRange.(type) {
	case string:
		invitation.PrivilegeRange = v
	case []interface{}:
		if len(v) == 0 {
			return invalidParameterResp(ctx, "不合法的权限范围")
		}

		groups, err := common.Marshal(v)
		if err != nil {
			return errResp(ctx, common.EncodeJSONErr, err.Error())
		}
		invitation.PrivilegeRange = groups
	default:
		return invalidParameterResp(ctx, "unsupported privilege_range")
	}

	if err := invitation.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	go s.sendInvitationEmail(invitation.EnterpriseID, invitation.ID, invitation.Email)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// PUT /agent_invitations/:invitation_id
func (s *IMService) CancelInvitation(ctx echo.Context) error {
	req := &UpdateInvitationReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Status != models.AgentInvitationCancelledStatus {
		return invalidParameterResp(ctx, "status not valid")
	}

	invitationID := ctx.Param("invitation_id")
	invitation, err := models.AgentInvitationByID(db.Mysql, invitationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "无效的邀请")
		}
		return dbErrResp(ctx, err.Error())
	}

	d := time.Now().UTC().Sub(invitation.ExpiredOn)
	if invitation.Status == models.AgentInvitationSuccessStatus ||
		invitation.Status == models.AgentInvitationCancelledStatus || d.Seconds() > 0 {
		return invalidParameterResp(ctx, "邀请已激活/已取消/已过期")
	}

	if req.Status != "" {
		invitation.Status = req.Status
		if err := invitation.Update(db.Mysql); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	v, err := dto.ConvertInvitation(invitation)
	if err != nil {
		return jsonResponse(ctx, err)
	}

	return jsonResponse(ctx, v)
}

// GET /api/agent/agent_invitations
func (s *IMService) GetInvitations(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	invitations, err := models.AgentInvitationsByEnterpriseID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &GetInvitationsResp{}
	for _, invt := range invitations {
		v, err := dto.ConvertInvitation(invt)
		if err != nil {
			return jsonResponse(ctx, err)
		}

		resp.Invitations = append(resp.Invitations, v)
	}

	return jsonResponse(ctx, resp)
}

// POST /agent/accept_invitations
func (s *IMService) AcceptInvitation(ctx echo.Context) error {
	req := &AcceptInvitationReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	invitation, err := models.AgentInvitationByID(db.Mysql, req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "无效的邀请")
		}
		return dbErrResp(ctx, err.Error())
	}

	d := time.Now().UTC().Sub(invitation.ExpiredOn)
	if invitation.Status == models.AgentInvitationSuccessStatus || invitation.Status == models.AgentInvitationCancelledStatus || d.Seconds() > 0 {
		return invalidParameterResp(ctx, "邀请已激活/已取消/已过期")
	}

	pwd, err := common.GenHashedPassword([]byte(req.Password))
	if err != nil {
		return internalServerErr(ctx, "add agent error")
	}

	now := time.Now().UTC()
	agentModel := &models.Agent{
		ID:             common.GenUniqueID(),
		EntID:          invitation.EnterpriseID,
		GroupID:        invitation.GroupID,
		RoleID:         invitation.Privilege,
		Avatar:         "",
		Username:       invitation.Realname,
		RealName:       invitation.Realname,
		NickName:       req.Nickname,
		HashedPassword: string(pwd),
		JobNumber:      invitation.WorkNum,
		ServeLimit:     invitation.ServingLimit,
		Email:          invitation.Email,
		Status:         models.AgentUnavailableStatus,
		IsAdmin:        0,
		AccountStatus:  models.AgentAccountValidStatus,
		CreateAt:       now,
		UpdateAt:       now,
	}

	var groups []string
	if invitation.PrivilegeRange == models.AgentPermsRangePersonalType || invitation.PrivilegeRange == models.AgentPermsRangeAllType {
		agentModel.PermsRangeType = invitation.PrivilegeRange
	} else {
		err := common.Unmarshal(invitation.PrivilegeRange, &groups)
		if err != nil {
			log.Logger.Warnf("Unmarshal PrivilegeRange: %s, error: %v", invitation.PrivilegeRange, err)
			agentModel.PermsRangeType = models.AgentPermsRangePersonalType
		}
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if err = agentModel.Insert(tx); err != nil {
		tx.Rollback()
		return dbErrResp(ctx, err.Error())
	}

	if len(groups) > 0 {
		if err = models.BulkAddAgentPermsRangeGroups(tx, agentModel.ID, groups); err != nil {
			tx.Rollback()
			return dbErrResp(ctx, err.Error())
		}
	}

	if err = models.UpdateAgentInvitationStatus(tx, req.Token, models.AgentInvitationSuccessStatus); err != nil {
		tx.Rollback()
		return dbErrResp(ctx, err.Error())
	}

	tx.Commit()

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// PUT /api/agent/resend_invitations/:invitation_id
func (s *IMService) ResendInvitation(ctx echo.Context) error {
	invitationID := ctx.Param("invitation_id")
	invitation, err := models.AgentInvitationByID(db.Mysql, invitationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "无效的邀请")
		}
		return dbErrResp(ctx, err.Error())
	}

	d := time.Now().UTC().Sub(invitation.ExpiredOn)
	if invitation.Status == models.AgentInvitationSuccessStatus ||
		invitation.Status == models.AgentInvitationCancelledStatus || d.Seconds() > 0 {
		return invalidParameterResp(ctx, "邀请已激活/已取消/已过期")
	}

	go s.sendInvitationEmail(invitation.EnterpriseID, invitation.ID, invitation.Email)
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func (s *IMService) sendInvitationEmail(entID, invitationID, email string) {
	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		log.Logger.Warnf("[sendInvitationEmail] error: %v", err)
		return
	}

	//content := fmt.Sprintf(
	//	"%s/activate-account?token=%s&ent_id=%s&fullname=%s&email=%s",
	//	conf.IMConf.Host,
	//	invitationID,
	//	entID,
	//	ent.FullName,
	//	email,
	//)

	link := fmt.Sprintf(conf.IMConf.Host+"/activate-account?token=%s&ent_id=%s&fullname=%s&email=%s",
		invitationID,
		entID,
		ent.FullName,
		email,
	)
	log.Logger.Infof("[SendInvitation] link: %s", link)

	err = submail.SendActivateEmail(email, `激活chat186客服帐号(请点击链接完成账号激活)`, "", link)
	if err != nil {
		log.Logger.Warnf("[SendInvitation] send invitation email error: %v", err)
	}
	return
}

// POST /api/resend_activate_account_email
// ResendActivateEmail ...
func (s *IMService) ResendActivateEmail(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if ent.IsActivated {
		return invalidParameterResp(ctx, "已激活")
	}

	req := ctx.Request()
	source := req.Referer()
	if err := sendActiveEntEmail(db.RedisClient, ent.Email, entID, ent.AdminID, source); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}
