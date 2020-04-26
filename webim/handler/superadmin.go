package handler

import (
	"database/sql"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/echo/v4"

	log "bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type QueryAllEnterpriseReq struct {
	Type   string `query:"type"` // all, paid, unpaid
	Email  string `query:"email"`
	Offset int    `query:"offset"`
	Limit  int    `query:"limit"`
}

type OverAllEntInfo struct {
	Enterprise *models.Enterprise `json:"enterprise"`
	Plan       *models.EntPlan    `json:"plan"`
}

type QueryAllEnterpriseResp struct {
	Total          int64             `json:"total"`
	AllEnterprises []*OverAllEntInfo `json:"all_enterprises"`
}

type StartTrialReq struct {
	Email      string `json:"email"`
	TrialInDay int    `json:"trial_in_day"`
}

type UpgradeEnterpriseReq struct {
	Email         string `json:"email"`
	Plan          int    `json:"plan"`
	AgentNum      int    `json:"agent_num"`
	DurationInDay int    `json:"duration_in_day"`
}

// GET /super_admin/enterprises
func (s *IMService) QueryAllEnterprise(ctx echo.Context) (err error) {
	req := &QueryAllEnterpriseReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var offset, limit = 0, 30
	if req.Offset > 0 {
		offset = req.Offset
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	total, ents, err := models.QueryEnterprises(db.Mysql, offset, limit, req.Type, req.Email)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	result := &QueryAllEnterpriseResp{Total: total}
	for _, ent := range ents {
		result.AllEnterprises = append(result.AllEnterprises, &OverAllEntInfo{
			Enterprise: ent,
		})
	}
	return jsonResponse(ctx, result)
}

// POST /super_admin/trials
func (s *IMService) StartTrial(ctx echo.Context) (err error) {
	req := &StartTrialReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Email == "" {
		return invalidParameterResp(ctx, "email 为空")
	}

	trialSet := mapset.NewSet(1, 3, 7)
	if !trialSet.Contains(req.TrialInDay) {
		return invalidParameterResp(ctx, "不支持的试用天数")
	}

	ent, err := models.EnterpriseByEmail(db.Mysql, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的客户")
		}

		return internalServerErr(ctx, err.Error())
	}

	updateEnt := func() (err error) {
		tx, err := db.Mysql.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				log.Logger.Warnf("update ent: %s, error: %v", req.Email, err)
				tx.Rollback()
				return
			}

			tx.Commit()
		}()

		d := time.Duration(req.TrialInDay) * time.Hour * 24
		if err = models.SetEntTrial(tx, ent.ID, models.EditionEnterprise, models.TrialIn, models.TrialAgentNum, time.Now().UTC().Add(d)); err != nil {
			return err
		}

		entPlan, err := models.EntPlanByEntID(tx, ent.ID)
		if err != nil {
			return err
		}

		entPlan.ExpirationTime = ent.ExpirationTime.UTC().Add(d)
		entPlan.TrialStatus = models.TrialIn
		entPlan.PlanType = models.EditionEnterprise
		entPlan.AgentServeLimit = models.EditionEnterpriseAgentServeLimit
		return entPlan.Update(tx)
	}

	if err = updateEnt(); err != nil {
		return internalServerErr(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// PUT /super_admin/enterprises/upgrade
func (s *IMService) UpgradeEnterprise(ctx echo.Context) (err error) {
	req := &UpgradeEnterpriseReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Email == "" {
		return invalidParameterResp(ctx, "email 为空")
	}

	ent, err := models.EnterpriseByEmail(db.Mysql, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的客户")
		}

		return internalServerErr(ctx, err.Error())
	}

	entPlan, err := models.EntPlanByEntID(db.Mysql, ent.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的客户")
		}

		return internalServerErr(ctx, err.Error())
	}

	planSet := mapset.NewSet(models.EditionFree, models.EditionStandard, models.EditionEnterprise)
	if !planSet.Contains(req.Plan) {
		return invalidParameterResp(ctx, "不支持的版本")
	}

	if req.AgentNum < 1 {
		return invalidParameterResp(ctx, "坐席数小于1")
	}

	if req.DurationInDay < 30 {
		return invalidParameterResp(ctx, "购买时长太短")
	}

	upgradeFn := func() (err error) {
		tx, err := db.Mysql.Begin()
		if err != nil {
			return err
		}

		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()

		exp := calculateExp(req.DurationInDay, req.AgentNum, ent)
		err = models.UpdateEntInfo(tx, ent.ID, map[string]interface{}{
			"plan":            req.Plan,
			"agent_num":       req.AgentNum,
			"expiration_time": exp,
			"trial_status":    models.TrialNone,
		})
		if err != nil {
			return
		}

		entPlan.ExpirationTime = exp
		entPlan.TrialStatus = models.TrialNone
		entPlan.PlanType = int8(req.Plan)
		entPlan.AgentNum = req.AgentNum
		switch req.Plan {
		case models.EditionStandard:
			entPlan.AgentServeLimit = models.EditionStandardAgentServeLimit
		case models.EditionEnterprise:
			entPlan.AgentServeLimit = models.EditionEnterpriseAgentServeLimit
		default:
			entPlan.AgentServeLimit = models.EditionFreeAgentServeLimit
		}
		return entPlan.Update(tx)
	}

	if err = upgradeFn(); err != nil {
		return internalServerErr(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// DELETE /super_admin/enterprises/tests/:ent_id
func (s *IMService) DeleteTestEnterprise(ctx echo.Context) (err error) {
	entID := ctx.Param("ent_id")
	if err = s.deleteTestEnterprise(entID); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func calculateExp(days int, agentNum int, ent *models.Enterprise) (exp time.Time) {
	var leftHours time.Duration
	hours := time.Duration(days*agentNum) * 24
	now := time.Now().UTC()

	if ent.TrialStatus == models.TrialNone && (ent.Plan == models.EditionStandard || ent.Plan == models.EditionEnterprise) {
		leftHours = time.Duration(ent.ExpirationTime.UTC().Sub(now).Hours())
	}

	if leftHours > 0 {
		avgHours := int64(hours+leftHours) / int64(agentNum)
		exp = now.Add(time.Duration(avgHours) * time.Hour)
		return
	}

	hs := hours / time.Duration(agentNum)
	exp = now.Add(hs * time.Hour)
	return
}

func (s *IMService) deleteTestEnterprise(entID string) (err error) {
	tx, err := db.Mysql.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	deleteSQL := `DELETE FROM custmchat.enterprise WHERE id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.agent WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.agent_group WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.agent_invitation WHERE enterprise_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.ent_all_configs WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.ent_app WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.ent_plan WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.role WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	deleteSQL = `DELETE FROM custmchat.user_group WHERE ent_id=?`
	if _, err = tx.Exec(deleteSQL, entID); err != nil {
		return
	}

	return
}
