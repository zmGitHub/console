package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/external/elasticsearch"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type attrs struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Gender  string `json:"gender"`
	Tel     string `json:"tel"`    // tel
	QQ      string `json:"qq"`     // qq
	Weixin  string `json:"weixin"` // weixin
	Weibo   string `json:"weibo"`
	Address string `json:"address"`
	Email   string `json:"email"`   // email
	Comment string `json:"comment"` // comment
}

type updateVisitorReq struct {
	ConvID  string  `json:"conv_id"`
	VisitID string  `json:"visit_id"`
	Attrs   attrsV1 `json:"attrs"`
}

type updateVisitorNameReq struct {
	TrackID string `json:"track_id"`
	ConvID  string `json:"conv_id"`
	Name    string `json:"name"`
}

// {
//    "avatar":"/static/client-avatar/01-03.png",
//    "created_on":1538145329,
//    "enterprise_id":5869,
//    "id":"d93bc6d11f34c40c59d0c9c359dd320a",
//    "name":"Visitor123",
//    "success":true,
//    "track_id":"1AqDeEPTLa5sANH6FigxTyrnQk8"
//}
type updateVisitorNameResp struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	CreatedOn    string `json:"created_on"`
	EnterpriseID string `json:"enterprise_id"`
	TrackID      string `json:"track_id"`
	Success      bool   `json:"success"`
}

type AddTagToVisitorReq struct {
	TagID  string `json:"tag_id"`
	ConvID string `json:"conv_id"`
}

type rule struct {
	Attribute string `json:"attribute"` // attribute
	Condition string `json:"condition"` // condition
	Type      string `json:"type"`      // type
	Value     string `json:"value"`     // value
}
type SearchVisitorsReq struct {
	Rules []*rule `json:"rules"`
	Page  int     `json:"page"`
	Count int     `json:"count"`
}

type SearchVisitorsResp struct {
	ExportLimit int                  `json:"export_limit"`
	Total       int64                `json:"total"`
	Info        []*dto.SearchVisitor `json:"info"`
}

// GET /api/client/:track_id/attrs
func (s *IMService) GetVisitorInfo(ctx echo.Context) (err error) {
	trackID := ctx.Param("track_id")
	if trackID == "" {
		return invalidParameterResp(ctx, "bad request")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	attrs := setVisitorAttrs(visitor)
	customAttrs, err := models.GetVisitorAttrs(db.Mysql, entID, trackID)
	if err == nil {
		for k, v := range customAttrs {
			attrs[k] = v
		}
	}

	return jsonResponse(ctx, attrs)
}

// UpdateVisitorInfo ...
// PUT /admin/api/v1/visitors
// POST /api/client/:track_id/attrs
func (s *IMService) UpdateVisitorInfo(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor", "add_update_visitor_card"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	trackID := ctx.Param("track_id")
	if trackID == "" {
		return invalidParameterResp(ctx, "track_id invalid")
	}

	req := new(updateVisitorReq)
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Attrs == nil {
		return invalidParameterResp(ctx, "attrs invalid")
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not found")
		}

		return dbErrResp(ctx, err.Error())
	}

	// NewAttrsChange
	var attrsValue = map[string]interface{}{}
	if err := s.updateVisitorAttrs(entID, trackID, visitor, req.Attrs, attrsValue); err != nil {
		log.Logger.Warnf("updateVisitorAttrs error: %v", err)
	}

	go s.sendVisitorAttrsUpdate(agentID, trackID, attrsValue)
	return jsonResponse(ctx, dto.ModelVisitorToVisitor(visitor))
}

