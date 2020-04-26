package handler

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var entMediaURL = func(entID string) string {
	return fmt.Sprintf("https://chat-im.s3.ap-southeast-1.amazonaws.com/%s/files/", entID)
}

type sendMsgReq struct {
	EntID       string
	TraceID     string
	Creator     string
	Agents      []string
	ConvID      string
	FromType    string
	Content     string
	ContentType string
	Internal    bool
	System      bool
	extra       *fileExtra
}

type sendMsgAction struct {
	Action  string           `json:"action"`
	Message *adapter.Message `json:"message"`
}
type TimeLineReq struct {
	EntID      string `query:"ent_id"`
	TrackID    string `query:"track_id"`
	FilterType string `query:"filter_type"`
	Dt         string `query:"dt"`
}

type TimeLineResp struct {
	Events   []interface{}      `json:"events"`
	Messages []*adapter.Message `json:"messages"`
}

func (req *sendMsgReq) validate() error {
	if _, ok := models.MessageContentTypes[req.ContentType]; !ok {
		return fmt.Errorf("unsupported message content type")
	}

	if req.TraceID == "" || req.Content == "" {
		return fmt.Errorf("trace_id/content is invalid")
	}

	return nil
}

func (s *IMService) sendMsg(req *sendMsgReq) (channelMsg *adapter.Message, err error) {
	msg := &models.Message{
		ID:             common.GenUniqueID(),
		EntID:          req.EntID,
		TraceID:        req.TraceID,
		AgentID:        req.Creator,
		ConversationID: req.ConvID,
		FromType:       sql.NullString{String: req.FromType, Valid: true},
		CreatedAt:      time.Now().UTC(),
		Content:        sql.NullString{String: req.Content, Valid: true},
		ContentType:    sql.NullString{String: req.ContentType, Valid: true},
	}

	var msgContent = req.Content
	var mediaURL = entMediaURL(req.EntID)
	var extra interface{}

	t := time.Time{}
	switch req.ContentType {
	case models.MessagePictureContentType:
		fileName := req.Content
		if strings.HasPrefix(fileName, "http") {
			_, fileName = filepath.Split(fileName)
		}

		picExtra := &fileExtra{
			Filename: fileName,
			Size:     0,
			Type:     echo.MIMEOctetStream,
			ExpireAt: t.String(),
		}
		extra = picExtra
		setMessagePictureExtra(picExtra, msg)
		mediaURL = mediaURL + picExtra.Filename
	case models.MessageFileContentType:
		extra = req.extra
		setMessageFileExtra(req.extra, msg)
		mediaURL = mediaURL + req.extra.Filename
	case models.MessageRichTextType:
		extra, msgContent = setMessageRichTextExtra(msg.ID, req.EntID, req.Content, msg)
	}

	if req.System {
		msg.MsgType = sql.NullString{String: models.MessageMsgSystemType, Valid: true}
	} else {
		msg.MsgType = sql.NullString{String: models.MessageMsgPublicType, Valid: true}
		if req.Internal {
			msg.MsgType = sql.NullString{String: models.MessageMsgInternalType, Valid: true}
		}
	}

	if err = msg.Insert(db.Mysql); err != nil {
		return nil, err
	}

	if err = models.IncrMessageCount(db.Mysql, req.ConvID, req.FromType); err != nil {
		return nil, err
	}

	agent, err := models.AgentByID(db.Mysql, req.Creator)
	if err != nil {
		log.Logger.Warnf("get agent error: %v", err)
	}

	channelMsg = &adapter.Message{
		ID:             msg.ID,
		Action:         "message",
		Agent:          adapter.ConvertAgentToAgentInfo(agent, true),
		FromType:       req.FromType,
		Content:        msgContent,
		ContentType:    req.ContentType,
		CreatedOn:      common.ConvertUTCToTimeString(msg.CreatedAt),
		AgentID:        msg.AgentID,
		ConversationID: req.ConvID,
		ContentRobot:   nil,
		EnterpriseID:   msg.EntID,
		Extra:          extra,
		QuestionID:     "",
		ReadOn:         nil,
		MediaURL:       mediaURL,
		SubType:        "",
		TraceStart:     -1,
		TrackID:        msg.TraceID,
		Type:           "message",
	}

	if req.Internal {
		channelMsg.Type = models.MessageMsgInternalType
	}

	content, err := common.MarshalUnescape(&sendMsgAction{Action: "new_message", Message: channelMsg})
	if err != nil {
		log.Logger.Warnf("marshal channel msg: %v, error: %v", channelMsg, err)
		return channelMsg, err
	}

	if req.Internal {
		s.sendMessageToMultiAgents(req.Agents, content)
		return channelMsg, nil
	}

	s.sendMessageToMultiAgents(req.Agents, content)
	sendMessageToVisitor(s.imCli, req.TraceID, req.EntID, content)
	return channelMsg, nil
}

