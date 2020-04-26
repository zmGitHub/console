package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/labstack/echo/v4"
	"github.com/mssola/user_agent"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/handler/monitor"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type CallerInfo struct {
	Caller     string `json:"caller,omitempty"`
	OperatorId int64  `json:"operator_id,omitempty"`
	SourceIp   string `json:"source_ip,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
}

type UAInfo struct {
	BrowserFamily        string `json:"browser_family,omitempty"`
	BrowserVersionString string `json:"browser_version_string,omitempty"`
	BrowserVersion       string `json:"browser_version,omitempty"`
	OsCategory           string `json:"os_category,omitempty"`
	OsFamily             string `json:"os_family,omitempty"`
	OsVersionString      string `json:"os_version_string,omitempty"`
	OsVersion            string `json:"os_version,omitempty"`
	Platform             string `json:"platform,omitempty"`
	UaString             string `json:"ua_string,omitempty"`
	DeviceFamily         string `json:"device_family,omitempty"`
}

type LocationInfo struct {
	Ip       string `json:"ip,omitempty"`
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
	Isp      string `json:"isp,omitempty"`
}

type InitVisitReq struct {
	EntID       string `json:"ent_id"`
	TraceID     string `json:"trace_id"`
	Keyword     string `json:"keyword"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	ReferrerURL string `json:"referrer_url"`
}

type InitVisitResp struct {
	VisitInfo *models.Visit `json:"visit_info"`
	VisitID   string        `json:"visit_id"`
	VisitorID string        `json:"visitor_id"`
	TraceID   string        `json:"trace_id"`
	Visitor   *models.Visitor
}

type InitVisitV1Req struct {
	EntID       string `query:"ent_id"`
	TrackID     string `query:"track_id"`
	Title       string `query:"title"`
	URL         string `query:"url"`
	ReferrerURL string `query:"referrer_url"`
	JsonpCb     string `query:"jsonp_cb"`
	V           int64  `query:"v"`
}

type OnlineVisitsReq struct {
	Province     string `query:"province"`
	SearchEngine string `query:"search_engine"` // direct_access
	PageCount    int    `query:"page_count"`
	EntID        string `query:"ent_id"`
}

type OnlineVisitsResp struct {
	VisitCnt int                  `json:"visit_cnt"`
	Visits   []*adapter.VisitInfo `json:"visits"`
}

// {"black": true, "success": false}
type blackResp struct {
	Black   bool `json:"black"`
	Success bool `json:"success"`
}

// InitVisit
// POST /api/v1/enterprises/:ent_id/visits
func (s *IMService) InitVisit(ctx echo.Context) (err error) {
	visit := new(InitVisitReq)
	if err = ctx.Bind(visit); err != nil {
		return
	}

	visit.URL = ctx.Request().RequestURI
	referer := ctx.Request().Header.Get("Referer")
	if referer != "" {
		visit.ReferrerURL = referer
	} else {
		visit.ReferrerURL = visit.URL
	}

	visit.EntID = ctx.Param("ent_id")
	if visit.EntID == "" || visit.Title == "" || visit.URL == "" {
		return invalidParameterResp(ctx, "ent_id/title/url 为空")
	}

	if visit.ReferrerURL == "" {
		visit.ReferrerURL = visit.URL
	}

	if visit.TraceID == "" {
		resp, err := s.createVisit(ctx, visit)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		return jsonResponse(ctx, &Resp{Code: 0, Body: resp})
	}

	resp, err := s.updateVisit(visit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: resp})
}