func (s *IMService) updateVisitorAttrs(entID, trackID string, visitor *models.Visitor, attrs attrsV1, attrsValue map[string]interface{}) (err error) {
	attrs.setAttrs(visitor, attrsValue)
	visitor.UpdatedAt = time.Now().UTC()
	if err = visitor.Update(db.Mysql); err != nil {
		return
	}

	if len(attrs) > 0 {
		var originAttrs map[string]interface{}
		originAttrs, err = models.GetVisitorAttrs(db.Mysql, entID, trackID)
		if err == nil {
			for k, v := range attrs {
				originAttrs[k] = v
			}
			err = models.UpdateVisitorAttrs(db.Mysql, entID, trackID, originAttrs)
		}
	}

	go CreateOrUpdateVisitorESDoc(db.Mysql, entID, trackID, visitor)
	return
}

// GET /api/visit/:visit_id/pages?track_id=xxxxx
func (s *IMService) GetVisitorPages(ctx echo.Context) (err error) {
	trackID := ctx.QueryParam("track_id")
	visitID := ctx.Param("visit_id")
	if trackID == "" || visitID == "" {
		return invalidParameterResp(ctx, "track_id/visit_id is empty")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	pages, err := models.VisitPagesByEntIDVisitID(db.Mysql, entID, visitID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &dto.VisitPages{}
	for _, page := range pages {
		resp.Pages = append(resp.Pages, &dto.VisitPage{
			CreatedOn:     *common.ConvertUTCToTimeString(page.CreatedAt.UTC()),
			Source:        page.Source,
			SourceKeyword: page.SourceKeyword,
			SourceURL:     page.SourceURL,
			Title:         page.Title,
			URL:           page.SourceURL,
		})
	}

	return jsonResponse(ctx, resp)
}

// PUT /api/visit/:visitor_id
func (s *IMService) UpdateVisitorName(ctx echo.Context) (err error) {
	req := &updateVisitorNameReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.ConvID == "" || req.Name == "" || req.TrackID == "" {
		return invalidParameterResp(ctx, "track_id/name/conv_id invalid")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, req.TrackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not found")
		}

		return dbErrResp(ctx, err.Error())
	}

	visitor.Name = req.Name
	visitor.UpdatedAt = time.Now().UTC()
	if err = visitor.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	go elasticsearch.UpdateConversationVisitorName(elasticsearch.ESClient, req.ConvID, visitor.Name)
	go CreateOrUpdateVisitorESDoc(db.Mysql, entID, req.TrackID, visitor)
	return jsonResponse(ctx, &updateVisitorNameResp{
		ID:           visitor.ID,
		Name:         visitor.Name,
		Avatar:       visitor.Avatar,
		CreatedOn:    *common.ConvertUTCToTimeString(visitor.CreatedAt),
		EnterpriseID: visitor.EntID,
		TrackID:      visitor.TraceID,
		Success:      true,
	})
}

// AddTagToVisitor ...
// POST /admin/api/v1/visitors/tags
// POST /api/agent/client/:track_id/tags
// {"tag_id":10907,"conv_id":502966874,"browser_id":"agent1552606025544"}
func (s *IMService) AddTagToVisitor(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor_and_conv", "add_del_client_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &AddTagToVisitorReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.TagID == "" || req.ConvID == "" {
		return invalidParameterResp(ctx, "tag_id/conv_id invalid")
	}
	trackID := ctx.Param("track_id")
	if trackID == "" {
		return invalidParameterResp(ctx, "invalid track_id")
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	r := &models.VisitorTagRelation{
		VisitorID: visitor.ID,
		TagID:     req.TagID,
	}
	if err = r.Insert(db.Mysql); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return invalidParameterResp(ctx, "tag 已存在")
		}

		return dbErrResp(ctx, err.Error())
	}

	if err = models.IncrVisitorTagUseCount(db.Mysql, r.TagID); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	tags, err := models.VisitorTagRelationsByVisitorIDs(db.Mysql, []string{visitor.ID})
	if err == nil {
		var ts []string
		for _, tag := range tags {
			ts = append(ts, tag.TagID)
		}
		elasticsearch.UpdateConversationTags(elasticsearch.ESClient, req.ConvID, ts)
	}

	go s.sendUseClientTagEvent(ClientAddTag, agentID, req.ConvID, req.TagID)

	return jsonResponse(ctx, &Resp{Code: 0})
}

