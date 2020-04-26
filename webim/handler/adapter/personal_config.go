package adapter

type config struct {
	AtmeList       *AtmeList       `json:"atme_list"`
	ColleagueConvs *ColleagueConvs `json:"colleague_convs"`
	MyConvs        *MyConvs        `json:"my_convs"`
	QueueingList   *QueueingList   `json:"queueing_list"`
	VisitorList    *VisitorList    `json:"visitor_list"`
}

type AtmeList struct {
	Enable bool `json:"enable"`
}

type ColleagueConvs struct {
	Enable  bool     `json:"enable"`
	GroupBy []string `json:"group_by"`
}

type MyConvs struct {
	Enable               bool     `json:"enable"`
	GroupBy              []string `json:"group_by"`
	NonresponseThreshold int      `json:"nonresponse_threshold"`
}

type QueueingList struct {
	Enable bool `json:"enable"`
}

type VisitorList struct {
	Enable bool `json:"enable"`
}

type ConvOrderConfig struct {
	OrderRule string `json:"order_rule"`
	SinkDown  bool   `json:"sink_down"`
}

type msgTone struct {
	ConvTurnInTone       bool `json:"conv_turn_in_tone"`
	ConvTurnOutTone      bool `json:"conv_turn_out_tone"`
	NewColleagueConvTone bool `json:"new_colleague_conv_tone"`
	NewConvTone          bool `json:"new_conv_tone"`
	NewMsgTone           bool `json:"new_msg_tone"`
}

type webMsgTone struct {
	Desktop *msgTone `json:"desktop"`
	Voice   *msgTone `json:"voice"`
}
type MessageTone struct {
	Web *webMsgTone `json:"web"`
}

type WelcomeMsg struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ConvGroupConfig struct {
	Config *config `json:"config"`
	Mode   string  `json:"mode"`
}

// personal_config
type PersonalConfig struct {
	ConvGroupConfig *ConvGroupConfig `json:"conv_group_config"`
	ConvOrderConfig *ConvOrderConfig `json:"conv_order_config"`
	MessageTone     *MessageTone     `json:"message_tone"`
	QuickReplyRule  string           `json:"quick_reply_rule"`
	WelcomeMsg      *WelcomeMsg      `json:"welcome_msg"`
}

var DefaultPersonalConfig = &PersonalConfig{
	ConvGroupConfig: &ConvGroupConfig{
		Config: &config{
			AtmeList:       &AtmeList{Enable: false},
			ColleagueConvs: &ColleagueConvs{Enable: false, GroupBy: []string{}},
			MyConvs:        &MyConvs{Enable: false, GroupBy: []string{}, NonresponseThreshold: 0},
			QueueingList:   &QueueingList{Enable: false},
			VisitorList:    &VisitorList{Enable: false},
		},
		Mode: "common",
	},
	ConvOrderConfig: &ConvOrderConfig{OrderRule: "new_msg_first", SinkDown: false},
	MessageTone: &MessageTone{Web: &webMsgTone{
		Desktop: &msgTone{
			ConvTurnInTone:       true,
			ConvTurnOutTone:      true,
			NewColleagueConvTone: true,
			NewConvTone:          true,
			NewMsgTone:           true,
		},
		Voice: &msgTone{
			ConvTurnInTone:       true,
			ConvTurnOutTone:      true,
			NewColleagueConvTone: true,
			NewConvTone:          true,
			NewMsgTone:           true,
		},
	}},
	QuickReplyRule: "add",
	WelcomeMsg:     &WelcomeMsg{Status: "close"},
}
