package adapter

import (
	"time"
)

type SchedulerRequest struct {
	EntID            string  `json:"ent_id"`
	TrackID          string  `json:"track_id"`
	VisitID          string  `json:"visit_id"`
	AgentToken       *string `json:"agent_token"`
	GroupToken       *string `json:"group_token"`
	URL              string  `json:"url"`
	Title            string  `json:"title"`
	Queueing         bool    `json:"queueing"`
	FromType         string  `json:"from_type"`
	ConvInitiateType string  `json:"conv_initiate_type"`
	ReferrerURL      string  `json:"referrer_url"`
}

type ContentRobot struct {
	RichText string `json:"rich_text,omitempty"`
	Text     string `json:"text,omitempty"`
	Type     string `json:"type"`
	Items    []struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"items,omitempty"`
}

type MessageExtra struct {
	ID           string          `json:"id"`
	EnterpriseID string          `json:"enterprise_id"`
	Content      string          `json:"content"`
	ContentRobot []*ContentRobot `json:"content_robot"`
	SubType      string          `json:"sub_type"`
	Summary      string          `json:"summary"`
	Thumbnail    string          `json:"thumbnail"`
	CreatedOn    string          `json:"created_on"`
	UpdatedOn    string          `json:"updated_on"`
}

type FileMessageExtra struct {
	Cancel   bool   `json:"cancel"`
	ExpireAt string `json:"expire_at"`
	FileName string `json:"file_name"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
}

type Message struct {
	ID             string          `json:"id"`
	Action         string          `json:"action"`
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
	Type           string          `json:"type"` // message / internal
}

type SchedulerResponse struct {
	Agent         *Agent          `json:"agent"`
	Conv          *Conversation   `json:"conv"`
	ConvNewCreate bool            `json:"conv_new_create"`
	Position      *int            `json:"position"`
	Ent           *EnterpriseResp `json:"ent"`
	Messages      []*Message      `json:"messages"`
	ReserveToken  string          `json:"reserve_token"`
	Result        string          `json:"result"` // new_conv
	Success       bool            `json:"success"`
}
