package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/gooxml/spreadsheet"

	"bitbucket.org/forfd/custm-chat/webim/auth"
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type CreateQuickReplyGroupReq struct {
	Title     string `json:"title"`
	Rank      int    `json:"rank"`
	OwnerType string `json:"owner_type"`
}

type UpdateQuickReplyGroupReq struct {
	Title    string `json:"title"`
	Position *int   `json:"position"`
}

type QuickReplyItem struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	HotKey      []int  `json:"hot_key"`
}

type QuickReplyGroup struct {
	*models.QuickreplyGroup
	Items []*models.QuickreplyItem `json:"items"`
}

type GetQuickRepliesResp struct {
	Groups []*QuickReplyGroup `json:"groups"`
}

// POST /admin/api/v1/quickreplies/
// POST /api/agent/quick_reply_groups
func (s *IMService) CreateQuickReplyGroup(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &CreateQuickReplyGroupReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}
	if req.Title == "" {
		return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: "title is invalid"})
	}

	var ownerID = agentID
	if req.OwnerType == "enterprise" {
		ownerID = entID
	}

	lastRank, _ := models.LastQkGroupRank(db.Mysql, entID)

	now := time.Now().UTC()
	reply := &models.QuickreplyGroup{
		ID:          common.GenUniqueID(),
		EntID:       entID,
		Title:       req.Title,
		Rank:        lastRank + 100000,
		CreatedBy:   ownerID,
		CreatorType: req.OwnerType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err = reply.Insert(db.Mysql); err != nil {
		return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: err.Error()})
	}

	group := adapter.ConvertModelQKReplyGroupToGroup(reply)

	go s.sendQkReplyGroupEvent(agentID, QuickReplyGroupCreate, group)

	return jsonResponse(ctx, group)
}

// PUT /admin/api/v1/quickreplies/:group_id
// PUT /api/agent/quick_reply_groups/:group_id
func (s *IMService) UpdateQuickReplyGroup(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	groupID := ctx.Param("group_id")
	req := &UpdateQuickReplyGroupReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if groupID == "" {
		return invalidParameterResp(ctx, "group_id is empty")
	}

	reply, err := models.QuickreplyGroupByID(db.Mysql, groupID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	if req.Title != "" {
		reply.Title = req.Title
	}

	if req.Position != nil {
		groups, err := models.QuickGroups(db.Mysql, entID, reply.CreatorType)
		if err != nil {
			return dbErrResp(ctx, err.Error())
		}

		var currentPosition = -1
		for i, g := range groups {
			if g.ID == groupID {
				currentPosition = i
				break
			}
		}

		targetPosition := *req.Position
		groupsRank := QuickReplyGroupsRankImpl(groups)
		newRank := getNewRank(currentPosition, targetPosition, groupsRank)
		if newRank != -1 {
			reply.Rank = newRank
		}
	}

	if err = reply.Update(db.Mysql); err != nil {
		return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: err.Error()})
	}

	group := adapter.ConvertModelQKReplyGroupToGroup(reply)
	go s.sendQkReplyGroupEvent(agentID, QuickReplyGroupUpdate, group)

	return jsonResponse(ctx, group)
}

