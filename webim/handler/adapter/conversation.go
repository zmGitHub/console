package adapter

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type VisitInfo struct {
	AgentToken                           string    `json:"agent_token"`
	AppName                              string    `json:"app_name"`
	AppVersion                           string    `json:"app_version"`
	Appkey                               string    `json:"appkey"`
	Avatar                               string    `json:"avatar"`
	BrowserFamily                        string    `json:"browser_family"`
	BrowserVersion                       string    `json:"browser_version"`
	BrowserVersionString                 string    `json:"browser_version_string"`
	City                                 string    `json:"city"`
	Country                              string    `json:"country"`
	CreatedOn                            time.Time `json:"created_on"`
	DeviceBrand                          string    `json:"device_brand"`
	DeviceFamily                         string    `json:"device_family"`
	DeviceModel                          string    `json:"device_model"`
	DeviceToken                          string    `json:"device_token"`
	EnterpriseID                         string    `json:"enterprise_id"`
	FirstVisitPageDomainBySession        string    `json:"first_visit_page_domain_by_session"`
	FirstVisitPageSourceBySession        string    `json:"first_visit_page_source_by_session"`
	FirstVisitPageSourceDomainBySession  string    `json:"first_visit_page_source_domain_by_session"`
	FirstVisitPageSourceKeywordBySession string    `json:"first_visit_page_source_keyword_by_session"`
	FirstVisitPageSourceURLBySession     string    `json:"first_visit_page_source_url_by_session"`
	FirstVisitPageTitleBySession         string    `json:"first_visit_page_title_by_session"`
	FirstVisitPageURLBySession           string    `json:"first_visit_page_url_by_session"`
	ID                                   string    `json:"id"`
	IP                                   string    `json:"ip"`
	Isp                                  string    `json:"isp"`
	LastTitle                            string    `json:"last_title"`
	LastURL                              string    `json:"last_url"`
	LastVisitID                          string    `json:"last_visit_id"`
	LastVisitPageTitleBySession          string    `json:"last_visit_page_title_by_session"`
	LastVisitPageURLBySession            string    `json:"last_visit_page_url_by_session"`
	Name                                 string    `json:"name"`
	NetType                              string    `json:"net_type"`
	OsCategory                           string    `json:"os_category"`
	OsFamily                             string    `json:"os_family"`
	OsLanguage                           string    `json:"os_language"`
	OsTimezone                           string    `json:"os_timezone"`
	OsVersion                            string    `json:"os_version"`
	OsVersionString                      string    `json:"os_version_string"`
	Platform                             string    `json:"platform"`
	Province                             string    `json:"province"`
	ResidenceTimeSec                     int       `json:"residence_time_sec"`
	ResidenceTimeSecBySession            int       `json:"residence_time_sec_by_session"`
	Resolution                           string    `json:"resolution"`
	SdkImageURL                          string    `json:"sdk_image_url"`
	SdkName                              string    `json:"sdk_name"`
	SdkSource                            string    `json:"sdk_source"`
	SdkVersion                           string    `json:"sdk_version"`
	Status                               int       `json:"status"`
	StatusOn                             string    `json:"status_on"`
	TrackID                              string    `json:"track_id"`
	UaString                             string    `json:"ua_string"`
	VisitCnt                             int       `json:"visit_cnt"`
	VisitID                              string    `json:"visit_id"`
	VisitPageCnt                         int       `json:"visit_page_cnt"`
	VisitPageCntBySession                int       `json:"visit_page_cnt_by_session"`
	VisitedOn                            time.Time `json:"visited_on"`
}

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
	EvaLevel              *int        `json:"eva_level"`
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

type ActiveConversationsResp struct {
	Conversations []*Conversation `json:"conversations"`
}

type HistoryConversation struct {
	*Conversation
	Messages []*Message `json:"messages"`
}

type HistoryConversationsResp struct {
	Conversations []*HistoryConversation `json:"conversations"`
}

type SearchConversation struct {
	*Conversation
	CreatedOn string  `json:"created_on"`
	EndedOn   *string `json:"ended_on"`
}

type ColleageConvs struct {
	ConvNums      int             `json:"conv_nums"`
	Conversations []*Conversation `json:"conversations"`
	FromCache     bool            `json:"from_cache"`
	HasNext       bool            `json:"has_next"`
	NextCursor    *int            `json:"next_cursor"`
}

