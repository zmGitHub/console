package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AutomaticMessage struct {
	ChannelType  string `json:"channel_type"`  // channel_type
	MsgType      string `json:"msg_type"`      // msg_type
	MsgContent   string `json:"msg_content"`   // msg_content
	AfterSeconds int    `json:"after_seconds"` // after_seconds
	Enabled      bool   `json:"enabled"`       // enabled
}

func (msg *AutomaticMessage) validate() error {
	if _, ok := models.AutoMessageChannelTypeMap[msg.ChannelType]; !ok {
		return fmt.Errorf("not supported channel_type")
	}

	if _, ok := models.AutoMessageMsgTypeMap[msg.MsgType]; !ok {
		return fmt.Errorf("not supported msg_type")
	}

	if msg.AfterSeconds < 0 {
		return fmt.Errorf("invalid after_seconds")
	}

	return nil
}

// GET /admin/api/v1/automessages
func (s *IMService) GetEntAutoMessages(ctx echo.Context) (err error) {
	entID := ctx.Get(middleware.AgentEntIDKey).(string)
	msgs, err := models.AutomaticMessagesByEntID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return jsonResponse(ctx, &Resp{Code: 0, Body: nil})
		}

		return dbErrResp(ctx, err.Error())
	}
	return jsonResponse(ctx, &Resp{Code: 0, Body: msgs})
}

// CreateAutoMessage create automessages
// POST /admin/api/v1/automessages
func (s *IMService) CreateAutoMessage(ctx echo.Context) (err error) {
	msg := &AutomaticMessage{}
	if err = ctx.Bind(msg); err != nil {
		return
	}

	if err = msg.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	autoMsg := &models.AutomaticMessage{
		ID:           common.GenUniqueID(),
		ChannelType:  msg.ChannelType,
		MsgType:      msg.MsgType,
		MsgContent:   msg.MsgContent,
		AfterSeconds: msg.AfterSeconds,
		Enabled:      msg.Enabled,
		CreatedAt:    time.Now().UTC(),
	}

	autoMsg.EntID = ctx.Get(middleware.AgentEntIDKey).(string)
	if err := autoMsg.Save(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: autoMsg})
}

// UpdateAutoMessage ...
// PUT /admin/api/v1/automessages/:msg_id
func (s *IMService) UpdateAutoMessage(ctx echo.Context) (err error) {
	msgID := ctx.Param("msg_id")
	if msgID == "" {
		return invalidParameterResp(ctx, "msg_id is empty")
	}

	msg := &AutomaticMessage{}
	if err = ctx.Bind(msg); err != nil {
		return
	}

	if err = msg.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	autoMsg, err := models.AutomaticMessageByID(db.Mysql, msgID)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	autoMsg.MsgContent = msg.MsgContent
	autoMsg.AfterSeconds = msg.AfterSeconds
	autoMsg.Enabled = msg.Enabled
	if err = autoMsg.Update(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0, Body: autoMsg})
}

// DELETE /admin/api/v1/automessages/:msg_id
func (s *IMService) DeleteAutoMessage(ctx echo.Context) (err error) {
	msgID := ctx.Param("msg_id")
	if msgID == "" {
		return invalidParameterResp(ctx, "msg_id is empty")
	}

	msg := &models.AutomaticMessage{ID: msgID}
	if err = msg.Delete(db.Mysql); err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, &Resp{Code: 0})
}
