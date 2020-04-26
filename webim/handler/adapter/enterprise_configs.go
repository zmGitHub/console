package adapter

// agents_permissions_config
type AgentsPermissionsConfig struct {
	AgentsPermissionsLevel int `json:"agents_permissions_level"`
}

//  "web": {
//      "content": "xxxxxx",
//      "countdown": 20,
//      "status": "open"
//    },
type webMsg struct {
	Content   string `json:"content"`
	CountDown int    `json:"countdown"`
	Status    string `json:"status"`
}

// auto_reply_msg_settings
type AutoReplyMsgSettings struct {
	Web         *webMsg `json:"web"`
	BaiduBCP    *webMsg `json:"baidu_bcp"`
	MiniProgram *webMsg `json:"mini_program"`
	SDK         *webMsg `json:"sdk"`
	Toutiao     *webMsg `json:"toutiao"`
	Weibo       *webMsg `json:"weibo"`
	Weixin      *webMsg `json:"weixin"`
}

// chat_link_auto_msg_config
type ChatLinkAutoMsgConfig struct {
	Web struct {
		Status string `json:"status"`
	} `json:"web"`
}

// client_waking_auto_msg
type ClientWakingAutoMsg struct {
	Web         *webMsg `json:"web"`
	BaiduBCP    *webMsg `json:"baidu_bcp"`
	MiniProgram *webMsg `json:"mini_program"`
	SDK         *webMsg `json:"sdk"`
	Toutiao     *webMsg `json:"toutiao"`
	Weibo       *webMsg `json:"weibo"`
	Weixin      *webMsg `json:"weixin"`
}

// conv_grade_config

type convLevel struct {
	AgentMsgCnt  int `json:"agent_msg_cnt"`
	ClientMsgCnt int `json:"client_msg_cnt"`
}

//"enable": true,
//    "first_level": {
//      "agent_msg_cnt": 0,
//      "client_msg_cnt": 19
//    },
//    "second_level": {
//      "agent_msg_cnt": 0,
//      "client_msg_cnt": 10
//    },
//    "third_level": {
//      "agent_msg_cnt": 0,
//      "client_msg_cnt": 5
//    }
type ConvGradeConfig struct {
	Enable      bool       `json:"enable"`
	FirstLevel  *convLevel `json:"first_level"`
	SecondLevel *convLevel `json:"second_level"`
	ThirdLevel  *convLevel `json:"third_level"`
}

type webConvEndConfig struct {
	NoMsgEnd   int `json:"no_msg_end"`
	OfflineEnd int `json:"offline_end"`
}

// end_conv_expire_config
// "mini_program": 3,
//    "sdk": 3,
//    "web": {
//      "no_msg_end": 3,
//      "offline_end": 30
//    },
//    "weibo": 3,
//    "weixin": 3
type EndConvExpireConfig struct {
	MiniProgram int               `json:"mini_program"`
	SDK         int               `json:"sdk"`
	Web         *webConvEndConfig `json:"web"`
	Weibo       int               `json:"weibo"`
	Weixin      int               `json:"weixin"`
}

// ending_msg_settings
// "web": {
//      "agent_ending_message": "感谢您的咨询，祝您生活工作顺利！（客服手动结束时）",
//      "auto_ending_message": "感谢您的咨询，祝您生活工作顺利！（系统自动关闭时给顾客发）",
//      "status": "open"
//    },
type endMsgSetting struct {
	AgentEndingMessage string `json:"agent_ending_message"`
	AutoEndingMessage  string `json:"auto_ending_message"`
	Status             string `json:"status"`
}
type EndingMsgSettings struct {
	Web         *endMsgSetting `json:"web"`
	BaiduBCP    interface{}    `json:"baidu_bcp"`
	MiniProgram interface{}    `json:"mini_program"`
	SDK         interface{}    `json:"sdk"`
	Toutiao     interface{}    `json:"toutiao"`
	Weibo       interface{}    `json:"weibo"`
	Weixin      interface{}    `json:"weixin"`
	// prompt_status: "open"
	PromptStatus string `json:"prompt_status"`
}

