package timedevent

import (
	"github.com/parnurzeal/gorequest"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	SendEntMessage                     = "send_message"
	SendEntMessageTaskTemplate         = "%s_send_auto_message" // {conv_id}_send_auto_message
	SendNoRespMessage                  = "send_no_resp_message"
	SendNoRespMessageTaskTemplate      = "%s_%s_send_no_resp_message" // "{conv_id}_{sender_type}_send_no_resp_message"
	EndConversation                    = "end_conversation"
	EndConversationTaskTemplate        = "%s_end_conversation" // {conv_id}_end_conversation
	OfflineEndConversation             = "offline_end_conversation"
	OfflineEndConversationTaskTemplate = "%s_offline_end_conversation" // // {conv_id}_offline_end_conversation

	DeleteJob = "delete_job"
)

type Handler struct {
	Host      string
	EndPoints map[string]string
}

type AddSendingEntMessageReq struct {
	EntID          string `json:"ent_id"`
	TraceID        string `json:"trace_id"`
	AgentID        string `json:"agent_id"`
	ConversationID string `json:"conversation_id"`
	MsgContent     string `json:"msg_content"`
	ContentType    string `json:"content_type"`
	AfterSeconds   int64  `json:"after_seconds"`
}

type AddSendingNoRespMessageReq struct {
	EntID          string `json:"ent_id"`
	AgentID        string `json:"agent_id"`
	TraceID        string `json:"trace_id"`
	ConversationID string `json:"conversation_id"`
	Sender         string `json:"sender"`
	SenderType     string `json:"sender_type"`
	MsgContent     string `json:"msg_content"`
	AfterSeconds   int    `json:"after_seconds"`
}

type AddEndingConversationTaskReq struct {
	EntID          string `json:"ent_id"`
	AgentID        string `json:"agent_id"`
	ConversationID string `json:"conversation_id"`
	AfterSeconds   int    `json:"after_seconds"`
}

type AddOfflineTaskReq struct {
	EntID          string `json:"ent_id"`
	AgentID        string `json:"agent_id"`
	TraceID        string `json:"trace_id"`
	ConversationID string `json:"conversation_id"`
	AfterSeconds   int64  `json:"after_seconds"`
}

type DeleteTaskReq struct {
	TaskNames []string `json:"task_names"`
}

func InitHandler(config *conf.TaskHandlerConfig) *Handler {
	return &Handler{
		Host: config.Host,
		EndPoints: map[string]string{
			SendEntMessage:         config.SendEntMessageEndPoint,
			SendNoRespMessage:      config.SendNoRespMessageEndPoint,
			EndConversation:        config.EndConversationEndPoint,
			OfflineEndConversation: config.OfflineEndConversationEndPoint,
			DeleteJob:              config.DeleteJobEndPoint,
		},
	}
}

func (h *Handler) AddSendingEntMessage(req *AddSendingEntMessageReq) error {
	content, err := common.Marshal(req)
	if err != nil {
		return err
	}

	request := gorequest.New()
	_, _, errs := request.Post(h.Host + h.EndPoints[SendEntMessage]).
		Send(content).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (h *Handler) AddEndingConversation(req *AddEndingConversationTaskReq) error {
	request := gorequest.New()
	reqContent, err := common.Marshal(req)
	if err != nil {
		log.Logger.Errorf("common.Marshal AddEndingConversationTaskReq error: %v", err)
		return err
	}

	_, _, errs := request.Post(h.Host + h.EndPoints[EndConversation]).
		Send(reqContent).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (h *Handler) AddSendingNoRespMessage(req *AddSendingNoRespMessageReq) error {
	request := gorequest.New()
	reqContent, err := common.Marshal(req)
	if err != nil {
		log.Logger.Errorf("common.Marshal req error: %v", err)
		return err
	}

	_, _, errs := request.Post(h.Host + h.EndPoints[SendNoRespMessage]).
		Send(reqContent).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (h *Handler) AddOfflineEndConversation(req *AddOfflineTaskReq) error {
	return h.addTask(req, OfflineEndConversation)
}

func (h *Handler) DeleteJob(req *DeleteTaskReq) error {
	request := gorequest.New()
	reqContent, err := common.Marshal(req)
	if err != nil {
		log.Logger.Errorf("common.Marshal DeleteJob req error: %v", err)
		return err
	}

	_, _, errs := request.Post(h.Host + h.EndPoints[DeleteJob]).
		Send(reqContent).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (h *Handler) addTask(task interface{}, endpoint string) error {
	request := gorequest.New()
	reqContent, err := common.Marshal(task)
	if err != nil {
		log.Logger.Errorf("common.Marshal req error: %v", err)
		return err
	}

	_, _, errs := request.Post(h.Host + h.EndPoints[endpoint]).
		Send(reqContent).
		End()

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
