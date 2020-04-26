package handler

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/dto"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

// GET /api/users/segments
func (s *IMService) UserSegments(ctx echo.Context) error {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	segments, err := models.UserSegmentsByEntID(db.Mysql, entID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &dto.GetSegmentsResp{}
	for _, seg := range segments {
		v, err := dto.ConvertToSegmentsResp(seg)
		if err != nil {
			return errResp(ctx, common.DecodeJSONErr, err.Error())
		}
		resp.Segments = append(resp.Segments, v)
	}

	return jsonResponse(ctx, resp)
}

// /api/users/segments/823
func (s *IMService) UpdateUserSegments(ctx echo.Context) error {
	segID := ctx.Param("seg_id")
	seg, err := models.UserSegmentByID(db.Mysql, segID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	req := &dto.UpdateSegmentsReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Name != "" {
		seg.Name = req.Name
		if err := seg.Update(db.Mysql); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	v, err := dto.ConvertToSegmentsResp(seg)
	if err != nil {
		return errResp(ctx, common.DecodeJSONErr, err.Error())
	}

	return jsonResponse(ctx, v)
}

// delete /api/users/segments/:seg_id
func (s *IMService) DeleteSegments(ctx echo.Context) error {
	segID := ctx.Param("seg_id")
	seg, err := models.UserSegmentByID(db.Mysql, segID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	if err := seg.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// POST /api/users/segments
func (s *IMService) CreateUserSegments(ctx echo.Context) error {
	req := &dto.Segments{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if len(req.Rules) == 0 {
		return invalidParameterResp(ctx, "invalid segments")
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	segments := models.UserSegment{
		ID:        common.GenUniqueID(),
		EntID:     entID,
		Name:      req.Name,
		CreatedAt: time.Now().UTC(),
	}

	bs, err := common.Marshal(req.Rules)
	if err != nil {
		return errResp(ctx, common.EncodeJSONErr, err.Error())
	}

	segments.Rules = sql.NullString{String: string(bs), Valid: true}

	if err := segments.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	resp := &dto.CreateSegmentsResp{
		ID:           segments.ID,
		CreatedOn:    *common.ConvertUTCToTimeString(segments.CreatedAt),
		UpdatedOn:    *common.ConvertUTCToTimeString(segments.CreatedAt),
		EnterpriseID: segments.EntID,
		Segments:     &dto.Segments{Name: segments.Name, Rules: req.Rules},
	}
	return jsonResponse(ctx, resp)
}