type InvitationConfig struct {
	Auto struct {
		Accept struct {
			Countdown int    `json:"countdown"`
			Status    string `json:"status"`
		} `json:"accept"`
		Countdown int `json:"countdown"`
		Reject    struct {
			Countdown int    `json:"countdown"`
			Type      string `json:"type"`
		} `json:"reject"`
		Status string `json:"status"`
	} `json:"auto"`
	Desktop struct {
		Actions []struct {
			Height   int    `json:"height"`
			ID       string `json:"id"`
			Position struct {
				Bottom string `json:"bottom"`
				Left   int    `json:"left"`
				Right  string `json:"right"`
				Top    int    `json:"top"`
			} `json:"position"`
			Type     int    `json:"type"`
			Width    int    `json:"width"`
			Href     string `json:"href,omitempty"`
			LinkType int    `json:"linkType,omitempty"`
		} `json:"actions"`
		Bgi struct {
			Height int    `json:"height"`
			Src    string `json:"src"`
			Width  int    `json:"width"`
		} `json:"bgi"`
		Position struct {
			Bottom int `json:"bottom"`
			Side   int `json:"side"`
			Type   int `json:"type"`
		} `json:"position"`
		Src  string `json:"src"`
		Text string `json:"text"`
		Type int    `json:"type"`
	} `json:"desktop"`
	FacadeStatus string `json:"facade_status"`
	Manual       struct {
		Accept struct {
			Countdown int    `json:"countdown"`
			Status    string `json:"status"`
		} `json:"accept"`
		Reject struct {
			Countdown int    `json:"countdown"`
			Type      string `json:"type"`
		} `json:"reject"`
		Status string `json:"status"`
	} `json:"manual"`
	Mobile struct {
		Actions []struct {
			Height   int    `json:"height"`
			ID       string `json:"id"`
			Position struct {
				Bottom string `json:"bottom"`
				Left   int    `json:"left"`
				Right  string `json:"right"`
				Top    int    `json:"top"`
			} `json:"position"`
			Type  int `json:"type"`
			Width int `json:"width"`
		} `json:"actions"`
		Bgi struct {
			Height int    `json:"height"`
			Src    string `json:"src"`
			Width  int    `json:"width"`
		} `json:"bgi"`
		Position struct {
			Type  int `json:"type"`
			Value int `json:"value"`
		} `json:"position"`
		Src  string `json:"src"`
		Text string `json:"text"`
		Type int    `json:"type"`
	} `json:"mobile"`
}

// oauth_settings
// "identity_key": "",
//    "retry_times": 0,
//    "secret_key": "",
//    "status": "close",
//    "success_result": "",
//    "url": ""
type OAuthSettings struct {
	IdentityKey   string `json:"identity_key"`
	RetryTimes    int    `json:"retry_times"`
	SecretKey     string `json:"secret_key"`
	Status        string `json:"status"`
	SuccessResult string `json:"success_result"`
	URL           string `json:"url"`
}

// promotion_msg_settings
// "content": [
//        21719
//      ],
//      "status": "open",
//      "stop_after_talk": false
type promotionMsgSetting struct {
	Content       []string `json:"content"`
	Status        string   `json:"status"`
	StopAfterTalk bool     `json:"stop_after_talk"`
}

type PromotionMsgSettings struct {
	Web         *promotionMsgSetting `json:"web"`
	BaiduBCP    *promotionMsgSetting `json:"baidu_bcp"`
	MiniProgram *promotionMsgSetting `json:"mini_program"`
	SDK         *promotionMsgSetting `json:"sdk"`
	Toutiao     *promotionMsgSetting `json:"toutiao"`
	Weibo       *promotionMsgSetting `json:"weibo"`
	Weixin      *promotionMsgSetting `json:"weixin"`
}

// queueing_settings
type QueueSettings struct {
	Intro     string `json:"intro"`
	QueueSize int    `json:"queue_size"`
	Status    string `json:"status"`
}

// reserve_clues_config
//  "enabled": true,
//    "fallback": "allocate_rule"
type ReserveCluesConfig struct {
	Enabled  bool   `json:"enabled"`
	Fallback string `json:"fallback"`
}

// robot_settings
type RobotSettings struct {
	AutoCompleteOnlyQueryMainQ bool   `json:"auto_complete_only_query_main_q"`
	Avatar                     string `json:"avatar"`
	CorrelationThreshold       int    `json:"correlation_threshold"`
	FailedThreshold            int    `json:"failed_threshold"`
	LeftMsgCnt                 int    `json:"left_msg_cnt"`
	ManualRedirect             bool   `json:"manual_redirect"`
	MoreLikeThisCount          int    `json:"more_like_this_count"`
	Nickname                   string `json:"nickname"`
	Provider                   string `json:"provider"`
	ResponseCantAnswer         []struct {
		Text string `json:"text"`
	} `json:"response_cant_answer"`
	ResponseEvalUseful []struct {
		Text string `json:"text"`
	} `json:"response_eval_useful"`
	ResponseEvalUseless []struct {
		Text string `json:"text"`
	} `json:"response_eval_useless"`
	ResponseManualRedirect []struct {
		Text string `json:"text"`
	} `json:"response_manual_redirect"`
	ResponseMoreLikeThis []struct {
		Text string `json:"text"`
	} `json:"response_more_like_this"`
	ResponseQueueing []struct {
		Text string `json:"text"`
	} `json:"response_queueing"`
	ResponseRedirect []struct {
		Text string `json:"text"`
	} `json:"response_redirect"`
	ResponseReply []struct {
		Text string `json:"text"`
	} `json:"response_reply"`
	Rule               string `json:"rule"`
	ShowSwitch         bool   `json:"show_switch"`
	Signature          string `json:"signature"`
	Status             string `json:"status"`
	UnmatchThreshold   int    `json:"unmatch_threshold"`
	WelcomeQuestionIds []int  `json:"welcome_question_ids"`
	WelcomeText        string `json:"welcome_text"`
}

