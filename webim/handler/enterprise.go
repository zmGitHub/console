package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/external/submail"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type ent struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
}

type RegisterEnterpriseReq struct {
	ContactTelephone string `json:"contact_telephone"`
	Email            string `json:"email"`
	Fullname         string `json:"fullname"`
	Password         string `json:"password"`
	Source           string `json:"source"`
}

// reg_token=naCAGraeiUZjKYJBPJFePNkoPRthKIhd&email=952244784@qq.com
type VerifyRegisterReq struct {
	RegToken string `query:"reg_token"`
	Email    string `query:"email"`
}

type superAdmin struct {
	*models.Agent
	DeletedAt string `json:"deleted_at,omitempty"`
}

type entInfo struct {
	EntDetail      *models.Enterprise     `json:"ent_detail"`
	AllocationRule *models.AllocationRule `json:"allocation_rule"`
}

type GetEntInfoResp struct {
	EntInfo    *entInfo    `json:"ent_info"`
	SuperAdmin *superAdmin `json:"super_admin"`
}

type UpdateEntInfoReq struct {
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Industry    string `json:"industry"`
	Mobile      string `json:"mobile"`
	Description string `json:"description"`
}

type RegisterEnterpriseResp struct {
	Token    string `json:"token"`
	AgentID  string `json:"agent_id"`
	EntID    string `json:"ent_id"`
	Fullname string `json:"fullname"`
}

func (s *IMService) RegisterEnterprise(ctx echo.Context) (err error) {
	//if errMsg := verifyCaptcha(ctx, "Captcha-Token", "Captcha-Value"); errMsg != nil {
	//	return ctx.JSON(http.StatusBadRequest, errMsg)
	//}

	e := new(RegisterEnterpriseReq)
	if err = ctx.Bind(e); err != nil {
		return
	}

	if e.Fullname == "" || e.Email == "" || e.ContactTelephone == "" || e.Password == "" {
		return invalidParameterResp(ctx, "企业名称或登录邮箱或手机号或密码为空")
	}

	entModel, err := models.GetEntByEmailORName(db.Mysql, e.Email, e.Fullname)
	if err != nil && err != sql.ErrNoRows {
		return dbErrResp(ctx, err.Error())
	}

	if entModel != nil {
		if e.Fullname == entModel.Name || e.Email == entModel.Email {
			return invalidParameterResp(ctx, "租户名/邮箱 重复")
		}

		return invalidParameterResp(ctx, "租户已存在")
	}

	agent, err := models.AgentByEmail(db.Mysql, e.Email)
	if err != nil && err != sql.ErrNoRows {
		return dbErrResp(ctx, err.Error())
	}

	if agent != nil {
		return invalidParameterResp(ctx, "邮箱已被使用")
	}

	password, err := common.GenHashedPassword([]byte(e.Password))
	if err != nil {
		log.Logger.Warnf("GenHashedPassword error: %v\n", err)
		return internalServerErr(ctx, "internal server error")
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		dbErr = rollBackOrCommit(tx, dbErr)
		if dbErr != nil {
			log.Logger.Errorf("rollBackOrCommit error: %v\n", dbErr)
		}
	}()

	now := time.Now().UTC()
	entID := common.GenUniqueID()
	adminID := common.GenUniqueID()
	ent := &models.Enterprise{
		ID:              entID,
		Name:            e.Fullname,
		FullName:        e.Fullname,
		AdminID:         adminID,
		AllocationRule:  models.OrderTakeTurnsAllocation,
		Email:           e.Email,
		Mobile:          e.ContactTelephone,
		CreatedAt:       now,
		Owner:           e.Email,
		Plan:            models.EditionFree,
		AgentNum:        1,
		TrialStatus:     models.TrialNone,
		ExpirationTime:  now.Add(15 * 24 * time.Hour),
		IsActivated:     false,
		LastActivatedAt: now,
	}

	if dbErr = ent.Insert(tx); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	errMsg := s.initSuperAdmin(tx, entID, adminID, e.Email, string(password))
	if errMsg != nil {
		dbErr = fmt.Errorf(errMsg.Message)
		return jsonResponse(ctx, errMsg)
	}

	if errMsg = s.initPlan(tx, entID, ent.ExpirationTime); errMsg != nil {
		dbErr = fmt.Errorf(errMsg.Message)
		return jsonResponse(ctx, errMsg)
	}

	if dbErr = s.InitEntConfigs(tx, entID); dbErr != nil {
		log.Logger.Warnf("InitEntConfigs: %v", dbErr)
		return internalServerErr(ctx, "初始化配置失败")
	}

	if _, formErr := s.addForm(entID, "询前表单", tx, dto.DefaultChatForms.FormDef); formErr != nil {
		log.Logger.Warnf("addForm error: %v", formErr)
	}

	req := ctx.Request()
	source := req.Referer()
	if err := sendActiveEntEmail(db.RedisClient, ent.Email, entID, adminID, source); err != nil {
		dbErr = err
		return ctx.JSON(http.StatusInternalServerError, &ErrMsg{Message: "internal server error"})
	}

	token, err := s.genJWT(entID, adminID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrMsg{Message: err.Error()})
	}
	if err = s.auth.RegisterLogin(adminID, token); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrMsg{Message: err.Error()})
	}

	monitor.EnterprisesRegisterCount.WithLabelValues("register_direct").Inc()
	return jsonResponse(ctx, &RegisterEnterpriseResp{Token: token, AgentID: adminID, EntID: entID, Fullname: ent.FullName})
}

