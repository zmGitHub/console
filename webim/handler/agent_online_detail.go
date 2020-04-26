package handler

import "github.com/labstack/echo/v4"

// {"begin":1556640000000,"end":1556726399999,"group_id":null,"browser_id":"agent1556722413695"}
type AgentOnlineDetailReq struct {
	Begin   int64   `json:"begin"`
	End     int64   `json:"end"`
	GroupID *string `json:"group_id"`
}

type statusDetail struct {
	Date      int    `json:"date"` // 1556088079
	Duration  int    `json:"duration"`
	IP        string `json:"ip"`
	Note      string `json:"note"`
	OriStatus string `json:"ori_status"`
	Platform  string `json:"platform"`
	Status    string `json:"status"`
}

type agentStatusCnt struct {
	AgentID       string          `json:"agent_id"`
	Detail        []*statusDetail `json:"detail"`
	OnlineOffduty int             `json:"online_offduty"`
	OnlineOnduty  int             `json:"online_onduty"`
	OnlineTotal   int             `json:"online_total"`
}

type AgentOnlineDetailResp struct {
	Data []*agentStatusCnt `json:"data"`
}

// POST /api/analytics/agent_online_detail
func (s *IMService) AgentOnlineDetail(ctx echo.Context) error {
	req := &AgentOnlineDetailReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	result := &AgentOnlineDetailResp{Data: []*agentStatusCnt{}}
	return jsonResponse(ctx, result)
}
