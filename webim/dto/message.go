package dto

import "time"

type ContentRobot struct {
	RichText string `json:"rich_text,omitempty"`
	Text     string `json:"text,omitempty"`
	Type     string `json:"type"`
	Items    []struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"items,omitempty"`
}

type Message struct {
	Action         string          `json:"action"`
	ID             string          `json:"id"`
	Agent          *Agent          `json:"agent"`
	AgentID        string          `json:"agent_id"`
	Content        string          `json:"content"`
	ContentType    string          `json:"content_type"`
	ConversationID string          `json:"conversation_id"`
	ContentRobot   []*ContentRobot `json:"content_robot"`
	CreatedOn      *string         `json:"created_on"`
	EnterpriseID   string          `json:"enterprise_id"`
	Extra          interface{}     `json:"extra"`
	FromType       string          `json:"from_type"`
	MediaURL       string          `json:"media_url"`
	QuestionID     string          `json:"question_id"`
	ReadOn         *time.Time      `json:"read_on"`
	SubType        string          `json:"sub_type"` // menu
	TraceStart     int64           `json:"trace_start"`
	TrackID        string          `json:"track_id"`
	Type           string          `json:"type"` // message
}