type SalescloudConfig struct {
	DataSync struct {
		Status     string   `json:"status"`
		SyncFields []string `json:"sync_fields"`
	} `json:"data_sync"`
	SalesInfo struct {
		FollowRecords string `json:"follow_records"`
		HistoryOrder  string `json:"history_order"`
		Status        string `json:"status"`
	} `json:"sales_info"`
	Status string `json:"status"`
	Token  string `json:"token"`
}

// send_file_settings
type SendFileSettings struct {
	WidgetStatus string `json:"widget_status"`
}

// service_evaluation_config
// "agent_invitation": "close",
//    "agent_visible": "close",
//    "auto_invitation": "open",
//    "prompt_text": "请您为我的服务做出评价"
type ServiceEvaluationConfig struct {
	AgentInvitation string `json:"agent_invitation"`
	AgentVisible    string `json:"agent_visible"`
	AutoInvitation  string `json:"auto_invitation"`
	PromptText      string `json:"prompt_text"`
}

// standalone_window_config
// "background": {
//      "color": "c4cf48",
//      "url": ""
//    },
//    "desktop": {
//      "customer_content": "<div>778</div><div><img src=https://s3-qcloud.meiqia.com/pics.meiqia.bucket/6ab41102558be9a21dd500dd6ce39a00.png style=\"max-width: 100%;\"></div><div><br></div><div><img src=https://s3-qcloud.meiqia.com/pics.meiqia.bucket/709e561b79d497069df012e12e44ffd7.jpg style=\"max-width: 100%;\"></div><div><br></div><div><br></div>",
//      "customer_photo_type": "small",
//      "theme": [
//        "8686a5",
//        "white",
//        "3ad531"
//      ],
//      "type": "fusion"
//    },
//    "mobile": {
//      "theme": [
//        "573942",
//        "white"
//      ],
//      "type": "mustang"
//    },
//    "removeBrand": "close",
//    "ring": "open"
type StandaloneWindowConfig struct {
	Background struct {
		Color string `json:"color"`
		URL   string `json:"url"`
	} `json:"background"`

	Desktop struct {
		CustomerContent   string   `json:"customer_content"`
		CustomerPhotoType string   `json:"customer_photo_type"`
		Theme             []string `json:"theme"`
		Type              string   `json:"type"`
	} `json:"desktop"`

	Mobile struct {
		Theme []string `json:"theme"`
		Type  string   `json:"type"`
	} `json:"mobile"`

	RemoveBrand string `json:"removeBrand"`
	Ring        string `json:"ring"`
}

type SurveyConfig struct {
	// has_submitted_form
	HasSubmittedForm bool `json:"has_submitted_form"`
	// status
	Status string `json:"status"`
}

type TicketConfig struct {
	Captcha                string `json:"captcha"`
	Category               string `json:"category"`
	ContactRule            string `json:"contactRule"`
	DefaultTemplate        string `json:"defaultTemplate"`
	DefaultTemplateContent string `json:"defaultTemplateContent"`
	Email                  string `json:"email"`
	Intro                  string `json:"intro"`
	Name                   string `json:"name"`
	Permission             string `json:"permission"`
	Qq                     string `json:"qq"`
	Tel                    string `json:"tel"`
	Wechat                 string `json:"wechat"`
}

// timeout_redirect_config
//  "countdown": 20,
//    "rules": [
//      {
//        "countdown": 20,
//        "type": "group_first"
//      }
//    ],
//    "status": "open"
type TimeoutRedirectConfig struct {
	CountDown int `json:"countdown"`
	Rules     []struct {
		CountDown int    `json:"countdown"`
		Type      string `json:"type"`
	} `json:"rules"`
	Status string `json:"status"`
}

