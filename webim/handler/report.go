package handler

import (
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type ReportReq struct {
	StartTime time.Time `query:"start_time"`
	EndTime   time.Time `query:"end_time"`
}

type VisitReportResp struct {
	VisitorCount int64 `json:"visitor_count"`
	VisitNum     int64 `json:"visit_num"`
}

// VisitReport 报表
// GET /admin/api/v1/enterprise/reports/visit
func (s *IMService) VisitReport(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req := &ReportReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return invalidParameterResp(ctx, "invalid start_time/end_time")
	}

	resp := &VisitReportResp{}
	resp.VisitorCount, resp.VisitNum, err = models.GetVisitStatsByDateRange(db.Mysql, entID, req.StartTime, req.EndTime)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: resp})
}

// ConversationReport
// GET /admin/api/v1/enterprise/reports/conversation
func (s *IMService) ConversationReport(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	req := &ReportReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return invalidParameterResp(ctx, "invalid start_time/end_time")
	}

	stats, err := models.GetConversationStatsByDateRange(db.Mysql, entID, req.StartTime, req.EndTime)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: stats})
}