func (s *IMService) createVisit(ctx echo.Context, visit *InitVisitReq) (resp *InitVisitResp, err error) {
	httpReq := ctx.Request()
	uaInfo := s.getUAInfoFromReq(httpReq)
	location, err := s.getLocationFromReq(ctx.RealIP())
	if err != nil {
		return
	}

	var platform = "other"
	if uaInfo.Platform == "Windows" || uaInfo.Platform == "Macintosh" || uaInfo.Platform == "Linux" {
		platform = "pc"
	}

	if uaInfo.Platform == "iPhone" || uaInfo.Platform == "Android" {
		platform = "mobile"
	}

	now := time.Now().UTC()
	visitInfo := &models.Visit{
		ID:                   common.GenUniqueID(),
		EntID:                visit.EntID,
		VisitPageCnt:         1,
		ResidenceTimeSec:     1,
		BrowserFamily:        uaInfo.BrowserFamily,
		BrowserVersion:       uaInfo.BrowserVersion,
		BrowserVersionString: uaInfo.BrowserVersionString,
		OsCategory:           uaInfo.OsCategory,
		OsFamily:             uaInfo.OsFamily,
		OsVersion:            uaInfo.OsVersion,
		OsVersionString:      uaInfo.OsVersionString,
		Platform:             platform,
		UaString:             uaInfo.UaString,
		IP:                   location.Ip,
		Country:              location.Country,
		Province:             location.Province,
		City:                 location.City,
		FirstPageTitle:       visit.Title,
		FirstPageSource:      visit.URL,
		LatestTitle:          visit.Title,
		LatestURL:            visit.URL,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	u, err := url.Parse(visit.URL)
	if err != nil {
		return
	}
	visitInfo.FirstPageDomain = u.Host
	visitInfo.FirstPageSourceKeyword = visit.Keyword
	visitInfo.FirstPageSourceURL = visit.ReferrerURL

	u, err = url.Parse(visit.ReferrerURL)
	if err != nil {
		return
	}
	visitInfo.FirstPageSourceDomain = u.Host

	tx, err := db.Mysql.Begin()
	if err != nil {
		return nil, err
	}

	var traceID string
	if visit.TraceID != "" {
		traceID = visit.TraceID
	} else {
		traceID = common.GenUniqueID()
	}
	visitInfo.TraceID = traceID

	var dbErr error
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Debug("recover from  db panic")
		}

		dbErr = rollBackOrCommit(tx, dbErr)
		if dbErr != nil {
			log.Logger.Errorf("rollBackOrCommit error: %v\n", dbErr)
		}
	}()

	if dbErr = visitInfo.Insert(tx); dbErr != nil {
		return nil, dbErr
	}

	visitor := &models.Visitor{
		ID:               common.GenUniqueID(),
		EntID:            visit.EntID,
		TraceID:          traceID,
		Name:             fmt.Sprintf("#visitor:%d", now.UTC().Unix()),
		Avatar:           "",
		VisitCnt:         1,
		VisitPageCnt:     1,
		ResidenceTimeSec: 1,
		LastVisitID:      visitInfo.ID,
		VisitedAt:        visitInfo.CreatedAt,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if dbErr = visitor.Insert(tx); dbErr != nil {
		return nil, dbErr
	}

	resp = &InitVisitResp{
		VisitInfo: visitInfo,
		VisitID:   visitInfo.ID,
		VisitorID: visitor.ID,
		TraceID:   traceID,
		Visitor:   visitor,
	}

	return resp, nil
}

func (s *IMService) updateVisit(visit *InitVisitReq) (resp *InitVisitResp, err error) {
	tx, err := db.Mysql.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("recover from db panic")
		}

		dbErr := rollBackOrCommit(tx, err)
		if dbErr != nil {
			log.Logger.Errorf("rollBackOrCommit error: %v\n", dbErr)
		}
	}()

	visitIDs, err := models.VisitIDsByEntIDTraceID(tx, visit.EntID, visit.TraceID)
	if err != nil {
		return
	}

	if len(visitIDs) == 0 {
		return nil, fmt.Errorf("no visit IDs")
	}

	visitor, err := models.VisitorByEntIDTraceID(tx, visit.EntID, visit.TraceID)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	visitor.LastVisitID = visitIDs[0]
	visitor.VisitedAt = now

	if err = visitor.UpdateVisitorCount(tx); err != nil {
		return
	}

	vt := &models.Visit{
		ID:          visitIDs[0],
		EntID:       visit.EntID,
		TraceID:     visit.TraceID,
		LatestURL:   visit.URL,
		LatestTitle: visit.Title,
		UpdatedAt:   now,
	}

	if err = vt.UpdateCountByID(tx); err != nil {
		return
	}

	resp = &InitVisitResp{VisitID: visitIDs[0], VisitorID: visitor.ID, TraceID: visit.TraceID, Visitor: visitor}
	visitInfo, err := models.VisitByID(tx, resp.VisitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}

		return nil, err
	}
	resp.VisitInfo = visitInfo
	return resp, nil
}

