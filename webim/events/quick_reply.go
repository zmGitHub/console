package events

import "bitbucket.org/forfd/custm-chat/webim/common"

// m"{\"action\":\"quick_replies_refresh\",
// \"enterprise_id\":5869,\
// "id\":\"quick_reply_imports-561255\",
// \"trace_start\":1574580019.023858}"
type QuickReplayUpload struct {
	Action       string  `json:"action"`
	EnterpriseID string  `json:"enterprise_id"`
	ID           string  `json:"id"`
	TraceStart   float64 `json:"trace_start"`
}

func NewQuickReplayUploadEvent(e *QuickReplayUpload) string {
	s, err := common.Marshal(e)
	if err != nil {
		return ""
	}

	return s
}
