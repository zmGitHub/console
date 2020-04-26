package handler

import (
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type LeaveMessage struct {
	EntID   string `json:"ent_id"`
	TrackID string `json:"track_id"`
	Mobile  string `json:"mobile"`
	Content string `json:"content"`
}

type LeaveMessages struct {
	*models.LeaveMessage
	LastOptionAgent *string            `json:"last_option_agent"`
	CreatedAt       time.Time          `json:"created_at"` // created_at
	UpdatedAt       time.Time          `json:"updated_at"` // updated_at
	VisitInfo       *adapter.VisitInfo `json:"visit_info"`
}

type LeaveMessageResp struct {
	Total         int64            `json:"total"`
	LeaveMessages []*LeaveMessages `json:"leave_messages"`
}

// CreateLeaveMessageConfig
// POST /admin/api/v1/leave_message_config
func (s *IMService) CreateLeaveMessageConfig(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	lmc := &models.LeaveMessageConfig{}
	if err = ctx.Bind(lmc); err != nil {
		return
	}

	if lmc.FillContact != models.LeaveMessageConfigFillSingleContact &&
		lmc.FillContact != models.LeaveMessageConfigFillMultiContact {
		return invalidParameterResp(ctx, "unsupported fill_contact")
	}

	lmc.EntID = entID
	if err = models.CreateOrUpdateLeaveMessageConfig(db.Mysql, lmc); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: lmc})
}

// CreateLeaveMessage ...
// POST /api/leave_messages
func (s *IMService) CreateLeaveMessage(ctx echo.Context) (err error) {
	req := &LeaveMessage{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	now := time.Now().UTC()
	lm := models.LeaveMessage{
		ID:        common.GenUniqueID(),
		EntID:     req.EntID,
		TrackID:   req.TrackID,
		Mobile:    req.Mobile,
		Content:   req.Content,
		Status:    models.LeaveMessageUnHandledStatus,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err = lm.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// LeaveMessagesByEnt ...
// GET /admin/api/v1/leave_messages?status=open/closed
// GET /api/leave_messages
func (s *IMService) LeaveMessagesByEnt(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var trackIDs []string
	var all bool
	switch agent.PermsRangeType {
	case models.AgentPermsRangePersonalType:
		trackIDs, err = models.TrackIDsByAgentID(db.Mysql, agentID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}
	case models.AgentPermsRangeAllType:
		all = true
	default:
		groups, err := models.PermsRangeGroupIDsByAgentID(db.Mysql, agentID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		if len(groups) > 0 {
			agentIDs, err := models.AgentIDsByPermGroupIDs(db.Mysql, groups)
			if err != nil {
				return dbErrResp(ctx, err.Error())
			}

			trackIDs, err = models.TrackIDsByAgentIDs(db.Mysql, agentIDs)
			if err != nil {
				return dbErrResp(ctx, err.Error())
			}
		}
	}

	var result = &LeaveMessageResp{
		LeaveMessages: []*LeaveMessages{},
	}

	if len(trackIDs) == 0 && !all {
		return jsonResponse(ctx, result)
	}

	offset, limit, err := getOffsetLimitFromCtx(ctx)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var lms []*models.LeaveMessage
	var total int64
	if all {
		total, lms, err = models.EntLeaveMessagesByStatus(db.Mysql, entID, offset, limit)
	} else {
		total, lms, err = models.LeaveMessagesByTrackIDs(db.Mysql, trackIDs, offset, limit)
	}

	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	result.Total = total

	var visitTrackIDs []string
	var m = map[string]struct{}{}
	for _, lm := range lms {
		if _, ok := m[lm.TrackID]; ok {
			continue
		}

		visitTrackIDs = append(visitTrackIDs, lm.TrackID)
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, visitTrackIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, visitTrackIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	for _, lm := range lms {
		var visitor *models.Visitor
		var visit *models.Visit

		for _, vt := range visitors {
			if vt.TraceID == lm.TrackID {
				visitor = vt
				break
			}
		}

		for _, vt := range visits {
			if vt.TraceID == lm.TrackID {
				visit = vt
				break
			}
		}

		lmResp := &LeaveMessages{
			LeaveMessage: lm,
			CreatedAt:    common.ConvertUTCToLocal(lm.CreatedAt),
			UpdatedAt:    common.ConvertUTCToLocal(lm.UpdatedAt),
		}

		if visitor != nil && visit != nil {
			lmResp.VisitInfo = adapter.ConvertModelVisitToVisit(visitor, visit)
		}

		if lm.LastOptionAgent.Valid {
			lmResp.LastOptionAgent = &lm.LastOptionAgent.String
		}

		result.LeaveMessages = append(result.LeaveMessages, lmResp)
	}

	return jsonResponse(ctx, result)
}

// UpdateLeaveMessageStatus ...
// PUT /api/leave_messages/:msg_id
func (s *IMService) UpdateLeaveMessageStatus(ctx echo.Context) (err error) {
	type statusReq struct {
		Status string `json:"status"`
	}

	req := &statusReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.Status != models.LeaveMessageHandledStatus && req.Status != models.LeaveMessageUnHandledStatus {
		return invalidParameterResp(ctx, "unsupported status")
	}

	msgID := ctx.Param("msg_id")
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if err = models.UpdateLeaveMessageStatus(db.Mysql, msgID, req.Status, agentID); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	go s.sendLeaveMessageUpdateEvent(agentID, msgID)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func (s *IMService) sendLeaveMessageUpdateEvent(agentID, msgID string) {
	msg, err := models.LeaveMessageByID(db.Mysql, msgID)
	if err != nil {
		return
	}

	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, []string{msg.TrackID})
	if err != nil {
		return
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, []string{msg.TrackID})
	if err != nil {
		return
	}

	if len(visitors) == 0 || len(visits) == 0 {
		return
	}

	event := events.NewLeaveMessageEvent(agent, msg, dto.ModelVisitInfoToVisitInfo(visits[0], visitors[0]))
	bs, err := common.Marshal(event)
	if err == nil {
		s.sendEventToAllAgents(agent.EntID, string(bs))
	}
}
