package handler

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type UpdateVisitorTagReq struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position *int   `json:"position"`
}

// CreateVisitorTag ...
// POST /admin/api/v1/visitor_tags
// POST /api/agent/client_tags
func (s *IMService) CreateVisitorTag(ctx echo.Context) (err error) {
	agentInfo := getAgentInfoFromJwtToken(ctx)
	if !conf.IMConf.Debug {
		if msg := hasPerm(agentInfo.EntID, agentInfo.UserID, "online_agent_config", "config_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	tag := &models.VisitorTag{}
	if err = ctx.Bind(tag); err != nil {
		return
	}
	if err = tag.Validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	lastRank, _ := models.LastTagRank(db.Mysql, agentInfo.EntID)

	now := time.Now().UTC()
	tag.ID = common.GenUniqueID()
	tag.Creator = agentInfo.UserID
	tag.EntID = agentInfo.EntID
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.UseCount = 0
	tag.Rank = lastRank + 100000
	if err = tag.Insert(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	clientTag := adapter.ConvertModelTagToClientTag(tag)
	go s.sendClientTagEvent(agentInfo.UserID, ClientTagCreate, clientTag)

	return jsonResponse(ctx, clientTag)
}

// GET /admin/api/v1/visitor_tags
// GET /api/agent/client_tags
func (s *IMService) GetEntVisitorTags(ctx echo.Context) (err error) {
	agentInfo := getAgentInfoFromJwtToken(ctx)
	if !conf.IMConf.Debug {
		if msg := hasPerm(agentInfo.EntID, agentInfo.UserID, "online_agent_config", "config_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	tags, err := models.VisitorTagsByEntID(db.Mysql, agentInfo.EntID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, adapter.ConvertVisitorTagsToClientTags(tags))
}

// PUT /admin/api/v1/visitor_tags/:tag_id
// PUT /api/agent/client_tags/:tag_id
func (s *IMService) UpdateVisitorTag(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	tagID := ctx.Param("tag_id")
	if tagID == "" {
		return invalidParameterResp(ctx, "tag_id is empty")
	}

	req := &UpdateVisitorTagReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	vt, err := models.VisitorTagByID(db.Mysql, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return dbErrResp(ctx, "tag not found")
		}

		return dbErrResp(ctx, err.Error())
	}

	var updateCount int
	if req.Name != "" {
		vt.Name = req.Name
		updateCount++
	}
	if req.Color != "" {
		vt.Color = req.Color
		updateCount++
	}

	if req.Position != nil {
		clientTags, err := models.GetClientTagRanks(db.Mysql, entID)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		var currentPosition = -1
		for i, tag := range clientTags {
			if tag.ID == tagID {
				currentPosition = i
				break
			}
		}

		targetPosition := *req.Position
		clientTagsRank := ClientTagsRankImpl(clientTags)
		newRank := getNewRank(currentPosition, targetPosition, clientTagsRank)
		if newRank != -1 {
			updateCount++
			vt.Rank = newRank
		}
	}

	if updateCount > 0 {
		vt.UpdatedAt = time.Now().UTC()
		if err = vt.Update(db.Mysql); err != nil {
			return dbErrResp(ctx, err.Error())
		}
	}

	clientTag := adapter.ConvertModelTagToClientTag(vt)
	go s.sendClientTagEvent(agentID, ClientTagUpdate, clientTag)

	return jsonResponse(ctx, clientTag)
}

// DELETE /admin/api/v1/visitor_tags/:tag_id
// DELETE /api/agent/client_tags/:tag_id
func (s *IMService) DeleteVisitorTag(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_tag"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	tagID := ctx.Param("tag_id")
	if tagID == "" {
		return invalidParameterResp(ctx, "tag_id is empty")
	}

	tag, err := models.VisitorTagByID(db.Mysql, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "tag not exists")
		}

		return dbErrResp(ctx, err.Error())
	}

	if err = tag.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	clientTag := adapter.ConvertModelTagToClientTag(tag)
	go s.sendClientTagEvent(agentID, ClientTagDelete, clientTag)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}