// MessagesByConversationID
// GET  /admin/api/v1/conversations/:conversation_id/messages
func (s *IMService) MessagesByConversationID(ctx echo.Context) (err error) {
	convID := ctx.Param("conversation_id")
	if convID == "" {
		return invalidParameterResp(ctx, "conversation_id is invalid")
	}

	offset, limit, err := getOffsetLimitFromCtx(ctx)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	messages, err := models.MessagesByConversationID(db.Mysql, convID, offset, limit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: messages})
}

// GET /client/timeline?ent_id=38523&track_id=1A1WbJY6Vo3DHED5pIy35xM7kDl&filter_type=after&dt=2019-05-02T11:16:10.749119
func (s *IMService) TimeLine(ctx echo.Context) (err error) {
	req := &TimeLineReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.Dt == "" {
		return invalidParameterResp(ctx, "bad time format")
	}

	dt, err := time.Parse("2006-01-02 15:04:05.999999999", req.Dt)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	msgs, err := models.TimeLimeMessages(db.Mysql, req.EntID, req.TrackID, dt)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	modelAgents, err := models.AgentsByMsgs(db.Mysql, msgs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentGroups, err := models.PermsRangeGroupIDsByAgents(db.Mysql, modelAgents)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	adapterAgents := adapter.ConvertAgentsToAdapterAgentsV1(modelAgents, agentGroups)

	result := &TimeLineResp{Events: []interface{}{}}
	result.Messages = convertModelMessagesToAdapterMessages(msgs, adapterAgents)
	return jsonResponse(ctx, result)
}

func (s *IMService) sendMessageToMultiAgents(agents []string, message string) {
	var wg sync.WaitGroup
	for _, agent := range agents {
		wg.Add(1)

		go func(id string) {
			sendMessageToAgent(s.imCli, id, message)
			wg.Done()
		}(agent)
	}

	wg.Wait()
}

func setMessageFileExtra(msgExtra *fileExtra, msg *models.Message) {
	if msgExtra == nil {
		return
	}

	content, err := common.Marshal(msgExtra)
	if err != nil {
		log.Logger.Warnf("marshal extra error: %v\n", err)
	}
	if content != "" {
		msg.Extra = sql.NullString{String: content, Valid: true}
	}
}

func setMessagePictureExtra(msgExtra *fileExtra, msg *models.Message) {
	if msgExtra == nil {
		return
	}

	content, err := common.Marshal(msgExtra)
	if err != nil {
		log.Logger.Warnf("marshal extra error: %v\n", err)
	}
	if content != "" {
		msg.Extra = sql.NullString{String: content, Valid: true}
	}
}

func setMessageRichTextExtra(msgID, entID, content string, msg *models.Message) (extra interface{}, promoContent string) {
	now := time.Now().UTC()
	promoContent = "[推广消息]"
	extra = &adapter.MessageExtra{
		ID:           msgID,
		EnterpriseID: entID,
		Content:      content,
		Summary:      "",
		Thumbnail:    "",
		CreatedOn:    *common.ConvertUTCToTimeString(now),
		UpdatedOn:    *common.ConvertUTCToTimeString(now),
	}

	msg.Content = sql.NullString{String: promoContent, Valid: true}
	extraStr, err := common.MarshalUnescape(extra)
	if err == nil {
		msg.Extra = sql.NullString{String: extraStr, Valid: true}
	}

	return
}