type Stream struct {
	Action         string      `json:"action"`
	Agent          *Agent      `json:"agent"`
	AgentID        string      `json:"agent_id"`
	Content        string      `json:"content"`
	ContentType    string      `json:"content_type"`
	ConversationID string      `json:"conversation_id"`
	CreatedOn      string      `json:"created_on"`
	EnterpriseID   string      `json:"enterprise_id"`
	Extra          interface{} `json:"extra"`
	FromType       string      `json:"from_type"`
	ID             int64       `json:"id"`
	ReadOn         *time.Time  `json:"read_on"`
	TrackID        string      `json:"track_id"`
	Type           string      `json:"type"`
}

type InitConvAction struct {
	Action        string        `json:"action"` // init_conv
	AgentID       string        `json:"agent_id"`
	AgentNickname string        `json:"agent_nickname"`
	Body          *Conversation `json:"body"`
	CreatedOn     time.Time     `json:"created_on"`
	EnterpriseID  string        `json:"enterprise_id"`
	ID            string        `json:"id"`
	TargetID      string        `json:"target_id"`
	TargetKind    string        `json:"target_kind"` // conv
	TrackID       string        `json:"track_id"`
}

type ConversationStreams struct {
	ConvIDs []string      `json:"conv_ids"`
	Streams []interface{} `json:"streams"`
	Total   int           `json:"total"`
}

type EndConversationMessageBody struct {
	AgentMsgNum  int    `json:"agent_msg_num"`
	ClientMsgNum int    `json:"client_msg_num"`
	ConvID       string `json:"conv_id"`
	EndedBy      string `json:"ended_by"`
	Evaluation   bool   `json:"evaluation"`
	MsgNum       int    `json:"msg_num"`
}

type EndConversationMessage struct {
	Action        string                      `json:"action"`
	AgentID       string                      `json:"agent_id"`
	AgentNickname string                      `json:"agent_nickname"`
	Body          *EndConversationMessageBody `json:"body"`
	CreatedOn     *string                     `json:"created_on"`
	EnterpriseID  string                      `json:"enterprise_id"`
	ID            string                      `json:"id"`
	Realname      string                      `json:"realname"`
	Source        string                      `json:"source"`
	TargetID      string                      `json:"target_id"`
	TargetKind    string                      `json:"target_kind"`
	TraceStart    int64                       `json:"trace_start"`
	TrackID       string                      `json:"track_id"`
}

func ModelConversationToConversation(conv *models.Conversation, visitID string, visitInfo *models.Visit, visitor *models.Visitor, tags []*models.VisitorTagRelation) *Conversation {
	var clientFirstSendTime, clientLastSendTime, firstMsgCreateTime *time.Time

	var endAt *time.Time
	if conv.EndedAt.Valid {
		t := common.ConvertUTCToLocal(conv.EndedAt.Time)
		endAt = &t
	}

	if conv.FirstMsgCreatedAt.Valid {
		clientFirstSendTime = &conv.FirstMsgCreatedAt.Time
		firstMsgCreateTime = &conv.FirstMsgCreatedAt.Time
	}

	if conv.LastMsgCreatedAt.Valid {
		clientLastSendTime = &conv.LastMsgCreatedAt.Time
	}

	var lastMsgCreateTime *time.Time
	if conv.LastMsgCreatedAt.Valid {
		lastMsgCreateTime = &conv.LastMsgCreatedAt.Time
	}

	var visit *VisitInfo
	if visitInfo != nil && visitor != nil {
		visit = ModelVisitInfoToVisitInfo(visitInfo, visitor)
	}

	var tagIDs []string
	if visitor != nil {
		for _, vr := range tags {
			if vr.VisitorID == visitor.ID {
				tagIDs = append(tagIDs, vr.TagID)
			}
		}
	}

	var evaLevel *int
	evaLevel = &conv.EvalLevel

	adapterConv := &Conversation{
		ID:                    conv.ID,
		AgentID:               conv.AgentID,
		EnterpriseID:          conv.EntID,
		AgentMsgNum:           int(conv.AgentMsgCount),
		AgentEffectiveMsgNum:  int(conv.AgentMsgCount),
		AgentType:             models.ConversationHumanAgentType,
		Assignee:              "",
		ClientFirstSendTime:   clientFirstSendTime,
		ClientLastSendTime:    clientLastSendTime,
		ClientMsgNum:          int(conv.ClientMsgCount),
		Clues:                 nil,
		ConverseDuration:      int(conv.Duration),
		CreatedOn:             common.ConvertUTCToLocal(conv.CreatedAt),
		EndedBy:               conv.EndedBy.String,
		EndedOn:               endAt,
		EvaContent:            conv.EvalContent,
		EvaLevel:              evaLevel,
		FirstMsgCreatedOn:     firstMsgCreateTime,
		FirstResponseWaitTime: conv.FirstResponseWaitTime.Int64,
		HasSummary:            false,
		LastMsgContent:        conv.LastMsgContent.String,
		LastMsgCreatedOn:      lastMsgCreateTime,
		LastUpdated:           &conv.UpdateAt,
		MsgNum:                int(conv.MsgCount),
		QualityGrade:          conv.QualityGrade.String,
		Tags:                  tagIDs,
		Title:                 conv.Title,
		TrackID:               conv.TraceID,
		URL:                   "",
		VisitID:               visitID,
		VisitInfo:             visit,
	}

	return adapterConv
}