func (s *IMService) initSuperAdmin(db models.XODB, entID, adminID, email, password string) *ErrMsg {
	now := time.Now().UTC()
	group := &models.AgentGroup{
		ID:          common.GenUniqueID(),
		EntID:       entID,
		Name:        "超级管理员",
		Description: "超管",
	}

	if err := group.Insert(db); err != nil {
		return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	role := &models.Role{
		ID:    common.GenUniqueID(),
		EntID: entID,
		Name:  "超管",
	}
	if err := role.Insert(db); err != nil {
		return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	user := &models.Agent{
		ID:             adminID,
		EntID:          entID,
		GroupID:        group.ID,
		RoleID:         role.ID,
		Username:       email,
		RealName:       "超级管理员",
		ServeLimit:     2,
		HashedPassword: password,
		Email:          email,
		Wechat:         "",
		IsAdmin:        1,
		PermsRangeType: models.AgentPermsRangeAllType,
		Status:         models.AgentUnavailableStatus,
		AccountStatus:  models.AgentAccountCreatedStatus,
		CreateAt:       now,
		UpdateAt:       now,
	}

	if dbErr := user.Insert(db); dbErr != nil {
		mysqlErr, ok := dbErr.(*mysql.MySQLError)
		if !ok {
			return &ErrMsg{Code: common.DBErr, Message: dbErr.Error()}
		}

		if mysqlErr.Number == 1062 {
			return &ErrMsg{Code: common.UserExistsErr, Message: "admin email or mobile exists"}
		}

		return &ErrMsg{Code: common.DBErr, Message: dbErr.Error()}
	}

	if err := s.initPerms(db, entID, role.ID); err != nil {
		return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	return nil

}
func (s *IMService) initPerms(db models.XODB, entID, roleID string) error {
	var apps []string
	for app := range perms {
		apps = append(apps, app)
	}
	if err := models.BulkCreateEntApps(db, entID, apps); err != nil {
		return err
	}

	permIDs, err := models.BulkCreateEntPerms(db, entID, perms)
	if err != nil {
		return err
	}

	return models.UpdateRolePerms(db, roleID, permIDs)
}

func (s *IMService) initPlan(db models.XODB, entID string, expireTime time.Time) *ErrMsg {
	now := time.Now().UTC()
	plan := &models.EntPlan{
		ID:              common.GenUniqueID(),
		EntID:           entID,
		PlanType:        models.EditionFree,
		TrialStatus:     models.TrialNone,
		AgentServeLimit: models.EditionFreeAgentServeLimit,
		LoginAgentLimit: 1,
		AgentNum:        1,
		PayAmount:       0,
		ExpirationTime:  expireTime,
		CreateAt:        now,
		UpdateAt:        now,
	}
	if err := plan.Insert(db); err != nil {
		return &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	return nil
}

func (s *IMService) InitEntConfigs(db models.XODB, entID string) error {
	now := time.Now().UTC()
	configs := &models.EntAllConfig{
		ID:            common.GenUniqueID(),
		EntID:         entID,
		ConfigContent: sql.NullString{String: string(conf.IMConf.DefaultEntConfigs), Valid: true},
		CreateAt:      now,
		UpdateAt:      now,
	}

	if err := configs.Insert(db); err != nil {
		return err
	}

	return nil
}

// GET /register_redirect
func (s *IMService) VerifyRegister(ctx echo.Context) error {
	req := &VerifyRegisterReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.RegToken == "" || req.Email == "" {
		return invalidParameterResp(ctx, "reg_token/email invalid")
	}

	v, err := db.RedisClient.HMGet(req.RegToken, "ent_id", "agent_id", "source").Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "链接过期或不存在")
		}

		log.Logger.Warnf("[VerifyRegister] redis get %s, error: %v", req.RegToken, err)
		return invalidParameterResp(ctx, "验证失败")
	}

	if len(v) < 2 {
		return invalidParameterResp(ctx, "invalid reg_token")
	}

	if v[0] == nil || v[1] == nil {
		return invalidParameterResp(ctx, "invalid reg_token")
	}

	entID, agentID := v[0].(string), v[1].(string)
	updateStatusFn := func() error {
		tx, err := db.Mysql.Begin()
		if err != nil {
			return err
		}

		var dbErr error
		defer func() {
			dbErr = rollBackOrCommit(tx, dbErr)
			if dbErr != nil {
				log.Logger.Warnf("[VerifyRegister] rollBackOrCommit error: %v", dbErr)
			}
		}()

		if dbErr = models.UpdateEntActivated(tx, entID); dbErr != nil {
			return dbErr
		}

		if dbErr = models.UpdateAgentAccountStatusValid(tx, agentID); dbErr != nil {
			return dbErr
		}

		if dbErr = models.UpdateEntTrialStatus(tx, entID, models.TrialNone); dbErr != nil {
			return dbErr
		}

		return nil
	}

	if err = updateStatusFn(); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	configs, errMsg := s.getEnterpriseConfigs(entID)
	if errMsg == nil {
		configs.IsActivated = true
		s.updateEntConfigs(entID, configs)
	}
	db.RedisClient.Del(req.RegToken)
	//host := common.GetHostFromURL(source)
	//if host == "" {
	//	host = conf.IMConf.Host
	//}
	return ctx.Redirect(http.StatusTemporaryRedirect, conf.IMConf.Host+"/signin")
}

func sendActiveEntEmail(redisClient *redis.Client, toAddress, entID, agentID, source string) error {
	randomStr := common.RandStringBytesMask(32)
	duration := conf.IMConf.AgentConf.ActivateCodeEffectiveDuration.Duration
	_, err := redisClient.HMSet(randomStr, map[string]interface{}{
		"ent_id":   entID,
		"agent_id": agentID,
		"source":   source,
	}).Result()
	if err != nil {
		log.Logger.Warnf("HMSet activate code key(%s) in redis error: %v", randomStr, err)
		return err
	}

	if err := redisClient.Expire(randomStr, duration).Err(); err != nil {
		log.Logger.Error("Expire activate code key in redis error: ", err)
		return err
	}

	title := `激活账号`
	registerURL := fmt.Sprintf(conf.IMConf.BackendHost+"/register_redirect?reg_token=%s&email=%s", randomStr, toAddress)
	log.Logger.Infof("registerURL: %s", registerURL)

	err = submail.SendActivateEnterprise(toAddress, title, "点击下面链接激活账号", registerURL)
	if err != nil {
		log.Logger.Warnf("send to: %s, email error: %v", toAddress, err)
		return err
	}
	return nil
}

// ActivateEnterprise
// GET /api/v1/activate_ent?activate_code=xxxxxxxxx
func (s *IMService) ActivateEnterprise(ctx echo.Context) (err error) {
	activateCode := ctx.QueryParam("activate_code")
	if activateCode == "" {
		return invalidParameterResp(ctx, "activate_code is invalid")
	}

	v, err := db.RedisClient.HMGet(activateCode, "ent_id", "agent_id").Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "invalid activate_code")
		}

		log.Logger.Error("get user id from redis err: ", err)
		return dbErrResp(ctx, err.Error())
	}

	if len(v) < 2 {
		return invalidParameterResp(ctx, "invalid activate_code")
	}

	if v[0] == nil || v[1] == nil {
		return invalidParameterResp(ctx, "invalid activate_code")
	}

	entID, agentID := v[0].(string), v[1].(string)
	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		dbErr = rollBackOrCommit(tx, dbErr)
		if dbErr != nil {
			log.Logger.Errorf("ActivateEnterprise rollBackOrCommit error: %v", dbErr)
		}
	}()

	if dbErr = models.UpdateAgentAccountStatusValid(tx, agentID); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	if dbErr = models.UpdateEntTrialStatus(tx, entID, models.TrialIn); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	db.RedisClient.Del(activateCode)
	return jsonResponse(ctx, &Resp{Code: 0})
}