func (s *IMService) getUAInfoFromReq(req *http.Request) *UAInfo {
	uaStr := req.UserAgent()
	ua := user_agent.New(uaStr)
	name, version := ua.Browser()
	osInfo := ua.OSInfo()

	info := &UAInfo{
		BrowserFamily:        name,
		BrowserVersion:       version,
		BrowserVersionString: version,
		OsCategory:           osInfo.Name,
		OsFamily:             osInfo.FullName,
		OsVersionString:      osInfo.Version,
		OsVersion:            osInfo.Version,
		Platform:             ua.Platform(),
		UaString:             uaStr,
		DeviceFamily:         ua.Platform(),
	}

	return info
}

func (s *IMService) getLocationFromReq(ip string) (location *LocationInfo, err error) {
	location = &LocationInfo{Ip: ip}
	location.Country, location.Province, location.City, location.Isp, err = s.loc.GetLocation(ip)
	if err != nil {
		log.Logger.Errorf("loc.GetLocation ip: %s, error: %v\n", ip, err)
	}
	return
}

// GetEntVisits this handler search visit through es, there will be a lot arguments
// GET /admin/api/v1/enterprise/visits
func (s *IMService) GetEntVisits(ctx echo.Context) (err error) {
	return nil
}

// UpdateVisitResidenceTimeSec
// POST  /api/v1/enterprises/:ent_id/visits/:visit_id/update_residence
func (s *IMService) UpdateVisitResidenceTimeSec(ctx echo.Context) (err error) {
	type req struct {
		ResidenceTimeSec int64  `json:"residence_time_sec"`
		VisitorID        string `json:"visitor_id"`
	}

	r := &req{}
	if err = ctx.Bind(r); err != nil {
		return
	}

	if r.ResidenceTimeSec <= 0 || r.VisitorID == "" {
		return invalidParameterResp(ctx, "residence_time_sec/visitor_id is invalid")
	}

	visitID := ctx.Param("visit_id")
	if visitID == "" {
		return invalidParameterResp(ctx, "visit_id is invalid")
	}

	if err = models.IncrVisitResidenceTimeSec(db.Mysql, visitID, r.ResidenceTimeSec); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if err = models.IncrVisitorResidenceTimeSec(db.Mysql, r.VisitorID, r.ResidenceTimeSec); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

// InitVisitV1 ...
// GET /visit/init
func (s *IMService) InitVisitV1(ctx echo.Context) (err error) {
	req := &InitVisitV1Req{}
	if err = ctx.Bind(req); err != nil {
		return ctx.JSONP(http.StatusOK, req.JsonpCb, &ErrMsg{common.InvalidParameterErr, err.Error()})
	}

	if req.EntID == "" || req.Title == "" || req.URL == "" {
		return ctx.JSONP(http.StatusOK, req.JsonpCb, &ErrMsg{common.InvalidParameterErr, "ent_id/title/url is invalid"})
	}

	if req.TrackID != "" && req.EntID != "" {
		black, _ := models.VisitBlacklistByEntIDTraceID(db.Mysql, req.EntID, req.TrackID)
		if black != nil {
			return ctx.JSONP(http.StatusOK, req.JsonpCb, &blackResp{Black: true, Success: false})
		}
	}

	visit := &InitVisitReq{
		EntID:       req.EntID,
		TraceID:     req.TrackID,
		URL:         req.URL,
		Title:       req.Title,
		Keyword:     "",
		ReferrerURL: req.ReferrerURL,
	}

	referer := visit.ReferrerURL
	if referer != "" {
		visit.ReferrerURL = referer
	} else {
		visit.ReferrerURL = visit.URL
	}

	configs, errMsg := s.getConfigsFromCache(visit.EntID)
	if errMsg != nil {
		return ctx.JSONP(http.StatusOK, req.JsonpCb, errMsg)
	}

	ent, err := models.EnterpriseByID(db.Mysql, req.EntID)
	if err != nil {
		return ctx.JSONP(http.StatusBadRequest, req.JsonpCb, &ErrMsg{common.EntNotExistErr, "ent not exists"})
	}

	result := &adapter.InitVisitResp{
		BrowserID:                   common.GenUniqueID(),
		EnterpriseInfo:              adapter.EnterpriseRespToEnterpriseInfo(adapter.ConvertEntToAdapterEnt(ent)),
		Survey:                      configs.Survey,
		InvitationConfig:            configs.InvitationConfig,
		RobotSettings:               configs.RobotSettings,
		QueueingSettings:            configs.QueueSettings,
		TicketConfig:                configs.TicketConfig,
		SendFileSettings:            configs.SendFileSettings,
		ServiceEvaluationConfig:     configs.ServiceEvaluationConfig,
		StandaloneWindowConfig:      configs.StandaloneWindowConfig,
		WidgetSettings:              configs.WidgetSettings,
		VisitorStatusAgentToken:     "",
		Facade:                      &adapter.Facade{},
		InQueue:                     false,
		SearchEngine:                "",
		Success:                     true,
		EntWelcomeMessage:           configs.WelcomeMsgSettings.Web.Content,
		Servability:                 true,
		SchedulerAfterClientSendMsg: false,
		BaiduBidBlackList:           []string{},
		VisitorStatus:               -1,
		TrackID:                     req.TrackID,
	}

	if req.TrackID != "" && models.IsVisitNotExists(db.Mysql, req.EntID, req.TrackID) {
		visit.TraceID = ""
	}

	if visit.TraceID == "" {
		resp, err := s.createVisit(ctx, visit)
		if err != nil {
			log.Logger.Warnf("createVisit error: %v", err)
			return ctx.JSONP(http.StatusOK, req.JsonpCb, &ErrMsg{common.InternalServerErr, "internal server error"})
		}

		s.createVisitPages(resp.VisitInfo)

		result.VisitID = resp.VisitInfo.ID
		result.BrowserFamily = resp.VisitInfo.BrowserFamily
		result.TrackID = resp.TraceID

		updateVisitCount(visit.EntID, 1, 1, 1)

		go s.sendVisitInitEvent(resp.Visitor, resp.VisitInfo)

		monitor.VisitorsComesIn.WithLabelValues("scheduler").Inc()
		return ctx.JSONP(http.StatusOK, req.JsonpCb, result)
	}

	resp, err := s.updateVisit(visit)
	if err != nil {
		log.Logger.Warnf("updateVisit error: %v", err)
		return ctx.JSONP(http.StatusOK, req.JsonpCb, &ErrMsg{common.InternalServerErr, "internal server error"})
	}

	updateVisitCount(visit.EntID, 0, 1, 1)

	if resp.VisitInfo != nil {
		result.VisitID = resp.VisitInfo.ID
		result.BrowserFamily = resp.VisitInfo.BrowserFamily
		s.createVisitPages(resp.VisitInfo)

		go s.sendVisitInitEvent(resp.Visitor, resp.VisitInfo)
	}

	return ctx.JSONP(http.StatusOK, req.JsonpCb, result)
}

func (s *IMService) createVisitPages(visitInfo *models.Visit) {
	now := time.Now().UTC()
	page := models.VisitPage{
		ID:            common.GenUniqueID(),
		EntID:         visitInfo.EntID,
		VisitID:       visitInfo.ID,
		IP:            visitInfo.IP,
		Source:        visitInfo.FirstPageSource,
		SourceKeyword: visitInfo.FirstPageSourceKeyword,
		SourceDomain:  visitInfo.FirstPageDomain,
		SourceURL:     visitInfo.FirstPageSourceURL,
		Title:         visitInfo.FirstPageTitle,
		Domain:        visitInfo.FirstPageDomain,
		URL:           visitInfo.FirstPageURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := page.Insert(db.Mysql); err != nil {
		log.Logger.Warnf("createVisitPages error: %v", err)
	}
}

// GET /api/visit/search?search_engine=&page_count=100&ent_id=38523&browser_id=agent1553935625467&v=1553935625880
func (s *IMService) OnlineVisits(ctx echo.Context) error {
	req := &OnlineVisitsReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.EntID == "" {
		return invalidParameterResp(ctx, "ent_id is empty")
	}

	traceIDs, visitors, err := s.onlineVisitors(req.EntID)
	if err != nil {
		return internalServerErr(ctx, err.Error())
	}

	findVisitor := func(traceID string) *models.Visitor {
		for _, visitor := range visitors {
			if visitor.TraceID == traceID {
				return visitor
			}
		}
		return nil
	}

	visits, err := models.TraceVisitsByConds(db.Mysql, traceIDs, req.Province, req.PageCount)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var results = []*adapter.VisitInfo{}
	for _, visit := range visits {
		visitor := findVisitor(visit.TraceID)
		visitInfo := adapter.ConvertModelVisitToVisit(visitor, visit)
		results = append(results, visitInfo)
	}

	return jsonResponse(ctx, &OnlineVisitsResp{VisitCnt: len(results), Visits: results})
}

// POST visit/:ent_id/:track_id/reject
func (s *IMService) RejectVisit(ctx echo.Context) error {
	entID := ctx.Param("ent_id")
	trackID := ctx.Param("track_id")
	if entID == "" || trackID == "" {
		return invalidParameterResp(ctx, "ent_id/track_id empty")
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, []string{trackID})
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	visit := visits[0]

	go s.sendInviteRejectEvent(entID, visit.ID, trackID, float64(visit.CreatedAt.Unix()))
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// POST /api/visit/:track_id/invite
func (s *IMService) InviteVisitor(ctx echo.Context) error {
	trackID := ctx.Param("track_id")
	if trackID == "" {
		return invalidParameterResp(ctx, "track_id empty")
	}

	q, err := models.VisitorQueueByTrackID(db.Mysql, trackID)
	if err != nil && err != sql.ErrNoRows {
		return invalidParameterResp(ctx, err.Error())
	}

	if q != nil {
		return invalidParameterResp(ctx, "不能邀请在排队的顾客")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, []string{trackID})
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	visit := visits[0]

	traceStart := float64(visit.CreatedAt.Unix())
	event := InviteClientEvent(agentID, entID, agent.NickName, agent.RealName, visit.ID, trackID, traceStart)
	content, err := common.Marshal(event)
	if err == nil {
		sendMessageToAgent(s.imCli, agent.ID, content)
		sendMessageToVisitor(s.imCli, trackID, entID, content)
	}
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GET /api/agent/get_visit_filter
func (s *IMService) GetVisitFilter(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	traceIDs, _, err := s.onlineVisitors(entID)
	if err != nil {
		return internalServerErr(ctx, err.Error())
	}

	visitsInfo, err := models.VisitsByTraceIDs(db.Mysql, traceIDs)

	var osCategory = &dto.OsCategory{}
	var locals = map[string]mapset.Set{}
	for _, visit := range visitsInfo {
		v, ok := locals[visit.Province]
		if ok {
			v.Add(visit.City)
			locals[visit.Province] = v
		} else {
			locals[visit.Province] = mapset.NewSet(visit.City)
		}

		var other = true
		if common.IsPC(visit.OsCategory) {
			other = false
			osCategory.Pc = append(osCategory.Pc, visit.OsCategory)
		}

		if common.IsMobile(visit.OsCategory) {
			other = false
			osCategory.Mobile = append(osCategory.Mobile, visit.OsCategory)
		}

		if other {
			osCategory.Other = append(osCategory.Other, visit.OsCategory)
		}
	}

	osCategory.Pc = mapset.NewSetFromSlice(osCategory.Pc).ToSlice()
	osCategory.Other = mapset.NewSetFromSlice(osCategory.Other).ToSlice()
	osCategory.Mobile = mapset.NewSetFromSlice(osCategory.Mobile).ToSlice()

	result := &dto.VisitFilter{
		OsCategoryHash: osCategory,
		SourceHash: &dto.SourceHash{
			DirectAccess: []interface{}{},
		},
	}
	result.Locals = map[string][]interface{}{}
	for k, v := range locals {
		result.Locals[k] = v.ToSlice()
	}

	return jsonResponse(ctx, result)
}

func (s *IMService) sendVisitInitEvent(visitor *models.Visitor, visit *models.Visit) {
	visitInfo := adapter.ConvertModelVisitToVisit(visitor, visit)
	event := VisitComeEvent(visitInfo)
	eventContent, err := common.Marshal(event)
	if err != nil {
		log.Logger.Warnf("marshal event error: %v", err)
	}

	s.sendEventToAllAgents(visitInfo.EnterpriseID, eventContent)
}

func updateVisitCount(entID string, visitor, visit, visitPage int) {
	var values = []string{
		"visit_count",
		"visit_page_count",
		"visitor_count",
	}

	convStats := &models.ConversationStat{
		EntID:          entID,
		VisitCount:     uint(visit),
		VisitorCount:   uint(visitor),
		VisitPageCount: uint(visitPage),
		CreatedAt:      getDate(),
	}

	if err := convStats.Insert(db.Mysql, values); err != nil {
		log.Logger.Warnf("updateVisitCount error: %v", err)
	}
}