// DeleteTagFromVisitor
// DELETE /api/agent/client/:track_id/tags/:tag_id?conv_id=xxx
func (s *IMService) DeleteTagFromVisitor(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor_and_conv", "add_del_client_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	trackID := ctx.Param("track_id")
	tagID := ctx.Param("tag_id")
	convID := ctx.QueryParam("conv_id")
	if trackID == "" || tagID == "" || convID == "" {
		return invalidParameterResp(ctx, "track_id/tag_id/conv_id empty")
	}

	conv, err := models.ConversationByID(db.Mysql, convID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "conv not exists")
		}
		return dbErrResp(ctx, err.Error())
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, conv.EntID, trackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not exists")
		}
		return dbErrResp(ctx, err.Error())
	}

	tagRelation := &models.VisitorTagRelation{VisitorID: visitor.ID, TagID: tagID}
	if err := tagRelation.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	tags, err := models.VisitorTagRelationsByVisitorIDs(db.Mysql, []string{visitor.ID})
	if err == nil {
		var ts []string
		for _, tag := range tags {
			ts = append(ts, tag.TagID)
		}
		elasticsearch.UpdateConversationTags(elasticsearch.ESClient, convID, ts)
	}

	go s.sendUseClientTagEvent(ClientRemoveTag, conv.AgentID, conv.ID, tagID)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// SearchVisitors
// POST /admin/api/v1/visitors/search
// POST /api/users/search
func (s *IMService) SearchVisitors(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req := &SearchVisitorsReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	var queries, notQueries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("enterprise_id", entID))

	for _, rule := range req.Rules {
		var fieldName string
		switch rule.Attribute {
		case "landing_page_title":
			fieldName = "landing_page.landing_page_title"
		case "landing_page_url":
			fieldName = "landing_page.landing_page_url"
		case "source_domain":
			fieldName = "source.source_domain"
		case "source_keyword":
			fieldName = "source.source_keyword"
		case "source_source":
			fieldName = "source.source_source"
		case "source_url":
			fieldName = "source.source_url"
		default:
			fieldName = rule.Attribute
		}

		switch rule.Condition {
		case common.EqualOperator:
			queries = append(queries, elastic.NewTermQuery(fieldName, rule.Value))
		case common.NotEqualOperator:
			notQueries = append(notQueries, elastic.NewTermQuery(fieldName, rule.Value))
		case common.ContainsOperator:
			queries = append(queries, elastic.NewWildcardQuery(fieldName, fmt.Sprintf("%s*", rule.Value)))
		case common.NotContainsOperator:
			notQueries = append(notQueries, elastic.NewWildcardQuery(fieldName, fmt.Sprintf("%s*", rule.Value)))
		case common.EmptyOperator:
			queries = append(queries, elastic.NewTermQuery(fieldName, ""))
		case common.NotEmptyOperator:
			notQueries = append(notQueries, elastic.NewTermQuery(fieldName, ""))
		default:
			return invalidParameterResp(ctx, "unsupported operator")
		}
	}

	var offset, limit int
	offset = defaultOffset
	if req.Page >= 0 {
		offset = req.Page
	}

	limit = defaultLimit
	if req.Count > 0 && req.Count <= 1000 {
		limit = req.Count
	}

	query := elastic.NewBoolQuery().Must(queries...).MustNot(notQueries...)
	visitors, count, err := elasticsearch.SearchVisitors(elasticsearch.ESClient, query, offset, limit)
	if err != nil {
		log.Logger.Warnf("[elasticsearch.SearchVisitors] error: %v", err)
		return jsonResponse(ctx, &SearchVisitorsResp{Total: 0, Info: nil})
	}

	return jsonResponse(ctx, &SearchVisitorsResp{Total: count, Info: visitors, ExportLimit: 50000})
}