// visitor_visible
type VisitorVisible struct {
	Region string `json:"region"`
}

// web_callback_settings
// "callback_switch": "open",
//    "captcha_switch": "open"
type WebCallbackSettings struct {
	CallbackSwitch string `json:"callback_switch"`
	CaptchaSwitch  string `json:"captcha_switch"`
}

// welcome_msg_settings
type welcomeMsgSetting struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

type WelcomeMsgSettings struct {
	Web         *welcomeMsgSetting `json:"web"`
	BaiduBCP    *welcomeMsgSetting `json:"baidu_bcp"`
	MiniProgram *welcomeMsgSetting `json:"mini_program"`
	SDK         *welcomeMsgSetting `json:"sdk"`
	Toutiao     *welcomeMsgSetting `json:"toutiao"`
	Weibo       *welcomeMsgSetting `json:"weibo"`
	Weixin      *welcomeMsgSetting `json:"weixin"`
}

// widget_settings
type WidgetSettings struct {
	Desktop struct {
		Btn struct {
			Icon struct {
				Offline int `json:"offline"`
				Online  int `json:"online"`
			} `json:"icon"`
			Picture struct {
				Offline string `json:"offline"`
				Online  string `json:"online"`
			} `json:"picture"`
			Position struct {
				Bottom     string `json:"bottom"`
				Horizontal string `json:"horizontal"`
				Type       string `json:"type"`
			} `json:"position"`
			Preview string `json:"preview"`
			Text    struct {
				Offline string `json:"offline"`
				Online  string `json:"online"`
			} `json:"text"`
			Theme string `json:"theme"`
			Type  string `json:"type"`
		} `json:"btn"`
		Panel struct {
			CustomerContent   string `json:"customer_content"`
			CustomerPhotoType string `json:"customer_photo_type"`
			Position          struct {
				Bottom     string `json:"bottom"`
				Horizontal string `json:"horizontal"`
				Type       string `json:"type"`
			} `json:"position"`
			Theme []string `json:"theme"`
			Type  string   `json:"type"`
		} `json:"panel"`
		Pop bool `json:"pop"`
	} `json:"desktop"`
	HTTPS  string `json:"https"`
	Mobile struct {
		Btn struct {
			Icon struct {
				Offline int `json:"offline"`
				Online  int `json:"online"`
			} `json:"icon"`
			Picture struct {
				Offline string `json:"offline"`
				Online  string `json:"online"`
			} `json:"picture"`
			Position struct {
				Bottom     string `json:"bottom"`
				Horizontal string `json:"horizontal"`
				Type       string `json:"type"`
			} `json:"position"`
			Preview string `json:"preview"`
			Text    struct {
				Offline string `json:"offline"`
				Online  string `json:"online"`
			} `json:"text"`
			Theme string `json:"theme"`
			Type  string `json:"type"`
		} `json:"btn"`
		Panel struct {
			Position struct {
				Bottom     int    `json:"bottom"`
				Horizontal int    `json:"horizontal"`
				Type       string `json:"type"`
			} `json:"position"`
			Theme []string `json:"theme"`
			Type  string   `json:"type"`
		} `json:"panel"`
		Pop bool `json:"pop"`
	} `json:"mobile"`
	RemoveBrand string `json:"removeBrand"`
	Ring        string `json:"ring"`
	TicketOnly  string `json:"ticketOnly"`
}

