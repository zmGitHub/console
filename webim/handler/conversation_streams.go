package handler

import (
	"sort"

	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

func (s *IMService) getConversationStreams(conv *models.Conversation, msgs []*models.Message, agents []*adapter.Agent, visitID string, visit *models.Visit, visitor *models.Visitor, tags []*models.VisitorTagRelation) []interface{} {
	var elements = []interface{}{}

	initConv := s.convertToNewConvAction(adapter.ModelConversationToConversation(conv, visitID, visit, visitor, tags))
	elements = append(elements, initConv)

	modelMessages := []*models.Message{}
	for _, msg := range msgs {
		if msg.ConversationID == conv.ID {
			modelMessages = append(modelMessages, msg)
		}
	}

	messages := convertModelMessagesToAdapterMessages(modelMessages, agents)
	for _, msg := range messages {
		elements = append(elements, msg)
	}

	return elements
}

func (s *IMService) getConversationVisitInfo(entID, trackID string) (visit *models.Visit, visitor *models.Visitor, visitID string, err error) {
	var visits []*models.Visit
	visits, err = models.VisitsByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		return
	}

	visitor, err = models.VisitorByEntIDTraceID(db.Mysql, entID, trackID)
	if err != nil {
		return
	}

	if len(visits) > 0 {
		visit = visits[0]
	}

	if visit != nil {
		visitID = visit.ID
	}
	return
}

func (s *IMService) findConversations(convs []*models.Conversation, limit int, req *ConvStreamsReq) (results []*models.Conversation) {
	if len(convs) == 0 {
		return nil
	}

	if len(convs) == 1 {
		if req.Include == 1 {
			return convs
		}

		return nil
	}

	sort.SliceStable(convs, func(i, j int) bool {
		return convs[i].CreatedAt.Before(convs[j].CreatedAt)
	})

	index := -1
	for i, conv := range convs {
		if conv.ID == req.ConvID {
			index = i
			break
		}
	}

	if req.Type == "earlier" {
		var total []*models.Conversation
		switch index {
		case 0:
			total = convs[:1]
		case len(convs) - 1:
			total = convs
		default:
			total = convs[:index+1]
		}

		if limit >= len(total) {
			results = total
		} else {
			results = total[len(total)-limit:]
		}

		if req.Include == 0 {
			return results[:len(results)-1]
		}

		return results
	}

	var total []*models.Conversation
	switch index {
	case 0:
		total = convs
	case len(convs) - 1:
		total = convs[len(convs)-1:]
	default:
		total = convs[index:]
	}

	if limit >= len(total) {
		results = total
	} else {
		results = total[:limit]
	}

	if req.Include == 0 {
		return results[1:]
	}

	return results
}

func (s *IMService) getVisitorConversationStreams(trackID string, req *ConvStreamsReq) (streams *adapter.ConversationStreams, errMsg *ErrMsg) {
	var limit = 1
	if req.Limit > 1 {
		limit = req.Limit
	}

	convs, dbErr := models.ConversationsByTraceID(db.Mysql, trackID, 0, 100, false)
	if dbErr != nil {
		return nil, &ErrMsg{Message: dbErr.Error()}
	}

	result := &adapter.ConversationStreams{ConvIDs: []string{}, Total: 0, Streams: []interface{}{}}

	resultConvs := s.findConversations(convs, limit, req)
	if resultConvs == nil || len(resultConvs) == 0 {
		return result, nil
	}

	entID := convs[0].EntID
	visit, visitor, visitID, err := s.getConversationVisitInfo(entID, trackID)
	if err != nil {
		return nil, &ErrMsg{Message: err.Error()}
	}

	tags, err := models.VisitorTagRelationsByVisitors(db.Mysql, []*models.Visitor{visitor})
	if err != nil {
		return nil, &ErrMsg{Message: err.Error()}
	}

	var convIDs = []string{}
	for _, v := range resultConvs {
		convIDs = append(convIDs, v.ID)
	}

	msgs, err := models.MessagesByConversationIDs(db.Mysql, convIDs)
	if err != nil {
		return nil, &ErrMsg{Message: err.Error()}
	}

	agents, err := models.AgentsByMsgs(db.Mysql, msgs)
	if err != nil {
		return nil, &ErrMsg{Message: err.Error()}
	}

	agentGroups, err := models.PermsRangeGroupIDsByAgents(db.Mysql, agents)
	if err != nil {
		return nil, &ErrMsg{Message: err.Error()}
	}

	adapterAgents := adapter.ConvertAgentsToAdapterAgentsV1(agents, agentGroups)

	result.ConvIDs = convIDs
	result.Total = len(convIDs)

	for _, conv := range resultConvs {
		result.Streams = append(result.Streams, s.getConversationStreams(
			conv,
			msgs,
			adapterAgents,
			visitID,
			visit,
			visitor,
			tags,
		))
	}

	return result, nil
}
