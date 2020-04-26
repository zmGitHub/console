package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type addBlacklistReq struct {
	EntID   string `json:"ent_id"`
	TraceID string `json:"trace_id"`
	VisitID string `json:"visit_id"`
	ConvID  string `json:"conv_id"`
}

type addBlacklistReqV1 struct {
	TrackID string `json:"track_id"`
	VisitID string `json:"visit_id"`
	ConvID  string `json:"conv_id"`
}

type visitBlack struct {
	AgentID      string `json:"agent_id"`
	Avatar       string `json:"avatar"`
	City         string `json:"city"`
	ConvID       string `json:"conv_id"`
	CreatedOn    string `json:"created_on"`
	EnterpriseID string `json:"enterprise_id"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Province     string `json:"province"`
	TrackID      string `json:"track_id"`
	VisitID      string `json:"visit_id"`
}

type addBlacklistResp struct {
	Success    bool        `json:"success"`
	VisitBlack *visitBlack `json:"visit_black"`
}

type DeleteVisitBlacklistReq struct {
	TrackID string `json:"track_id"`
}

type GetVisitBlackListReq struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

// GetVisitBlackList
type GetVisitBlackListResp struct {
	Success        bool          `json:"success"`
	TotalCount     int           `json:"total_count"`
	VisitBlacklist []*visitBlack `json:"visit_blacklist"`
}

func (req *addBlacklistReq) validate() error {
	if req.EntID == "" || req.TraceID == "" || req.VisitID == "" || req.ConvID == "" {
		return fmt.Errorf("ent_id/trace_id/visit_id/conv_id is invalid")
	}

	return nil
}

type checkVisitorAllowedResp struct {
	Allow bool `json:"allow"`
}

// AddVisitorToBlacklist
// POST /admin/api/v1/blacklists
func (s *IMService) AddVisitorToBlacklist(ctx echo.Context) (err error) {
	req := &addBlacklistReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if err = req.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	blk := &models.VisitBlacklist{
		ID:        common.GenUniqueID(),
		EntID:     req.EntID,
		TraceID:   req.TraceID,
		VisitID:   req.VisitID,
		ConvID:    req.ConvID,
		CreatedAt: time.Now().UTC(),
	}
	if err = blk.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: blk})
}

// PUT /api/agent/visit/blacklist
func (s *IMService) AddVisitorToBlacklistV1(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor_and_conv", "add_del_client_blacklist"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &addBlacklistReqV1{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.TrackID == "" || req.VisitID == "" || req.ConvID == "" {
		return invalidParameterResp(ctx, "track_id/visit_id/conv_id invalid")
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, req.TrackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not exists")
		}
		return dbErrResp(ctx, err.Error())
	}

	visit, err := models.VisitByID(db.Mysql, visitor.LastVisitID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	blk := &models.VisitBlacklist{
		ID:        common.GenUniqueID(),
		EntID:     entID,
		TraceID:   req.TrackID,
		VisitID:   req.VisitID,
		ConvID:    req.ConvID,
		CreatedAt: time.Now().UTC(),
	}
	if err = blk.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &addBlacklistResp{
		Success: true,
		VisitBlack: &visitBlack{
			AgentID:      agentID,
			Avatar:       visitor.Avatar,
			City:         visit.City,
			ConvID:       req.ConvID,
			CreatedOn:    *common.ConvertUTCToTimeString(blk.CreatedAt),
			EnterpriseID: entID,
			ID:           blk.ID,
			Name:         visitor.Name,
			Province:     visit.Province,
			TrackID:      req.TrackID,
			VisitID:      req.VisitID,
		},
	}

	go s.sendVisitBlackEvent(VisitBlackAdd, req.TrackID, visitor.ID, agentID, *resp.VisitBlack)

	return jsonResponse(ctx, resp)
}

// DeleteVisitBlacklist
// DELETE /api/agent/visit/blacklist
func (s *IMService) DeleteVisitBlacklist(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "visitor_and_conv", "add_del_client_blacklist"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &DeleteVisitBlacklistReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.TrackID == "" {
		return invalidParameterResp(ctx, "track_id empty")
	}

	blackList, err := models.VisitBlacklistByEntIDTraceID(db.Mysql, entID, req.TrackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的黑名单")
		}

		return dbErrResp(ctx, err.Error())
	}

	if err := blackList.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visitor, err := models.VisitorByEntIDTraceID(db.Mysql, entID, req.TrackID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visitor not exists")
		}
		return dbErrResp(ctx, err.Error())
	}

	visit, err := models.VisitByID(db.Mysql, visitor.LastVisitID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "visit not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	visitBlack := &visitBlack{
		AgentID:      agentID,
		Avatar:       visitor.Avatar,
		City:         visit.City,
		ConvID:       "",
		CreatedOn:    *common.ConvertUTCToTimeString(visitor.CreatedAt),
		EnterpriseID: entID,
		ID:           visitor.ID,
		Name:         visitor.Name,
		Province:     visit.Province,
		TrackID:      req.TrackID,
		VisitID:      visit.ID,
	}

	go s.sendVisitBlackEvent(VisitBlackRemove, req.TrackID, visitor.ID, agentID, *visitBlack)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// CheckChatLinkAllowed check if a visitor can start conversation
// GET /api/v1/enterprises/:ent_id/visitor_allowed?trace_id=xxxxx
func (s *IMService) CheckVisitorAllowed(ctx echo.Context) (err error) {
	entID := ctx.Param("ent_id")
	traceID := ctx.QueryParam("trace_id")
	if entID == "" || traceID == "" {
		return invalidParameterResp(ctx, "ent_id/trace_id is empty")
	}

	blk, err := models.VisitBlacklistByEntIDTraceID(db.Mysql, entID, traceID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if blk == nil {
		return jsonResponse(ctx, &Resp{Code: 0, Body: &checkVisitorAllowedResp{Allow: true}})
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: &checkVisitorAllowedResp{Allow: false}})
}

//GET /api/agent/visit/blacklist?limit=15&offset=0&browser_id=agent1550496929419&v=1550498042183
func (s *IMService) GetVisitBlackList(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_blacklist"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &GetVisitBlackListReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var offset, limit = 0, 15
	if req.Offset > 0 {
		offset = req.Offset
	}
	if req.Limit >= 1 {
		limit = req.Limit
	}

	list, err := models.VisitorBlacklistsByEntID(db.Mysql, entID, offset, limit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var vbList = make([]*visitBlack, 0, len(list))
	for _, vb := range list {
		vbList = append(vbList, &visitBlack{
			AgentID:      vb.AgentID,
			Avatar:       "",
			City:         "",
			ConvID:       vb.ConvID,
			CreatedOn:    *common.ConvertUTCToTimeString(vb.CreatedAt),
			EnterpriseID: entID,
			ID:           vb.ID,
			Name:         "",
			Province:     "",
			TrackID:      vb.TraceID,
			VisitID:      vb.VisitID,
		})
	}

	return jsonResponse(ctx, &GetVisitBlackListResp{Success: true, TotalCount: len(vbList), VisitBlacklist: vbList})
}