// scheduler_after_client_send_msg
// "is_activated": true,
//  "is_from_baidu_open": true,
//  "is_group": false,
type EnterpriseConfigs struct {
	AgentsPermissionsConfig     *AgentsPermissionsConfig `json:"agents_permissions_config"`
	AutoReplyMsgSettings        *AutoReplyMsgSettings    `json:"auto_reply_msg_settings"`
	ChatLinkAutoMsgConfig       *ChatLinkAutoMsgConfig   `json:"chat_link_auto_msg_config"`
	ClientWalkingAutoMsg        *ClientWakingAutoMsg     `json:"client_waking_auto_msg"`
	ConvGradeConfig             *ConvGradeConfig         `json:"conv_grade_config"`
	EndConvExpireConfig         *EndConvExpireConfig     `json:"end_conv_expire_config"`
	EndingMsgSettings           *EndingMsgSettings       `json:"ending_msg_settings"`
	InvitationConfig            *InvitationConfig        `json:"invitation_config"`
	OAuthSettings               *OAuthSettings           `json:"oauth_settings"`
	PromotionMsgSettings        *PromotionMsgSettings    `json:"promotion_msg_settings"`
	QueueSettings               *QueueSettings           `json:"queueing_settings"`
	ReserveCluesConfig          *ReserveCluesConfig      `json:"reserve_clues_config"`
	RobotSettings               *RobotSettings           `json:"robot_settings"`
	SalesCloudConfig            *SalescloudConfig        `json:"sales_cloud_config"`
	SendFileSettings            *SendFileSettings        `json:"send_file_settings"`
	ServiceEvaluationConfig     *ServiceEvaluationConfig `json:"service_evaluation_config"`
	StandaloneWindowConfig      *StandaloneWindowConfig  `json:"standalone_window_config"`
	Survey                      *SurveyConfig            `json:"survey"`
	TicketConfig                *TicketConfig            `json:"ticket_config"`
	TimeoutRedirectConfig       *TimeoutRedirectConfig   `json:"timeout_redirect_config"`
	VisitorVisible              *VisitorVisible          `json:"visitor_visible"`
	WebCallbackSettings         *WebCallbackSettings     `json:"web_callback_settings"`
	WelcomeMsgSettings          *WelcomeMsgSettings      `json:"welcome_msg_settings"`
	WidgetSettings              *WidgetSettings          `json:"widget_settings"`
	IsActivated                 bool                     `json:"is_activated"`
	IsFromBaiduOpen             bool                     `json:"is_from_baidu_open"`
	IsGroup                     bool                     `json:"is_group"`
	SchedulerAfterClientSendMsg bool                     `json:"scheduler_after_client_send_msg"`
}

type ConfigOption func(*EnterpriseConfigs)

func WithAgentsPermissionsConfig(s *AgentsPermissionsConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.AgentsPermissionsConfig = s
	}
}

func WithAutoReplyMsgSettings(s *AutoReplyMsgSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.AutoReplyMsgSettings = s
	}
}

func WithChatLinkAutoMsgConfig(s *ChatLinkAutoMsgConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.ChatLinkAutoMsgConfig = s
	}
}

func WithClientWalkingAutoMsg(s *ClientWakingAutoMsg) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.ClientWalkingAutoMsg = s
	}
}

func WithConvGradeConfig(s *ConvGradeConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.ConvGradeConfig = s
	}
}

func WithEndConvExpireConfig(s *EndConvExpireConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.EndConvExpireConfig = s
	}
}

func WithEndingMsgSettings(s *EndingMsgSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.EndingMsgSettings = s
	}
}

func WithInvitationConfig(s *InvitationConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.InvitationConfig = s
	}
}

func WithOAuthSettings(s *OAuthSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.OAuthSettings = s
	}
}

func WithPromotionMsgSettings(s *PromotionMsgSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.PromotionMsgSettings = s
	}
}

func WithQueueSettings(s *QueueSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.QueueSettings = s
	}
}

func WithReserveCluesConfig(s *ReserveCluesConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.ReserveCluesConfig = s
	}
}

func WithRobotSettings(s *RobotSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.RobotSettings = s
	}
}

func WithSalesCloudConfig(s *SalescloudConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.SalesCloudConfig = s
	}
}

func WithSendFileSettings(s *SendFileSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.SendFileSettings = s
	}
}

func WithServiceEvaluationConfig(s *ServiceEvaluationConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.ServiceEvaluationConfig = s
	}
}

func WithStandaloneWindowConfig(s *StandaloneWindowConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.StandaloneWindowConfig = s
	}
}

func WithTicketConfig(s *TicketConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.TicketConfig = s
	}
}

func WithTimeoutRedirectConfig(s *TimeoutRedirectConfig) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.TimeoutRedirectConfig = s
	}
}

func WithVisitorVisible(s *VisitorVisible) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.VisitorVisible = s
	}
}

func WithWebCallbackSettings(s *WebCallbackSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.WebCallbackSettings = s
	}
}

func WithWelcomeMsgSettings(s *WelcomeMsgSettings) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.WelcomeMsgSettings = s
	}
}

func WithIsActivated(s bool) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.IsActivated = s
	}
}

func WithIsFromBaiduOpen(s bool) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.IsFromBaiduOpen = s
	}
}

func WithIsGroup(s bool) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.IsGroup = s
	}
}

func WithSchedulerAfterClientSendMsg(s bool) ConfigOption {
	return func(configs *EnterpriseConfigs) {
		configs.SchedulerAfterClientSendMsg = s
	}
}

func BuildEnterpriseConfigs(opts ...ConfigOption) *EnterpriseConfigs {
	entConfig := &EnterpriseConfigs{}
	for _, opt := range opts {
		opt(entConfig)
	}

	return entConfig
}
