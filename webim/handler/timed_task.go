package handler

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/external/timedevent"
	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var (
	defaultNoReplyMsgCountdownInSecond = 20
	defaultEndConvMinCountdownInMinute = 3
	defaultEndConvMaxCountdownInMinute = 100

	defaultOfflineEndConvMinCountdownInSec = 10
	defaultOfflineEndConvMaxCountdownInSec = 300

	defaultPingInterval int64 = 25
)

// POST /api/v1/system/send_ent_message
func (s *IMService) SendEntMessage(ctx echo.Context) (err error) {
	req := &timedevent.AddSendingEntMessageReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	agents, err := models.AgentsByConversationID(db.Mysql, req.ConversationID)
	if err != nil {
		agents = []string{req.AgentID}
	}

	sendMsgReq := &sendMsgReq{
		EntID:       req.EntID,
		TraceID:     req.TraceID,
		Creator:     req.AgentID,
		Agents:      agents,
		ConvID:      req.ConversationID,
		FromType:    models.MessageFromAgentType,
		Content:     req.MsgContent,
		ContentType: models.MessageTextContentType,
		Internal:    false,
		System:      true,
	}
	if req.ContentType != "" {
		sendMsgReq.ContentType = req.ContentType
	}

	if _, err := s.sendMsg(sendMsgReq); err != nil {
		log.Logger.Warnf("sendMsg error: %v", err)
		return jsonResponse(ctx, err)
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

// SendNoRespMessage if visitor or agent no resp in some duration,will send message to channel
func (s *IMService) SendNoRespMessage(ctx echo.Context) (err error) {
	req := &timedevent.AddSendingNoRespMessageReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.SenderType != "client" && req.SenderType != "agent" {
		return invalidParameterResp(ctx, "invalid sender type")
	}

	convEnd, err := models.IsConversationEnd(db.Mysql, req.ConversationID)
	if err != nil {
		log.Logger.Warnf("models.IsConversationEnd error: %v", err)
		return dbErrResp(ctx, err.Error())
	}

	if convEnd {
		log.Logger.Debug("SendNoRespMessage: conversation is end")
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	var lastMsgCreateTime time.Time
	if req.SenderType == "client" {
		visitor, err := models.VisitorByID(db.Mysql, req.Sender)
		if err != nil {
			return jsonResponse(ctx, &Resp{Code: 0, Body: err})
		}

		lastMsgCreateTime, err = models.VisitorLastMessage(db.Mysql, visitor.TraceID)
	} else {
		lastMsgCreateTime, err = models.LastMessageByFromType(db.Mysql, req.ConversationID, models.MessageFromAgentType)
	}
	if err != nil {
		return jsonResponse(ctx, &Resp{Code: 0, Body: err})
	}

	if lastMsgCreateTime.IsZero() {
		log.Logger.Infof("lastMsgCreateTime: %v", lastMsgCreateTime)
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	agents, err := models.AgentsByConversationID(db.Mysql, req.ConversationID)
	if err != nil {
		agents = []string{req.AgentID}
	}

	now := time.Now().UTC()
	duration := int64(now.Sub(lastMsgCreateTime.UTC()).Seconds())
	afterSec := int64(req.AfterSeconds)
	if duration < afterSec {
		log.Logger.Infof("duration: %d, afterSec: %d", duration, afterSec)
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	sendMsgReq := &sendMsgReq{
		EntID:       req.EntID,
		TraceID:     req.TraceID,
		Creator:     req.AgentID,
		Agents:      agents,
		ConvID:      req.ConversationID,
		FromType:    models.MessageFromAgentType,
		Content:     req.MsgContent,
		ContentType: models.MessageTextContentType,
		Internal:    false,
		System:      true,
	}

	if _, err := s.sendMsg(sendMsgReq); err != nil {
		log.Logger.Warnf("SendNoRespMessage error: %v", err)
	}
	return jsonResponse(ctx, &Resp{Code: 0})
}

// SysEndConversation
// POST /api/v1/system/end_conversation
func (s *IMService) SysEndConversation(ctx echo.Context) (err error) {
	req := &timedevent.AddEndingConversationTaskReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.AgentID == "" || req.ConversationID == "" || req.EntID == "" {
		return invalidParameterResp(ctx, "invalid request params")
	}

	createdAt, endedAt, err := models.ConversationCreatedAtByID(db.Mysql, req.ConversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return jsonResponse(ctx, &Resp{Code: 0})
		}

		return dbErrResp(ctx, err.Error())
	}
	if endedAt.Valid {
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	lastMsg, dbErr := models.LastMessageByConversationID(db.Mysql, req.ConversationID)
	if dbErr != nil {
		log.Logger.Warnf("LastMessageByConversationID error: %v", dbErr)
		return dbErrResp(ctx, dbErr.Error())
	}

	var duration int64
	now := time.Now().UTC()
	if lastMsg == nil {
		duration = int64(now.Sub(createdAt.UTC()).Seconds())
	} else {
		duration = int64(now.Sub(lastMsg.CreatedAt.UTC()).Seconds())
	}

	if duration >= int64(req.AfterSeconds) {
		if errMsg := s.endConversation("system", req.ConversationID); errMsg != nil {
			log.Logger.Warnf("endConversation error: %v", errMsg)
		}
		go s.sendEndConvEvent(EndConvTimeout, req.AgentID, req.ConversationID)
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

// POST /system/offline_end_conversation
func (s *IMService) SysOfflineEndConversation(ctx echo.Context) (err error) {
	req := &timedevent.AddOfflineTaskReq{}
	if err = ctx.Bind(req); err != nil {
		return
	}

	if req.AgentID == "" || req.ConversationID == "" || req.EntID == "" || req.TraceID == "" {
		return invalidParameterResp(ctx, "invalid request params")
	}

	_, endedAt, err := models.ConversationCreatedAtByID(db.Mysql, req.ConversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return jsonResponse(ctx, &Resp{Code: 0})
		}

		return dbErrResp(ctx, err.Error())
	}
	if endedAt.Valid {
		return jsonResponse(ctx, &Resp{Code: 0})
	}

	onlineVisitor, err := models.OnlineVisitorByEntIDTraceID(db.Mysql, req.EntID, req.TraceID)
	if err != nil && err != sql.ErrNoRows {
		log.Logger.Warnf("[OnlineVisitorByEntIDTraceID] error: %v", err)
		return dbErrResp(ctx, err.Error())
	}

	var duration float64
	if onlineVisitor != nil {
		duration = math.Abs(time.Now().UTC().Sub(onlineVisitor.UpdatedAt.UTC()).Seconds())
	} else {
		duration = math.MaxFloat64
	}

	if int64(duration) >= (req.AfterSeconds + defaultPingInterval) {
		if errMsg := s.endConversation("system", req.ConversationID); errMsg != nil {
			log.Logger.Warnf("endConversation error: %v", errMsg)
		}
		go s.sendEndConvEvent(EndConvOffline, req.AgentID, req.ConversationID)
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}

func (s *IMService) AddTimedTasks(entID, trackID, agentID, convID, visitorID string) {
	configs, errMsg := s.getEnterpriseConfigs(entID)
	if errMsg != nil {
		log.Logger.Warnf("getEnterpriseConfigs error: %+v", *errMsg)
		return
	}

	if configs == nil {
		return
	}

	s.SendAutoMessages(entID, trackID, agentID, convID, configs)
	s.SendNoReplyMessages(entID, trackID, agentID, convID, visitorID, configs)
	s.AddEndConversationTask(entID, agentID, convID, trackID, configs)
}

// SendAutoMessages send promotion msg & welcome msg after scheduler
func (s *IMService) SendAutoMessages(entID, trackID, agentID, convID string, configs *adapter.EnterpriseConfigs) *ErrMsg {
	promotions, err := models.PromotionMsgsByEnterpriseID(db.Mysql, entID)
	if err != nil {
		log.Logger.Warnf("PromotionMsgsByEnterpriseID error: %v", err)
	}

	if promotions != nil && len(promotions) > 0 {
		go s.asyncSendPromotionMessage(entID, trackID, agentID, convID, promotions)
	}

	if configs.WelcomeMsgSettings != nil {
		go s.asyncSendWelcomeMessage(entID, trackID, agentID, convID, configs.WelcomeMsgSettings)
	}

	return nil
}

func (s *IMService) asyncSendPromotionMessage(entID, trackID, agentID, convID string, msgs []*models.PromotionMsg) {
	if len(msgs) == 0 {
		return
	}

	agents, err := models.AgentsByConversationID(db.Mysql, convID)
	if err != nil {
		agents = []string{agentID}
	}

	sendMsgReq := &sendMsgReq{
		EntID:       entID,
		TraceID:     trackID,
		Creator:     agentID,
		Agents:      agents,
		ConvID:      convID,
		FromType:    models.MessageFromAgentType,
		Content:     "",
		ContentType: models.MessageRichTextType,
		Internal:    false,
		System:      true,
	}

	var instantMsgs, delayedMsgs []*models.PromotionMsg
	for _, msg := range msgs {
		if msg != nil && msg.Enabled && msg.Content.String != "" {
			if msg.Countdown > 2 {
				delayedMsgs = append(delayedMsgs, msg)
				continue
			}

			instantMsgs = append(instantMsgs, msg)
		}
	}

	if len(instantMsgs) > 0 {
		for _, msg := range instantMsgs {
			sendMsgReq.Content = msg.Content.String
			if _, err := s.sendMsg(sendMsgReq); err != nil {
				log.Logger.Warnf("asyncSendPromotionMessage error: %v", err)
			}
		}
	}

	if len(delayedMsgs) > 0 {
		task := &timedevent.AddSendingEntMessageReq{
			EntID:          entID,
			TraceID:        trackID,
			AgentID:        agentID,
			ConversationID: convID,
			ContentType:    models.MessageRichTextType,
		}

		for _, msg := range delayedMsgs {
			task.AfterSeconds = int64(msg.Countdown)
			task.MsgContent = msg.Content.String

			if err := s.taskHandler.AddSendingEntMessage(task); err != nil {
				log.Logger.Warnf("asyncSendPromotionMessage error: %v", err)
			}
		}
	}
}

func (s *IMService) asyncSendWelcomeMessage(entID, trackID, agentID, convID string, welcomeMsg *adapter.WelcomeMsgSettings) {
	if welcomeMsg.Web != nil && welcomeMsg.Web.Status == "open" {
		time.Sleep(500 * time.Millisecond)

		sendMsgReq := &sendMsgReq{
			EntID:       entID,
			TraceID:     trackID,
			Creator:     agentID,
			Agents:      []string{agentID},
			ConvID:      convID,
			FromType:    models.MessageFromAgentType,
			Content:     welcomeMsg.Web.Content,
			ContentType: models.MessageTextContentType,
			Internal:    false,
			System:      true,
		}

		if _, err := s.sendMsg(sendMsgReq); err != nil {
			log.Logger.Warnf("asyncSendWelcomeMessage error: %v", err)
		}
	}
}

func (s *IMService) SendNoReplyMessages(entID, trackID, agentID, convID, visitorID string, configs *adapter.EnterpriseConfigs) {
	if configs == nil {
		return
	}

	agentNoReplyMsg, clientNoReplyMsg := configs.AutoReplyMsgSettings, configs.ClientWalkingAutoMsg
	task := &timedevent.AddSendingNoRespMessageReq{
		EntID:          entID,
		AgentID:        agentID,
		TraceID:        trackID,
		ConversationID: convID,
	}

	if agentNoReplyMsg != nil {
		web := agentNoReplyMsg.Web
		if web != nil && web.Status == "open" {
			countDown := web.CountDown
			if countDown < defaultNoReplyMsgCountdownInSecond {
				countDown = defaultNoReplyMsgCountdownInSecond
			}

			task.Sender = agentID
			task.SenderType = "agent"
			task.AfterSeconds = countDown
			task.MsgContent = web.Content
			if err := s.taskHandler.AddSendingNoRespMessage(task); err != nil {
				log.Logger.Warnf("taskHandler.AddSendingNoRespMessage error: %v", err)
			}
		}
	}

	if clientNoReplyMsg != nil {
		web := clientNoReplyMsg.Web
		if web != nil && web.Status == "open" {
			countDown := web.CountDown
			if countDown < defaultNoReplyMsgCountdownInSecond {
				countDown = defaultNoReplyMsgCountdownInSecond
			}

			task.Sender = visitorID
			task.SenderType = "client"
			task.AfterSeconds = countDown
			task.MsgContent = web.Content
			if err := s.taskHandler.AddSendingNoRespMessage(task); err != nil {
				log.Logger.Warnf("taskHandler.AddSendingNoRespMessage error: %v", err)
			}
		}
	}
}

func (s *IMService) AddEndConversationTask(entID, agentID, convID, traceID string, configs *adapter.EnterpriseConfigs) {
	if configs == nil {
		return
	}

	endConv := configs.EndConvExpireConfig
	if endConv == nil {
		return
	}

	web := endConv.Web
	if web != nil && web.NoMsgEnd != -1 {
		countDown := web.NoMsgEnd
		if countDown < defaultEndConvMinCountdownInMinute || countDown > defaultEndConvMaxCountdownInMinute {
			countDown = defaultEndConvMinCountdownInMinute
		}

		task := &timedevent.AddEndingConversationTaskReq{
			EntID:          entID,
			AgentID:        agentID,
			ConversationID: convID,
			AfterSeconds:   countDown * 60,
		}
		if err := s.taskHandler.AddEndingConversation(task); err != nil {
			log.Logger.Warnf("AddEndConversationTask error: %v", err)
		}
	}
	if web != nil && web.OfflineEnd != -1 {
		countDown := web.OfflineEnd
		if countDown < defaultOfflineEndConvMinCountdownInSec || countDown > defaultOfflineEndConvMaxCountdownInSec {
			countDown = defaultOfflineEndConvMinCountdownInSec
		}

		task := &timedevent.AddOfflineTaskReq{
			EntID:          entID,
			AgentID:        agentID,
			TraceID:        traceID,
			ConversationID: convID,
			AfterSeconds:   int64(countDown),
		}
		if err := s.taskHandler.AddOfflineEndConversation(task); err != nil {
			log.Logger.Warnf("AddOfflineEndConversation error: %v", err)
		}
	}
}

func (s *IMService) DeleteConversationTasks(convID string) {
	deleteJobReq := &timedevent.DeleteTaskReq{
		TaskNames: []string{
			fmt.Sprintf(timedevent.SendEntMessageTaskTemplate, convID),
			fmt.Sprintf(timedevent.EndConversationTaskTemplate, convID),
			fmt.Sprintf(timedevent.OfflineEndConversationTaskTemplate, convID),
			fmt.Sprintf(timedevent.SendNoRespMessageTaskTemplate, convID, models.MessageFromAgentType),
			fmt.Sprintf(timedevent.SendNoRespMessageTaskTemplate, convID, models.MessageFromVisitorType),
		},
	}
	if err := s.taskHandler.DeleteJob(deleteJobReq); err != nil {
		log.Logger.Warnf("DeleteConversationTasks error: %v", err)
	}
}