// DELETE /admin/api/v1/quickreplies/:group_id
// DELETE /api/agent/quick_reply_groups/:group_id
func (s *IMService) DeleteQuickReplyGroup(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	groupID := ctx.Param("group_id")
	if groupID == "" {
		return invalidParameterResp(ctx, "group_id is invalid")
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		if err := rollBackOrCommit(tx, dbErr); err != nil {
			log.Logger.Warnf("DeleteQuickReplyGroup rollBackOrCommit error: %v\n", err)
		}
	}()

	group := &models.QuickreplyGroup{ID: groupID, EntID: entID}
	if dbErr = group.Delete(tx); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}
	if dbErr = models.DeleteQuickReplyItemsByGroupID(tx, groupID); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	go s.sendQkReplyGroupEvent(agentID, QuickReplyGroupDelete, adapter.ConvertModelQKReplyGroupToGroup(group))

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// CreateQuickReplyItem ...
// POST /admin/api/v1/quickreplies/:group_id/items
// POST /api/agent/quick_replies
func (s *IMService) CreateQuickReplyItem(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	req := &adapter.QkReply{}
	if err = ctx.Bind(req); err != nil {
		return
	}
	if req.Content == "" || req.GroupID == "" {
		return invalidParameterResp(ctx, "content/group_id is invalid")
	}

	lastRank, _ := models.LastQkReplyRank(db.Mysql, req.GroupID)
	now := time.Now().UTC()

	hotKey, err := common.Marshal(req.HotKey)
	if err != nil {
		hotKey = `[]`
	}

	item := &models.QuickreplyItem{
		ID:                common.GenUniqueID(),
		QuickreplyGroupID: req.GroupID,
		Title:             req.Title,
		Content:           req.Content,
		ContentType:       req.ContentType,
		RichContent:       sql.NullString{String: req.RichContent, Valid: true},
		Rank:              lastRank + 100000,
		HotKey:            hotKey,
		CreatedBy:         agentID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err = item.Insert(db.Mysql); err != nil {
		return jsonResponse(ctx, &ErrMsg{Code: common.DBErr, Message: err.Error()})
	}

	reply := adapter.ConvertModelQkItemToReply(entID, item)
	go s.sendQuickReplyEvent(agentID, QuickReplyCreate, reply)

	return jsonResponse(ctx, reply)
}

// UpdateQuickReplyItem
// PUT /admin/api/v1/quickreplies/items/:item_id
// PUT /api/agent/quick_replies/:reply_id
func (s *IMService) UpdateQuickReplyItem(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	itemID := ctx.Param("reply_id")
	req := &QuickReplyItem{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if itemID == "" || req.Content == "" {
		return invalidParameterResp(ctx, "reply_id/content is invalid")
	}

	item, err := models.QuickreplyItemByID(db.Mysql, itemID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	item.Title = req.Title
	item.Content = req.Content

	if req.ContentType != "" {
		item.ContentType = req.ContentType
	}

	if len(req.HotKey) > 0 {
		content, err := common.Marshal(req.HotKey)
		if err == nil {
			item.HotKey = content
		}
	}

	if err = item.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	reply := adapter.ConvertModelQkItemToReply(entID, item)
	go s.sendQuickReplyEvent(agentID, QuickReplyUpdate, reply)

	return jsonResponse(ctx, reply)
}

// DeleteQuickReplyItem
// DELETE /admin/api/v1/quickreplies/items/:reply_id
// DELETE /api/agent/quick_replies/:reply_id
func (s *IMService) DeleteQuickReplyItem(ctx echo.Context) (err error) {
	itemID := ctx.Param("reply_id")
	if itemID == "" {
		return invalidParameterResp(ctx, "item_id is empty")
	}

	item, err := models.QuickreplyItemByID(db.Mysql, itemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return invalidParameterResp(ctx, "不存在的快捷回复")
		}
		return invalidParameterResp(ctx, err.Error())
	}

	if err = item.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	reply := adapter.ConvertModelQkItemToReply(entID, item)
	go s.sendQuickReplyEvent(agentID, QuickReplyDelete, reply)

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// GetQuickReplies 获取快捷回复 quick replies
// GET /admin/api/v1/quickreplies/
// GET /api/agent/quick_replies
func (s *IMService) GetQuickReplies(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)
	if !conf.IMConf.Debug {
		if msg := hasPerm(entID, agentID, "online_agent_config", "config_ent_quick_reply"); msg != nil {
			return noPermResp(ctx, msg)
		}
	}

	groups, errMsg := s.getEntQuickReplies(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	return jsonResponse(ctx, convertQkGroupsToAdapterGroups(groups))
}

func convertQkGroupsToAdapterGroups(groups []*QuickReplyGroup) *adapter.QkReplyGroupResp {
	resp := &adapter.QkReplyGroupResp{
		AgentQuickReplyGroup: make([]*adapter.QkReplyGroup, 0, len(groups)),
		EntQuickReplyGroup:   make([]*adapter.QkReplyGroup, 0, len(groups)),
	}

	for _, group := range groups {
		quickReplies := make([]*adapter.QkReply, len(group.Items))

		for j, item := range group.Items {
			var hotKeys []int
			if item.HotKey != "" {
				if err := common.Unmarshal(item.HotKey, &hotKeys); err != nil {
					log.Logger.Warnf("[convertQkGroupsToAdapterGroups] Unmarshal error: %v", err)
				}
			}

			quickReplies[j] = &adapter.QkReply{
				ID:                 item.ID,
				EnterpriseID:       group.EntID,
				GroupID:            item.QuickreplyGroupID,
				Content:            item.Content,
				ContentType:        item.ContentType,
				CreatedOn:          item.CreatedAt,
				HotKey:             hotKeys,
				KnowledgeConverted: false,
				LastUpdated:        item.UpdatedAt,
				Rank:               item.Rank,
				RichContent:        item.RichContent.String,
				Title:              item.Title,
			}
		}

		if group.CreatorType == models.QuickReplyEnterpriseCreatorType {
			resp.EntQuickReplyGroup = append(resp.EntQuickReplyGroup, &adapter.QkReplyGroup{
				ID:           group.ID,
				EnterpriseID: group.EntID,
				OwnerID:      group.EntID,
				OwnerType:    models.QuickReplyEnterpriseCreatorType,
				Rank:         group.Rank,
				Title:        group.Title,
				CreatedOn:    group.CreatedAt,
				LastUpdated:  group.UpdatedAt,
				QuickReplies: quickReplies,
			})
			continue
		}

		resp.AgentQuickReplyGroup = append(resp.AgentQuickReplyGroup, &adapter.QkReplyGroup{
			ID:           group.ID,
			EnterpriseID: group.EntID,
			OwnerID:      group.CreatedBy,
			OwnerType:    models.QuickReplyAgentCreatorType,
			Rank:         group.Rank,
			Title:        group.Title,
			CreatedOn:    group.CreatedAt,
			LastUpdated:  group.UpdatedAt,
			QuickReplies: quickReplies,
		})
	}

	return resp
}

// GET /admin/api/v1/quickreplies/export
// GET api/agent/quick_replies/export
func (s *IMService) ExportQuickReply(ctx echo.Context) (err error) {
	ownerType := ctx.QueryParam("owner_type")
	if ownerType == "" {
		ownerType = models.QuickReplyEnterpriseCreatorType
	}

	tk := ctx.QueryParam("tk")
	if tk == "" {
		return &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  "invalid or expired jwt",
			Internal: fmt.Errorf("invalid or expired jwt"),
		}
	}

	entID, _, err := auth.ParseUserToken(conf.IMConf.JWTKey, tk)
	if err != nil {
		return err
	}

	groups, errMsg := s.getEntQuickReplies(entID)
	if errMsg != nil {
		return jsonResponse(ctx, errMsg)
	}

	if len(groups) == 0 {
		return invalidParameterResp(ctx, "无快速回复")
	}

	var resultGroups []*QuickReplyGroup
	for _, group := range groups {
		if group.CreatorType == ownerType {
			resultGroups = append(resultGroups, group)
		}
	}

	if len(resultGroups) == 0 {
		return invalidParameterResp(ctx, "无相关快速回复")
	}

	workbook, err := s.saveQuickReplies(entID, resultGroups)
	if err != nil {
		return errResp(ctx, common.ExportFileErr, err.Error())
	}

	contentBytes := &bytes.Buffer{}
	if err := workbook.Save(contentBytes); err != nil {
		return errResp(ctx, 0, err.Error())
	}

	w := ctx.Response()
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	r := ctx.Request()
	t := time.Now()
	layout := "2006-01-02 15:04:05"
	fileName := fmt.Sprintf("快速回复-企通-%s-%s.xlsx", t.Format(layout), ownerType)
	http.ServeContent(w, r, fileName, time.Time{}, bytes.NewReader(contentBytes.Bytes()))
	return nil
}

// POST /api/agent/quick_replies/import
func (s *IMService) ImportQuickReply(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	agentID := ctx.Get(middleware.AgentIDKey).(string)

	// fileForm := ctx.FormValue("file_from")
	ownerType := ctx.FormValue("owner_type")
	if ownerType != models.QuickReplyAgentCreatorType && ownerType != models.QuickReplyEnterpriseCreatorType {
		ownerType = models.QuickReplyEnterpriseCreatorType
	}

	file, err := ctx.FormFile("quick_replies")
	if err != nil {
		return errResp(ctx, common.UploadFileErr, err.Error())
	}

	if file.Size >= defaultMaxFileSize {
		return errResp(ctx, common.UploadFileErr, "file too large")
	}

	src, err := file.Open()
	if err != nil {
		return errResp(ctx, common.UploadFileErr, "open file error")
	}

	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			log.Logger.Error("close file error: ", closeErr)
		}
	}()

	groups, err := s.parseQuickReplyFile(entID, agentID, src, ownerType)
	if err != nil {
		log.Logger.Error("parseQuickReplyFile error: ", err)
		return jsonResponse(ctx, &SuccessResp{Success: false})
	}

	if len(groups) == 0 {
		return jsonResponse(ctx, &SuccessResp{Success: false})
	}

	var qkGroups []*models.QuickreplyGroup
	var qkItems []*models.QuickreplyItem

	for _, group := range groups {
		qkGroups = append(qkGroups, group.QuickreplyGroup)
		for _, item := range group.Items {
			qkItems = append(qkItems, item)
		}
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var dbErr error
	defer func() {
		if err := rollBackOrCommit(tx, dbErr); err != nil {
			log.Logger.Errorf("ImportQuickReply rollBackOrCommit error: %v", err)
		}
	}()

	if dbErr = models.InsertQuickReplyGroups(tx, qkGroups); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}
	if dbErr = models.InsertQuickReplyItems(tx, qkItems); dbErr != nil {
		return dbErrResp(ctx, dbErr.Error())
	}

	go s.sendQuickReplyUpload(entID, agentID)
	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func (s *IMService) sendQuickReplyUpload(entID, agentID string) {
	event := events.NewQuickReplayUploadEvent(&events.QuickReplayUpload{
		Action:       "quick_replies_refresh",
		EnterpriseID: entID,
		ID:           "",
		TraceStart:   float64(time.Now().UnixNano()),
	})

	agentIDs, err := models.AgentIDsByEntID(db.Mysql, entID)
	if err != nil {
		agentIDs = []string{agentID}
	}

	s.sendMessageToMultiAgents(agentIDs, event)
}

func (s *IMService) parseQuickReplyFile(entID, agentID string, r io.Reader, ownerType string) (groupContainer map[string]*QuickReplyGroup, err error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		log.Logger.Warn("excelize error", err)
		return
	}

	groupContainer = make(map[string]*QuickReplyGroup)
	now := time.Now().UTC()
	rows := xlsx.GetRows("快速回复")

	if len(rows) > 1 {
		rows = rows[1:] // 排除列名
	}

	var createBy = agentID
	if ownerType == models.QuickReplyEnterpriseCreatorType {
		createBy = entID
	}

	// rows: 组名  标题	内容
	for _, row := range rows {
		if len(row) >= 3 {
			groupTitle, title, content := row[0], row[1], row[2]
			if groupTitle == "" {
				continue
			}

			qkGroup, ok := groupContainer[groupTitle]
			if !ok {
				qkGroup = &QuickReplyGroup{QuickreplyGroup: &models.QuickreplyGroup{
					ID:          common.GenUniqueID(),
					EntID:       entID,
					Title:       groupTitle,
					Rank:        100000,
					CreatedBy:   createBy,
					CreatorType: ownerType,
					CreatedAt:   now,
					UpdatedAt:   now,
				}}
			}

			qkGroup.Items = append(qkGroup.Items, &models.QuickreplyItem{
				ID:                common.GenUniqueID(),
				QuickreplyGroupID: qkGroup.ID,
				Title:             title,
				Content:           content,
				ContentType:       models.QuickReplyTextContentType,
				RichContent:       sql.NullString{String: "", Valid: true},
				Rank:              100000,
				HotKey:            "[]",
				CreatedBy:         createBy,
				CreatedAt:         now,
				UpdatedAt:         now,
			})

			groupContainer[groupTitle] = qkGroup
		}
	}

	return groupContainer, nil
}

