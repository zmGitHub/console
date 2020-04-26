package events

const (
	InviteEvaluationAction = "invite_evaluation"
	LeaveMessageUpdate     = "leave_message_update"
	QueueingAddAction      = "queueing_add"
	QueueingRemoveAction   = "queueing_remove"
	AgentChangeAtrr        = "agent_change_attr"
)

type Event struct {
	ID            string  `json:"id"`
	Action        string  `json:"action"`
	AgentID       string  `json:"agent_id"`
	AgentNickname string  `json:"agent_nickname"`
	CreatedOn     string  `json:"created_on"`
	EnterpriseID  string  `json:"enterprise_id"`
	RealName      string  `json:"realname"`
	Source        string  `json:"source"`
	TargetID      string  `json:"target_id"`
	TargetKind    string  `json:"target_kind"`
	TraceStart    float64 `json:"trace_start"`
	TrackID       string  `json:"track_id"`
}
