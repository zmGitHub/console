package dto

import "time"

type Conversation struct {
	AgentEffectiveMsgNum  int         `json:"agent_effective_msg_num"`
	AgentID               string      `json:"agent_id"`
	AgentMsgNum           int         `json:"agent_msg_num"`
	AgentType             string      `json:"agent_type"`
	Assignee              string      `json:"assignee"`
	ClientFirstSendTime   *time.Time  `json:"client_first_send_time"`
	ClientLastSendTime    *time.Time  `json:"client_last_send_time"`
	ClientMsgNum          int         `json:"client_msg_num"`
	Clues                 interface{} `json:"clues"`
	ConverseDuration      int         `json:"converse_duration"`
	CreatedOn             time.Time   `json:"created_on"`
	EndedBy               string      `json:"ended_by"`
	EndedOn               *time.Time  `json:"ended_on"`
	EnterpriseID          string      `json:"enterprise_id"`
	EvaContent            string      `json:"eva_content"`
	EvaLevel              int         `json:"eva_level"`
	FirstMsgCreatedOn     *time.Time  `json:"first_msg_created_on"`
	FirstResponseWaitTime int64       `json:"first_response_wait_time"`
	HasSummary            bool        `json:"has_summary"`
	ID                    string      `json:"id"`
	IsClientOnline        bool        `json:"is_client_online"`
	LastMsgContent        string      `json:"last_msg_content"`
	LastMsgCreatedOn      *time.Time  `json:"last_msg_created_on"`
	LastUpdated           *time.Time  `json:"last_updated"`
	Messages              []*Message  `json:"messages"`
	MsgNum                int         `json:"msg_num"`
	QualityGrade          string      `json:"quality_grade"`
	Tags                  []string    `json:"tags,omitempty"`
	Title                 string      `json:"title"`
	TrackID               string      `json:"track_id"`
	URL                   string      `json:"url"`
	VisitID               string      `json:"visit_id"`
	VisitInfo             *VisitInfo  `json:"visit_info,omitempty"`
}
