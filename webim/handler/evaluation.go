package handler

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/external/elasticsearch"
	"bitbucket.org/forfd/custm-chat/webim/middleware"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type evaluationResp struct {
	Level int `json:"level"`
}

type evaluation struct {
	AgentID string `json:"agent_id"`
	Level   int    `json:"level"`
	Content string `json:"content"`
}

type evalResp struct {
	Success bool `json:"success"`
}

func (eval *evaluation) validate() error {
	if eval.AgentID == "" {
		return fmt.Errorf("agent_id is invalid")
	}

	if !models.EvalLevels.Contains(eval.Level) {
		return fmt.Errorf("invalid evaluation level")
	}

	return nil
}

// CreateEvaluation agent eval
// POST /api/v1/enterprises/:ent_id/conversations/:conversation_id/evaluations
func (s *IMService) CreateEvaluation(ctx echo.Context) (err error) {
	entID := ctx.Param("ent_id")
	convID := ctx.Param("conversation_id")
	if entID == "" {
		return invalidParameterResp(ctx, "ent_id is empty")
	}

	eval := &evaluation{}
	if err = ctx.Bind(eval); err != nil {
		return
	}
	if err = eval.validate(); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	err = models.CreateConversationEvalV1(db.Mysql, convID, eval.Content, eval.Level)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	return jsonResponse(ctx, eval)
}

// POST /conversation/:conversation_id/evaluation
// resp: {"success":true}
func (s *IMService) EvaluateConversation(ctx echo.Context) (err error) {
	convID := ctx.Param("conversation_id")
	if convID == "" {
		return invalidParameterResp(ctx, "conversation_id is empty")
	}

	eval := &evaluation{}
	if err = ctx.Bind(eval); err != nil {
		return
	}
	if eval.Level < 0 || eval.Level > 2 {
		return invalidParameterResp(ctx, "invalid eval level")
	}

	err = models.CreateConversationEvalV1(db.Mysql, convID, eval.Content, eval.Level)
	if err != nil {
		return dbErrResp(ctx, err.Error())
	}

	go elasticsearch.UpdateConversationEvalLevel(elasticsearch.ESClient, convID, &eval.Level, eval.Content)

	go s.sendCreateEvalEvent(ClientEvaluation, convID, &evaluationEventBody{
		Content: eval.Content,
		EvaType: "create",
		Level:   eval.Level,
	})

	return jsonResponse(ctx, &evalResp{Success: true})
}

// POST /api/conversation/:conversation_id/invite_evaluation
func (s *IMService) InviteEval(ctx echo.Context) (err error) {
	convID := ctx.Param("conversation_id")
	if convID == "" {
		return invalidParameterResp(ctx, "conversation_id invalid")
	}

	agentID := ctx.Get(middleware.AgentIDKey).(string)
	agent, err := models.AgentByID(db.Mysql, agentID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	conv, err := models.ConversationByID(db.Mysql, convID)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	event := &events.InviteEval{
		Event: &events.Event{
			ID:            common.GenUniqueID(),
			Action:        events.InviteEvaluationAction,
			AgentID:       agentID,
			AgentNickname: agent.NickName,
			RealName:      agent.RealName,
			CreatedOn:     *common.ConvertUTCToTimeString(time.Now().UTC()),
			EnterpriseID:  agent.EntID,
			TargetID:      conv.ID,
			TargetKind:    "conv",
			TraceStart:    float64(conv.CreatedAt.Unix()),
			TrackID:       conv.TraceID,
		},
		Body: struct{}{},
	}

	content, err := common.Marshal(event)
	if err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	var agents = []string{conv.AgentID}
	if agentID != conv.AgentID {
		agents = append(agents, agentID)
	}

	s.sendMessageToMultiAgents(agents, content)
	sendMessageToVisitor(s.imCli, conv.TraceID, conv.EntID, content)

	msgV := &msg{
		CreatedOn: event.CreatedOn,
		ID:        event.ID,
	}

	return jsonResponse(ctx, &InviteEvalResp{Msg: msgV, Success: true})
}
