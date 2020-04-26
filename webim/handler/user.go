package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/tomasen/realip"
	"golang.org/x/crypto/bcrypt"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var defaultAgentOnlineExpire = 1 * time.Minute

type signInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInResp struct {
	EntID     string `json:"ent_id"`
	UserID    string `json:"user_id"`
	UserToken string `json:"user_token"`
}

type logoutReq struct {
	AllDevices bool `json:"all_devices"`
}

type AddAgentReq struct {
	Email        string   `json:"email"`
	RealName     string   `json:"real_name"`
	NickName     string   `json:"nick_name"`
	InitPassWord string   `json:"init_password"`
	JobNum       string   `json:"job_num"`
	PermsRange   string   `json:"perms_range"`
	PermsGroups  []string `json:"perms_groups"`
	ServeLimit   int      `json:"serve_limit"`
	GroupID      string   `json:"group_id"`
	RoleID       string   `json:"role_id"`
}

type ResetPasswordReq struct {
	Token          string `json:"token"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type ResetForgetPassword struct {
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type onlineAgentsInfoReq struct {
	WithConvNum int    `query:"with_conv_num"`
	BrowserID   string `query:"browser_id"`
}

type onlienAgent struct {
	ActiveConvNum       int    `json:"active_conv_num"`
	HasAccessPermission bool   `json:"has_access_permission"`
	ID                  string `json:"id"`
	Status              string `json:"status"`
}

type onlineAgentsInfoResp struct {
	Agents []*onlienAgent `json:"agents"`
}

type QueueVisitorsResp struct {
	Queue []*adapter.VisitInfo `json:"queue"`
}

func (req *AddAgentReq) validate() error {
	if req.Email == "" || req.RealName == "" || req.JobNum == "" || req.GroupID == "" || req.RoleID == "" {
		return fmt.Errorf("invalid")
	}

	if req.ServeLimit < 0 {
		return fmt.Errorf("serve_limit invalid(can't be negetive)")
	}

	if req.PermsRange != models.AgentPermsRangeAllType && req.PermsRange != models.AgentPermsRangePartType && req.PermsRange != models.AgentPermsRangePersonalType {
		return fmt.Errorf("unsupported perms type")
	}

	if req.PermsRange == models.AgentPermsRangePartType && len(req.PermsGroups) == 0 {
		return fmt.Errorf("groups not selected")
	}

	return nil
}

func ClearExpiredAgentList(key string) error {
	var expire time.Duration

	now := time.Now()
	d := conf.IMConf.AgentConf.AgentOnlineExpire.Duration
	if d > 0 {
		expire = d
	} else {
		expire = defaultAgentOnlineExpire
	}

	maxTime := now.Add(-1 * (expire + time.Second))
	max := strconv.FormatInt(maxTime.Unix(), 10)
	return db.RedisClient.ZRemRangeByScore(key, "0", max).Err()
}

// SignIn agent Sign In
// POST /signin
func (s *IMService) SignIn(ctx echo.Context) (err error) {
	u := new(signInReq)
	err = ctx.Bind(u)
	if err != nil {
		return
	}

	if u.Email == "" || u.Password == "" {
		return invalidParameterResp(ctx, "邮箱或用户名为空")
	}

	user, err := models.AgentByEmail(db.Mysql, u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "用户不存在")
		}

		return dbErrResp(ctx, err.Error())
	}

	if user == nil {
		return invalidParameterResp(ctx, "用户不存在")
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(u.Password)); err != nil {
		return invalidParameterResp(ctx, "密码错误")
	}

	s.RecordLogin(ctx.Request(), user.ID, user.EntID)
	token, err := s.genJWT(user.EntID, user.ID)
	if err != nil {
		return errResp(ctx, common.GenTokenErr, err.Error())
	}

	if err = s.auth.RegisterLogin(user.ID, token); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &signInResp{EntID: user.EntID, UserID: user.ID, UserToken: token}
	if err = s.RefreshAgentOnline(token, user.EntID); err != nil {
		log.Logger.Warnf("RefreshAgentOnline error: %v", err)
	}

	return jsonResponse(ctx, resp)
}

func (s *IMService) CheckAgentNumAndClearExpiredAgent(entID string) (code int, msg string) {
	onlineCount, err := s.getEntOnlineAgentCount(entID)
	if err != nil {
		return common.DBErr, err.Error()
	}

	agentNum, err := models.GetEntAgentNum(db.Mysql, entID)
	if err != nil {
		return common.DBErr, err.Error()
	}

	if err = ClearExpiredAgentList(fmt.Sprintf(common.EntOnlineAgentList, entID)); err != nil {
		log.Logger.Warnf("ClearExpiredAgentUUIDList error: %v", err)
	}

	if onlineCount >= int64(agentNum) {
		return common.UserLoginCountExceedErr, "online agents number exceed"
	}

	return 0, ""
}

func (s *IMService) RefreshAgentOnline(userToken, entID string) error {
	key := fmt.Sprintf(common.EntOnlineAgentList, entID)
	z := redis.Z{Score: float64(time.Now().Unix()), Member: userToken}
	if err := db.RedisClient.ZAdd(key, z).Err(); err != nil {
		log.Logger.Warnf("set ent online agent list error: %v", err)
		return err
	}

	return nil
}

func (s *IMService) RecordLogin(req *http.Request, agentID, entID string) {
	record := &models.LoginRecord{
		ID:          common.GenUniqueID(),
		AgentID:     agentID,
		EntID:       entID,
		LoginAt:     time.Now().UTC(),
		LoginClient: "web",
		LoginIP:     realip.RealIP(req),
		DeviceInfo:  req.UserAgent(),
	}

	if err := record.Insert(db.Mysql); err != nil {
		log.Logger.Error("record agent login error: ", err)
	}
}

func (s *IMService) getEntOnlineAgentCount(entID string) (count int64, err error) {
	var expire time.Duration

	key := fmt.Sprintf(common.EntOnlineAgentList, entID)
	now := time.Now()
	d := conf.IMConf.AgentConf.AgentOnlineExpire.Duration
	if d > 0 {
		expire = d
	} else {
		expire = defaultAgentOnlineExpire
	}

	min := strconv.FormatInt(now.Add(-1*expire).Unix(), 10)
	max := strconv.FormatInt(now.Unix(), 10)
	count, err = db.RedisClient.ZCount(key, min, max).Result()
	return
}

func (s *IMService) genJWT(entID, userID string) (string, error) {
	token, err := newJwtToken(&auth.Claims{
		EntID:  entID,
		UserID: userID,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 10 * time.Hour).Unix(),
			Issuer:    "CustmIM",
		},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// AddAgent add agent to specific ent
// POST /admin/api/v1/agents
func (s *IMService) AddAgent(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentLimit, err := models.GetAgentLimitByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentNum, err := models.AgentNumOfEnt(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if agentNum >= agentLimit {
		return errResp(ctx, common.AgentNumExceedErr, "agent num exceed")
	}

	req := &AddAgentReq{}
	if err = ctx.Bind(&req); err != nil {
		return
	}
	if err = req.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	serveLimit, err := models.GetAgentServeLimitByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if req.ServeLimit >= serveLimit {
		return errResp(ctx, common.AgentServeLimitExceedErr, "agent server limit exceed")
	}

	agent, err := models.AgentByEmail(db.Mysql, req.Email)
	if err == nil && agent != nil {
		return errResp(ctx, common.AgentAlreadyExistsErr, "agent already exists")
	}

	password, err := common.GenHashedPassword([]byte(req.InitPassWord))
	if err != nil {
		log.Logger.Error("gen hashed password err: ", err)
		return internalServerErr(ctx, "add agent error")
	}

	now := time.Now().UTC()
	agentModel := &models.Agent{
		ID:             common.GenUniqueID(),
		EntID:          entID,
		GroupID:        req.GroupID,
		RoleID:         req.RoleID,
		Avatar:         "",
		Username:       req.Email,
		RealName:       req.RealName,
		NickName:       req.NickName,
		HashedPassword: string(password),
		JobNumber:      req.JobNum,
		ServeLimit:     req.ServeLimit,
		Email:          req.Email,
		Status:         models.AgentUnavailableStatus,
		IsAdmin:        0,
		PermsRangeType: req.PermsRange,
		AccountStatus:  models.AgentAccountCreatedStatus,
		CreateAt:       now,
		UpdateAt:       now,
	}

	if err = agentModel.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if agentModel.PermsRangeType == models.AgentPermsRangePartType {
		if err = models.BulkAddAgentPermsRangeGroups(db.Mysql, agentModel.ID, req.PermsGroups); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	agentModel.HashedPassword = ""
	go sendAgentActivateEmail(db.RedisClient, s.mailClient, agentModel.Email, agentModel.ID)

	return jsonResponse(ctx, &Resp{Code: 0, Body: modelAgentToViewAgent(agentModel)})
}

// UpdateAgentInfo
// PUT /admin/api/v1/agents/:agent_id
func (s *IMService) UpdateAgentInfo(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Param("agent_id")
	if agentID == "" {
		return invalidParameterResp(ctx, "invalid agent_id")
	}

	req := &AddAgentReq{}
	if err = ctx.Bind(&req); err != nil {
		return
	}
	if err = req.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	serveLimit, err := models.GetAgentServeLimitByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if req.ServeLimit >= serveLimit {
		return errResp(ctx, common.AgentServeLimitExceedErr, "agent server limit exceed")
	}

	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errResp(ctx, common.UserNotExistErr, "agent not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	if agent.DeletedAt.Valid {
		return errResp(ctx, common.UserNotExistErr, "agent not exists")
	}

	now := time.Now().UTC()
	agent.RealName = req.RealName
	agent.NickName = req.NickName
	agent.JobNumber = req.JobNum
	agent.ServeLimit = req.ServeLimit
	agent.CreateAt = agent.CreateAt.UTC()
	agent.UpdateAt = now

	if req.GroupID != "" {
		agent.GroupID = req.GroupID
	}

	if req.RoleID != "" {
		agent.RoleID = req.RoleID
	}

	if err = agent.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}
	agent.HashedPassword = ""

	return jsonResponse(ctx, &Resp{Code: 0, Body: modelAgentToViewAgent(agent)})
}

func sendAgentActivateEmail(redisClient *redis.Client, emailClient EMailSender, address, agentID string) {
	randomStr := common.RandStringBytesMask(32)
	duration := conf.IMConf.AgentConf.ActivateCodeEffectiveDuration.Duration
	_, err := redisClient.Set(randomStr, agentID, duration).Result()
	if err != nil {
		log.Logger.Error("set activate code key in redis error: ", err)
		return
	}

	title := `激活账户`
	content := fmt.Sprintf(`%s/api/v1/activate?activate_code=%s`, conf.IMConf.Host, randomStr)
	if err := emailClient.SendEmail(address, title, content); err != nil {
		log.Logger.Error("send email error: ", err, "address: ", address)
	}
}

// GET /api/v1/activate?activate_code=xxxxxxxxx
func (s *IMService) ActivateUser(ctx echo.Context) (err error) {
	activateCode := ctx.QueryParam("activate_code")
	if activateCode == "" {
		return invalidParameterResp(ctx, "activate_code is invalid")
	}

	v, err := db.RedisClient.Get(activateCode).Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "invalid activate_code")
		}

		log.Logger.Error("get user id from redis err: ", err)
		return errResp(ctx, common.RedisErr, err.Error())
	}

	if err = models.UpdateAgentAccountStatusValid(db.Mysql, v); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	db.RedisClient.Del(activateCode)
	return jsonResponse(ctx, &Resp{Code: 0})
}

// GenConnectionToken
// GET /admin/api/v1/agents/connection_token
func (s *IMService) GenConnectionToken(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if agentID == "" {
		return invalidParameterResp(ctx, "agent_id invalid")
	}

	token, err := newConnJwtToken(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		Subject:   agentID,
	})

	type resp struct {
		Token string `json:"token"`
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: &resp{Token: token}})
}

func (s *IMService) signOut(entID, userID, token string) (err error) {
	pipe := db.RedisClient.Pipeline()
	defer func() {
		if err := pipe.Close(); err != nil {
			log.Logger.Warnf("close redis pipe error: %v", err)
		}
	}()

	loginCountKey := fmt.Sprintf(common.AgentLoginCount, userID)
	pipe.Decr(loginCountKey)

	tkListKey := fmt.Sprintf(common.AgentTokenList, userID)
	pipe.SRem(tkListKey, token)

	pipe.ZRem(fmt.Sprintf(common.EntOnlineAgentList, entID), token)

	pipe.Del(token)

	if _, err := pipe.Exec(); err != nil {
		log.Logger.Warnf("SignOut redis pipe execute error: %v", err)
		return err
	}

	return nil
}

// SignOutAll 退出全部登录设备
// POST /api/logout
func (s *IMService) Logout(ctx echo.Context) (err error) {
	agentInfo := getAgentInfoFromJwtToken(ctx)
	req := &logoutReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	sendMsg := func(tokens []string) {
		go func() {
			s.sendAgentStatusUpdateEvent(agentInfo.UserID, models.AgentUnavailableStatus, false)
			time.Sleep(200 * time.Millisecond)
			s.sendAgentKickedEvent(agentInfo.UserID, tokens)
		}()
	}

	tkListKey := fmt.Sprintf(common.AgentTokenList, agentInfo.UserID)
	tokens, err := db.RedisClient.SMembers(tkListKey).Result()
	if err != nil {
		log.Logger.Warnf("Get Agent Token List(%s), error: %v", tkListKey, err)
	}

	if req.AllDevices {
		err = auth.SignOutAll(db.RedisClient, agentInfo.Token, agentInfo.EntID, agentInfo.UserID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		sendMsg(tokens)
		return jsonResponse(ctx, &SuccessResp{Success: true})
	}

	if err := s.signOut(agentInfo.EntID, agentInfo.UserID, agentInfo.Token); err != nil {
		return internalServerErr(ctx, err.Error())
	}
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// ForgetPassword
// POST /forget_password
func (s *IMService) ForgetPassword(ctx echo.Context) (err error) {
	reqData := &struct {
		Email string `json:"email"`
	}{}

	if err = ctx.Bind(reqData); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	agent, err := models.AgentByEmail(db.Mysql, reqData.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return errResp(ctx, common.UserNotExistErr, "account not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	if agent.DeletedAt.Valid {
		return errResp(ctx, common.UserNotExistErr, "account not exists")
	}

	go sendForgetPasswordEmail(db.RedisClient, s.mailClient, agent.Email, agent.ID)

	return jsonResponse(ctx, &Resp{Code: 0})
}

func sendForgetPasswordEmail(redisClient *redis.Client, emailClient EMailSender, address, agentID string) {
	randomStr := common.RandStringBytesMask(32)
	_, err := redisClient.Set(randomStr, agentID, 6*time.Hour).Result()
	if err != nil {
		log.Logger.Error("set activate code key in redis error: ", err)
		return
	}

	title := `验证码`
	if err := emailClient.SendEmail(address, title, randomStr); err != nil {
		log.Logger.Error("send email error: ", err, "address: ", address)
	}
}

// POST /api/v1/reset_password
func (s *IMService) ResetForgetPassword(ctx echo.Context) (err error) {
	req := &ResetForgetPassword{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Code == "" || req.NewPassword == "" {
		return invalidParameterResp(ctx, "invalid code/new_password")
	}

	if len(req.NewPassword) < 8 {
		return invalidParameterResp(ctx, "password too short")
	}

	v, err := db.RedisClient.Get(req.Code).Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "invalid code")
		}

		log.Logger.Error("ResetForgetPassword: get agent_id from redis err: ", err)
		return errResp(ctx, common.RedisErr, err.Error())
	}

	_, err = models.AgentByID(db.Mysql, v)
	if err != nil {
		if err == sql.ErrNoRows {
			return errResp(ctx, common.UserNotExistErr, "user not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	password, err := common.GenHashedPassword([]byte(req.NewPassword))
	if err != nil {
		log.Logger.Errorf("gen hashed password err: %v", err)
		return internalServerErr(ctx, "reset password error")
	}

	if err = models.UpdateAgentPassword(db.Mysql, v, string(password)); err != nil {
		return dbErrResp(ctx, err.Error())
	}
	db.RedisClient.Del(req.Code)

	return jsonResponse(ctx, &Resp{Code: 0})
}

// ResetPassword
// POST /reset
// {"token":"SHjSCDFLqmhrdYjnkdBbVGRicAXzgeVxUQnpoUhcLBhveGlJcHFuMrzBhCMsNent","password":"abc123@$","repeat_password":"abc123@$"}
func (s *IMService) ResetPassword(ctx echo.Context) (err error) {
	req := &ResetPasswordReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Token == "" || req.Password == "" || req.RepeatPassword == "" {
		return invalidParameterResp(ctx, "邮箱或密码为空")
	}

	if req.Password != req.RepeatPassword {
		return invalidParameterResp(ctx, "密码不一致")
	}

	email, err := db.RedisClient.Get(req.Token).Result()
	if err != nil {
		if err == redis.Nil {
			return invalidParameterResp(ctx, "链接已过期")
		}
		return invalidParameterResp(ctx, err.Error())
	}

	_, err = models.AgentByEmail(db.Mysql, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "用户不存在")
		}

		return dbErrResp(ctx, err.Error())
	}

	password, err := common.GenHashedPassword([]byte(req.RepeatPassword))
	if err != nil {
		log.Logger.Errorf("gen hashed password err: %v", err)
		return internalServerErr(ctx, "reset password error")
	}

	if err = models.UpdateAgentPasswordByEmail(db.Mysql, string(password), email); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// CheckOnlineAgentCount check if the ent agent online count is exceed
// GET /admin/api/v1/check_online_count
func (s *IMService) CheckOnlineAgentCount(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	code, msg := s.CheckAgentNumAndClearExpiredAgent(entID)
	if code != 0 {
		return errResp(ctx, code, msg)
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

// GET /api/agent/online_agents_info
// GET /api/agent/online_agents
func (s *IMService) OnlineAgentsInfo(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)

	agentInfo := &models.AgentInfo{Mysql: db.Mysql}
	users, err := agentInfo.GetEntAgentIDs(entID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var agents []*onlienAgent
	for _, user := range users {
		convNum := s.getAgentConvNum(user.AgentID)
		if convNum < 0 {
			convNum = 0
		}

		agents = append(agents, &onlienAgent{
			ActiveConvNum:       convNum,
			HasAccessPermission: true,
			ID:                  user.AgentID,
			Status:              models.AgentAvailableStatus,
		})
	}

	return jsonResponse(ctx, &onlineAgentsInfoResp{Agents: agents})
}

// GET /api/agent/online
func (s *IMService) SendOnlineEvent(ctx echo.Context) (err error) {
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	go s.sendAgentStatusUpdateEvent(agentID, models.AgentAvailableStatus, true)
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GET /api/agent/queue
func (s *IMService) QueueVisitors(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	qs, err := models.VisitorQueuesByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var trackIDs []string
	for _, q := range qs {
		trackIDs = append(trackIDs, q.TrackID)
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, trackIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, trackIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var res = []*adapter.VisitInfo{}
	for _, q := range qs {
		var visitor *models.Visitor
		for _, v := range visitors {
			if v.TraceID == q.TrackID {
				visitor = v
				break
			}
		}

		var vt *models.Visit
		for _, v := range visits {
			if v.TraceID == q.TrackID {
				vt = v
				break
			}
		}

		if vt != nil && visitor != nil {
			res = append(res, adapter.ConvertModelVisitToVisit(visitor, vt))
		}

	}

	return ctx.JSON(http.StatusOK, &QueueVisitorsResp{Queue: res})
}