// PUT /admin/api/v1/enterprises/info
// PUT /api/enterprise
func (s *IMService) UpdateEntInfo(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if errMsg := hasPerm(entID, agentID, "ent_info", "config_ent_team"); errMsg != nil {
			return ctx.JSON(http.StatusForbidden, errMsg)
		}
	}

	req := &adapter.EnterpriseResp{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	var values = map[string]interface{}{}

	if req.FullName != "" {
		values["full_name"] = req.FullName
	}

	if req.Province != "" {
		values["province"] = req.Province
	}

	if req.City != "" {
		values["city"] = req.City
	}

	if req.Name != "" {
		values["name"] = req.Name
	}

	if req.PublicNickname != "" {
		values["nick_name"] = req.PublicNickname
	}

	if req.PublicCellphone != "" {
		values["mobile"] = req.PublicCellphone
	}

	if req.PublicQQ != "" {
		values["contact_qq"] = req.PublicQQ
	}

	// 座机
	if req.PublicTelephone != "" {
		values["mobile"] = req.PublicTelephone
	}

	if req.PublicWeixin != "" {
		values["contact_wechat"] = req.PublicWeixin
	}

	if req.PublicSignature != "" {
		values["contact_signature"] = req.PublicSignature
	}

	if req.PublicEmail != "" {
		values["email"] = req.PublicEmail
	}

	if req.Avatar != "" {
		values["avatar"] = req.Avatar
	}

	if req.Industry != "" {
		values["industry"] = req.Industry
	}

	if req.ContactEmail != "" {
		values["contact_email"] = req.ContactEmail
	}

	if req.ContactName != "" {
		values["contact_name"] = req.ContactName
	}

	if req.ContactTelephone != "" {
		values["contact_mobile"] = req.ContactTelephone
	}

	if req.MailingAddress != "" {
		values["address"] = req.MailingAddress
	}

	var newAdminID string
	if req.SuperadminID != nil {
		v, err := models.IsAdmin(db.Mysql, agentID)
		if err != nil {
			if err == sql.ErrNoRows {
				return invalidParameterResp(ctx, "坐席不存在")
			}

			return dbErrResp(ctx, err.Error())
		}

		if !v {
			return ctx.JSON(http.StatusForbidden, &ErrMsg{Message: "只有管理员才可以更改super admin"})
		}

		if agentID != *req.SuperadminID {
			values["admin_id"] = *req.SuperadminID
			newAdminID = *req.SuperadminID
		}
	}

	changeAdmin := func() error {
		tx, err := db.Mysql.Begin()
		if err != nil {
			return err
		}

		if err = models.UpdateEntInfo(tx, entID, values); err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				tx.Rollback()
				return fmt.Errorf("contact_name/public_email 已存在")
				// return invalidParameterResp(ctx, )
			}
			tx.Rollback()
			return err
		}

		if newAdminID == "" {
			return tx.Commit()
		}

		if err := models.ChangeAdmin(tx, agentID, newAdminID); err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	}

	if err := changeAdmin(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	ent, err := models.EnterpriseByID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errResp(ctx, common.EntNotExistErr, "企业不存在")
		}

		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, adapter.ConvertEntToAdapterEnt(ent))
}
