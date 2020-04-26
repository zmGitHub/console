package handler

import (
	"sort"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type GetHistoryConversationsReq struct {
	EntID   string `query:"ent_id"`
	TrackID string `query:"track_id"`
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
}

type ConvStreamsReq struct {
	ConvID  string `query:"conv_id"`
	Include int    `query:"include"`
	Order   int    `query:"order"`
	Limit   int    `query:"limit"`
	Type    string `query:"type"`
}

// GET /client/history_conversation?ent_id=5869&track_id=1EACfgGhNoogG9YFlTVr2OJt9lK&page=1&limit=5
func (s *IMService) GetHistoryConversation(ctx echo.Context) (err error) {
	req := &GetHistoryConversationsReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}
	if req.EntID == "" || req.TrackID == "" {
		return invalidParameterResp(ctx, "ent_id/track_id invalid")
	}
	var offset, limit = 0, 5
	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.Page > 1 {
		offset = (req.Page - 1) * limit
	}

	conversations, err := models.ConversationsByEntIDTraceID(db.Mysql, req.EntID, req.TrackID, offset, limit)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	var traceIDs []string
	var convIDs []string
	for _, conv := range conversations {
		convIDs = append(convIDs, conv.ID)
		traceIDs = append(traceIDs, conv.TraceID)
	}

	visits, err := models.VisitsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	visitors, err := models.VisitorsByTraceIDs(db.Mysql, traceIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	tagRelations, err := models.VisitorTagRelationsByVisitors(db.Mysql, visitors)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	msgs, err := models.MessagesByConversationIDs(db.Mysql, convIDs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agents, err := models.AgentsByMsgs(db.Mysql, msgs)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	agentGroups, err := models.PermsRangeGroupIDsByAgents(db.Mysql, agents)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}
	adapterAgents := adapter.ConvertAgentsToAdapterAgentsV1(agents, agentGroups)

	var adapterConvs = make([]*adapter.HistoryConversation, len(conversations))
	for i, conversation := range conversations {

		var visit *models.Visit
		for _, vst := range visits {
			if vst.TraceID == conversation.TraceID {
				visit = vst
				break
			}
		}

		var visitor *models.Visitor
		for _, vstor := range visitors {
			if vstor.TraceID == conversation.TraceID {
				visitor = vstor
				break
			}
		}

		var convMessages []*models.Message
		for _, msg := range msgs {
			if msg.MsgType.String == models.MessageMsgInternalType {
				continue
			}

			if msg.ConversationID == conversation.ID {
				convMessages = append(convMessages, msg)
			}
		}

		sort.SliceStable(convMessages, func(i, j int) bool {
			msg1, msg2 := convMessages[i], convMessages[j]
			return msg1.CreatedAt.After(msg2.CreatedAt)
		})

		historyConv := &adapter.HistoryConversation{
			Conversation: adapter.ModelConversationToConversation(conversation, visit.ID, visit, visitor, tagRelations),
			Messages:     convertModelMessagesToAdapterMessages(convMessages, adapterAgents),
		}

		if offset == 0 && i == 0 && len(historyConv.Messages) > 0 {
			historyConv.Messages = historyConv.Messages[1:]
		}
		adapterConvs[i] = historyConv
	}

	return jsonResponse(ctx, &adapter.HistoryConversationsResp{Conversations: adapterConvs})
}

// GET /api/conversation/:track_id/streams?conv_id=458953411&include=1&order=0&limit=1&type=earlier&browser_id=agent1548679440405
func (s *IMService) GetStreamConversations(ctx echo.Context) (err error) {
	trackID := ctx.Param("track_id")
	if trackID == "" {
		return invalidParameterResp(ctx, "track id invalid")
	}

	req := &ConvStreamsReq{}
	if err = ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	result, errMsg := s.getVisitorConversationStreams(trackID, req)
	if errMsg != nil {
		return errResp(ctx, errMsg.Code, errMsg.Message)
	}

	return jsonResponse(ctx, result)
}

func convertModelMessagesToAdapterMessages(msgs []*models.Message, agents []*adapter.Agent) []*adapter.Message {
	if len(msgs) == 0 {
		return nil
	}

	findAgent := func(id string) *adapter.Agent {
		for _, agent := range agents {
			if agent.ID == id {
				return agent
			}
		}

		return nil
	}

	var adapterMessages []*adapter.Message
	for _, msg := range msgs {
		modelAgent := findAgent(msg.AgentID)
		adapterMessages = append(adapterMessages, modelMsgToAdapterMessage(modelAgent, msg))
	}

	return adapterMessages
}

func modelMsgToAdapterMessage(agent *adapter.Agent, msg *models.Message) *adapter.Message {
	var extra interface{}
	var err error
	if msg.Extra.Valid {
		if msg.ContentType.String == models.MessageRichTextType {
			var obj = &adapter.MessageExtra{}
			if err := common.Unmarshal(msg.Extra.String, &obj); err == nil {
				extra = obj
			}
		}
	}

	var mediaURL string
	var msgExtra = &fileExtra{}
	if msg.ContentType.String == models.MessageFileContentType || msg.ContentType.String == models.MessagePictureContentType {
		mediaURL = entMediaURL(msg.EntID)
		if err = common.Unmarshal(msg.Extra.String, &msgExtra); err == nil {
			extra = msgExtra
			mediaURL = mediaURL + msgExtra.Filename
		}
	}

	result := &adapter.Message{
		Action:         "message",
		ID:             msg.ID,
		Agent:          agent,
		FromType:       msg.FromType.String,
		MediaURL:       mediaURL,
		Content:        msg.Content.String,
		ContentType:    msg.ContentType.String,
		CreatedOn:      common.ConvertUTCToTimeString(msg.CreatedAt),
		AgentID:        msg.AgentID,
		ConversationID: msg.ConversationID,
		ContentRobot:   nil,
		EnterpriseID:   msg.EntID,
		Extra:          extra,
		QuestionID:     "",
		ReadOn:         nil,
		SubType:        "",
		TraceStart:     -1,
		TrackID:        msg.TraceID,
		Type:           "message",
	}

	if msg.MsgType.String == models.MessageMsgInternalType {
		result.Type = models.MessageMsgInternalType
	}

	return result
}