// GET /api/users/config
func (s *IMService) UsersConfig(ctx echo.Context) error {
	fss := map[string][]string{
		"main":   {"tag", "country", "qq", "weibo", "gender", "weixin", "province", "city"},
		"source": {"source_url", "source_keyword", "source_source", "os_family", "source_domain"},
	}

	fields := []*dto.VisitorFieldsSetting{
		{
			Category: "main",
			Key:      "name",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "tel",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "email",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "landing_page",
			Key:      "landing_page_title",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "landing_page",
			Key:      "landing_page_url",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "browser_family",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "address",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "comment",
			Type:     "string",
			Visible:  "yes",
		},
		{
			Category: "main",
			Key:      "created_on",
			Type:     "date",
			Visible:  "no",
		},
		{
			Category: "main",
			Key:      "updated_on",
			Type:     "date",
			Visible:  "no",
		},
		{
			Category: "main",
			Key:      "age",
			Type:     "int",
			Visible:  "no",
		},
	}
	for category, fs := range fss {
		for _, f := range fs {
			fields = append(fields, &dto.VisitorFieldsSetting{
				Category: category,
				Key:      f,
				Type:     "string",
				Visible:  "yes",
			})
		}

	}
	return jsonResponse(ctx, &dto.UserConfig{
		Setting:      fields,
		SyncSettings: nil,
	})
}

// OnlineVisitors ...
// GET /admin/api/v1/visitors
func (s *IMService) OnlineVisitors(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	ch := fmt.Sprintf(conf.IMConf.CentrifugoConf.VisitorQueueChannel, entID)
	visitIDs, err := s.imCli.ChannelUsers(context.Background(), ch)
	if err != nil {
		return internalServerErr(ctx, err.Error())
	}
	return jsonResponse(ctx, &Resp{Code: 0, Body: visitIDs})
}

// GenVisitorConnectionToken
// GET /api/v1/connection_token?visit_id=xxxxxx
func (s *IMService) GenVisitorConnectionToken(ctx echo.Context) (err error) {
	visitID := ctx.QueryParam("visit_id")
	if visitID == "" {
		return invalidParameterResp(ctx, "visit_id invalid")
	}

	token, err := newConnJwtToken(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		Subject:   visitID,
	})

	type resp struct {
		Token string `json:"token"`
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: &resp{Token: token}})
}

func CreateOrUpdateVisitorESDoc(db *sql.DB, entID, traceID string, modelVisitor *models.Visitor) {
	visits, err := models.VisitsByEntIDTraceID(db, entID, traceID)
	if err != nil {
		log.Logger.Errorf("VisitsByEntIDTraceID, ent: %s, trace: %s, error: %v", entID, traceID, err)
	}

	var visit *models.Visit
	if len(visits) > 0 {
		visit = visits[0]
	}

	tags, err := models.VisitorTagsByVisitorID(db, modelVisitor.ID)
	if err != nil {
		log.Logger.Errorf("VisitorTagsByVisitID error: %v", err)
	}

	doc := &dto.SearchVisitor{
		ID:           modelVisitor.ID,
		EnterpriseID: entID,
		TrackID:      traceID,
		Name:         modelVisitor.Name,
		Age:          modelVisitor.Age,
		Gender:       modelVisitor.Gender,
		Qq:           modelVisitor.QqNum,
		Email:        modelVisitor.Email,
		Tel:          modelVisitor.Mobile,
		Weibo:        modelVisitor.Weibo,
		Weixin:       modelVisitor.Wechat,
		Address:      modelVisitor.Address,
		Comment:      modelVisitor.Remark,
		Tag:          strings.Join(tags, ","),
		CreatedOn:    modelVisitor.CreatedAt,
		UpdatedOn:    modelVisitor.UpdatedAt,
	}

	if visit != nil {
		doc.VisitID = visit.ID
		doc.BrowserFamily = visit.BrowserFamily
		doc.OsFamily = visit.OsFamily
		doc.Country = visit.Country
		doc.Province = visit.Province
		doc.City = visit.City

		doc.Source = &dto.Source{
			SourceDomain:  visit.FirstPageSourceDomain,
			SourceKeyword: visit.FirstPageSourceKeyword,
			SourceSource:  visit.FirstPageSource,
			SourceURL:     visit.FirstPageSourceURL,
		}

		doc.LandingPage = &dto.LandingPage{
			LandingPageTitle: visit.FirstPageTitle,
			LandingPageURL:   visit.FirstPageURL,
		}
	}

	retryFunc := func() error {
		return elasticsearch.CreateOrUpdateVisitorDoc(elasticsearch.ESClient, doc)
	}

	err = retry.Do(
		retryFunc,
		retry.Attempts(3),
		retry.Delay(10*time.Millisecond),
	)
	if err != nil {
		log.Logger.Errorf("elastic search CreateOrUpdateVisitorDoc error: %v\n", err)
	}
}

