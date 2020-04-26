package adapter

import (
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AutoMsg struct {
	AutoReplyMsgSettings  *AutoReplyMsgSettings  `json:"auto_reply_msg_settings"`
	ChatLinkAutoMsgConfig *ChatLinkAutoMsgConfig `json:"chat_link_auto_msg_config"`
	ClientWakingAutoMsg   *ClientWakingAutoMsg   `json:"client_waking_auto_msg"`
	EndingMsgSettings     *EndingMsgSettings     `json:"ending_msg_settings"`
	PromotionMsgSettings  *PromotionMsgSettings  `json:"promotion_msg_settings"`
	WelcomeMsgSettings    *WelcomeMsgSettings    `json:"welcome_msg_settings"`
}

// Content      string `json:"content"`
//	ContentSdk   string `json:"content_sdk"`
//	Countdown    int    `json:"countdown"`
//	CreatedOn    string `json:"created_on"`
//	Enabled      bool   `json:"enabled"`
//	EnterpriseID string `json:"enterprise_id"`
//	ID           string `json:"id"`
//	Source       string `json:"source"`
//	Summary      string `json:"summary"`
//	Thumbnail    string `json:"thumbnail"`
//	UpdatedOn    string `json:"updated_on"`
type PromotionMessage struct {
	*models.PromotionMsg
	Content   string  `json:"content"`    // content
	CreatedOn *string `json:"created_on"` // created_on
	UpdatedOn *string `json:"updated_on"` // updated_on
}

func GetAutoMsgConfigsFromEntConfigs(configs *EnterpriseConfigs) *AutoMsg {
	if configs == nil {
		return &AutoMsg{}
	}

	return &AutoMsg{
		AutoReplyMsgSettings:  configs.AutoReplyMsgSettings,
		ChatLinkAutoMsgConfig: configs.ChatLinkAutoMsgConfig,
		ClientWakingAutoMsg:   configs.ClientWalkingAutoMsg,
		EndingMsgSettings:     configs.EndingMsgSettings,
		PromotionMsgSettings:  configs.PromotionMsgSettings,
		WelcomeMsgSettings:    configs.WelcomeMsgSettings,
	}
}