func (s *IMService) saveQuickReplies(entID string, groups []*QuickReplyGroup) (ss *spreadsheet.Workbook, err error) {
	if len(groups) == 0 {
		return
	}

	ss = spreadsheet.New()
	// add a single sheet
	sheet := ss.AddSheet()

	// rows
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.SetString("组名")
	cell = row.AddCell()
	cell.SetString("标题")
	cell = row.AddCell()
	cell.SetString("内容")

	for _, group := range groups {
		for _, item := range group.Items {
			row := sheet.AddRow()
			// and cells
			cell := row.AddCell()
			cell.SetString(group.Title)

			cell = row.AddCell()
			cell.SetString(item.Title)

			cell = row.AddCell()
			cell.SetString(item.Content)
		}
	}

	if err = ss.Validate(); err != nil {
		return
	}

	return
}

func (s *IMService) getEntQuickReplies(entID string) (groups []*QuickReplyGroup, errMsg *ErrMsg) {
	replies, err := models.QuickReplyGroupsByEntID(db.Mysql, entID)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	var groupIDs []string
	for _, reply := range replies {
		groupIDs = append(groupIDs, reply.ID)
	}

	items, err := models.QuickReplyItemsByGroupIDs(db.Mysql, groupIDs)
	if err != nil {
		return nil, &ErrMsg{Code: common.DBErr, Message: err.Error()}
	}

	for _, reply := range replies {
		var groupItems []*models.QuickreplyItem
		for _, item := range items {
			if item.QuickreplyGroupID == reply.ID {
				groupItems = append(groupItems, item)
			}
		}

		groups = append(groups, &QuickReplyGroup{
			QuickreplyGroup: reply,
			Items:           groupItems,
		})
	}

	return groups, nil
}