func (s *IMService) onlineVisitors(entID string) (traceIDs []string, visitors []*models.Visitor, err error) {
	onlineVisitors, err := models.OnlineVisitorsByEntID(db.Mysql, entID)
	if err != nil {
		return nil, nil, err
	}

	for _, v := range onlineVisitors {
		traceIDs = append(traceIDs, v.TraceID)
	}

	visitors, err = models.VisitorsByTraceIDs(db.Mysql, traceIDs)
	return
}

func (s *IMService) sendVisitorAttrsUpdate(agentID string, trackID string, attrs map[string]interface{}) {
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	event := events.NewAttrsChange(agent, trackID, attrs)
	bs, err := common.Marshal(event)
	if err == nil {
		s.sendEventToAllAgents(agent.EntID, string(bs))
	}
}

type GetVisitorPositionReq struct {
	EntID   string `query:"ent_id"`
	TrackID string `query:"track_id"`
}

type GetVisitorPositionRes struct {
	Position int `json:"position"`
}

type GetVisitorPositionNotInQueueRes struct {
	InQueue bool `json:"in_queue"`
}

// TODO 当坐席对话数小于限制之后
//  1. queue/position 接口返回 {in_queue:false}
//  2. 将访客从排队中删除，发删除(visit_delete)排队的事件
// GET /client/queue/position
func (s *IMService) GetVisitorPosition(ctx echo.Context) error {
	req := &GetVisitorPositionReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.TrackID == "" {
		return invalidParameterResp(ctx, "track_id is empty")
	}

	agentIDs := models.GetEntOnlineAgentsV1(req.EntID)
	if len(agentIDs) > 0 {
		activeConvCnts, err := models.ActiveConversationNumOfAgents(db.Mysql, agentIDs)
		if err != nil {
			log.Logger.Warnf("[GetVisitorPosition] ActiveConversationNumOfAgents error=%v", err)
			return dbErrResp(ctx, err.Error())
		}

		serveLimits, err := models.ServeLimitByAgentIDs(db.Mysql, agentIDs)
		if err != nil {
			log.Logger.Warnf("[GetVisitorPosition] ServeLimitByAgentIDs error=%v", err)
			return dbErrResp(ctx, err.Error())
		}

		var hasAvailableAgent bool
		for _, agentID := range agentIDs {
			convCount := activeConvCnts[agentID]
			serveLimit := serveLimits[agentID]

			if convCount < serveLimit {
				hasAvailableAgent = true
				break
			}
		}

		if hasAvailableAgent {
			q := &models.VisitorQueue{TrackID: req.TrackID}
			if err := q.Delete(db.Mysql); err != nil {
				log.Logger.Warnf("delete viitor queue error: %v", err)
			}

			go s.sendVisitorStatusUpdateEvent(VisitorOffline, req.EntID, req.TrackID, 0)
			return jsonResponse(ctx, &GetVisitorPositionNotInQueueRes{InQueue: false})
		}
	}

	pos, err := models.GetVisitorPosition(db.Mysql, req.TrackID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &GetVisitorPositionRes{Position: pos})
}