func ModelVisitInfoToVisitInfo(visitInfo *models.Visit, visitor *models.Visitor) *VisitInfo {
	var province = visitInfo.Province
	if common.IsMunicipality(visitInfo.City) {
		province = visitInfo.City
	}

	return &VisitInfo{
		ID:                                   visitInfo.ID,
		EnterpriseID:                         visitInfo.EntID,
		AgentToken:                           "",
		AppName:                              "",
		AppVersion:                           "",
		Appkey:                               "",
		Avatar:                               visitor.Avatar,
		BrowserFamily:                        visitInfo.BrowserFamily,
		BrowserVersion:                       visitInfo.BrowserVersion,
		BrowserVersionString:                 visitInfo.BrowserVersionString,
		City:                                 visitInfo.City,
		Country:                              visitInfo.Country,
		CreatedOn:                            common.ConvertUTCToLocal(visitInfo.CreatedAt),
		DeviceBrand:                          visitInfo.OsCategory,
		DeviceFamily:                         visitInfo.OsFamily,
		DeviceModel:                          visitInfo.OsVersion,
		DeviceToken:                          "",
		FirstVisitPageDomainBySession:        visitInfo.FirstPageDomain,
		FirstVisitPageSourceBySession:        visitInfo.FirstPageSource,
		FirstVisitPageSourceDomainBySession:  visitInfo.FirstPageSourceDomain,
		FirstVisitPageSourceKeywordBySession: visitInfo.FirstPageSourceKeyword,
		FirstVisitPageSourceURLBySession:     visitInfo.FirstPageSourceURL,
		FirstVisitPageTitleBySession:         visitInfo.FirstPageTitle,
		FirstVisitPageURLBySession:           visitInfo.FirstPageURL,
		IP:                                   visitInfo.IP,
		Isp:                                  visitInfo.Isp,
		LastTitle:                            visitInfo.LatestTitle,
		LastURL:                              visitInfo.LatestURL,
		LastVisitID:                          "",
		LastVisitPageTitleBySession:          visitInfo.FirstPageTitle,
		LastVisitPageURLBySession:            visitInfo.FirstPageURL,
		Name:                                 visitor.Name,
		NetType:                              "",
		OsCategory:                           visitInfo.OsCategory,
		OsFamily:                             visitInfo.OsFamily,
		OsLanguage:                           "",
		OsTimezone:                           "",
		OsVersion:                            visitInfo.OsVersion,
		OsVersionString:                      visitInfo.OsVersionString,
		Platform:                             visitInfo.Platform,
		Province:                             province,
		ResidenceTimeSec:                     visitInfo.ResidenceTimeSec,
		ResidenceTimeSecBySession:            visitInfo.ResidenceTimeSec,
		Resolution:                           "",
		SdkImageURL:                          "",
		SdkName:                              "",
		SdkSource:                            "",
		SdkVersion:                           "",
		Status:                               0,
		StatusOn:                             "",
		TrackID:                              visitInfo.TraceID,
		UaString:                             visitInfo.UaString,
		VisitCnt:                             visitInfo.VisitPageCnt,
		VisitID:                              visitInfo.ID,
		VisitPageCnt:                         visitInfo.VisitPageCnt,
		VisitPageCntBySession:                0,
		VisitedOn:                            common.ConvertUTCToLocal(visitInfo.CreatedAt),
	}
}
