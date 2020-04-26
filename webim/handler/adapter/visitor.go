package adapter

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type Facade struct {
	SearchEngine string `json:"search_engine"`
	LandingPage  string `json:"landing_page"`
	Location     int    `json:"location"`
	Keyword      string `json:"keyword"`
	SourceType   string `json:"source_type"`
	Returning    bool   `json:"returning"`
	SourcePage   string `json:"source_page"`
	Direct       bool   `json:"direct"`
}

type EnterpriseInfo struct {
	PublicSignature string `json:"public_signature"`
	PublicQq        string `json:"public_qq"`
	PublicEmail     string `json:"public_email"`
	PromptStatus    string `json:"prompt_status"`
	PublicCellphone string `json:"public_cellphone"`
	Name            string `json:"name"`
	PublicTelephone string `json:"public_telephone"`
	Location        string `json:"location"`
	PublicWeixin    string `json:"public_weixin"`
	PublicNickname  string `json:"public_nickname"`
	Fullname        string `json:"fullname"`
	Avatar          string `json:"avatar"`
}

type Survey struct {
	Status           string `json:"status"`
	HasSubmittedForm bool   `json:"has_submitted_form"`
}

type InitVisitResp struct {
	BrowserID                   string                   `json:"browser_id"`
	InvitationConfig            *InvitationConfig        `json:"invitation_config"`
	RobotSettings               *RobotSettings           `json:"robot_settings"`
	QueueingSettings            *QueueSettings           `json:"queueing_settings"`
	TicketConfig                *TicketConfig            `json:"ticket_config"`
	SendFileSettings            *SendFileSettings        `json:"send_file_settings"`
	VisitorStatusAgentToken     string                   `json:"visitor_status_agent_token"`
	Facade                      *Facade                  `json:"facade"`
	InQueue                     bool                     `json:"in_queue"`
	SearchEngine                string                   `json:"search_engine"`
	ServiceEvaluationConfig     *ServiceEvaluationConfig `json:"service_evaluation_config"`
	VisitID                     string                   `json:"visit_id"`
	Success                     bool                     `json:"success"`
	EntWelcomeMessage           string                   `json:"ent_welcome_message"`
	VisitPageID                 string                   `json:"visit_page_id"`
	TrackID                     string                   `json:"track_id"`
	BrowserFamily               string                   `json:"browser_family"`
	EnterpriseInfo              *EnterpriseInfo          `json:"enterprise_info"`
	Survey                      *SurveyConfig            `json:"survey"`
	Servability                 bool                     `json:"servability"`
	StandaloneWindowConfig      *StandaloneWindowConfig  `json:"standalone_window_config"`
	SchedulerAfterClientSendMsg bool                     `json:"scheduler_after_client_send_msg"`
	BaiduBidBlackList           []string                 `json:"baidu_bid_black_list"`
	VisitorStatus               int                      `json:"visitor_status"`
	WidgetSettings              *WidgetSettings          `json:"widget_settings"`
}

func EnterpriseRespToEnterpriseInfo(resp *EnterpriseResp) *EnterpriseInfo {
	return &EnterpriseInfo{
		Name:            resp.Name,
		Fullname:        resp.FullName,
		Avatar:          resp.Avatar,
		Location:        resp.Location,
		PromptStatus:    "open",
		PublicCellphone: resp.PublicCellphone,
		PublicEmail:     resp.PublicEmail,
		PublicNickname:  resp.PublicNickname,
		PublicQq:        resp.PublicQQ,
		PublicSignature: resp.PublicSignature,
		PublicWeixin:    resp.PublicWeixin,
	}
}

func ConvertModelVisitToVisit(visitor *models.Visitor, visitInfo *models.Visit) *VisitInfo {
	var name, avatar = "", ""
	if visitor != nil {
		name, avatar = visitor.Name, visitor.Avatar
	}

	return &VisitInfo{
		ID:                                   visitInfo.ID,
		EnterpriseID:                         visitInfo.EntID,
		AgentToken:                           "",
		AppName:                              "",
		AppVersion:                           "",
		Appkey:                               "",
		Avatar:                               avatar,
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
		Name:                                 name,
		NetType:                              "",
		OsCategory:                           visitInfo.OsCategory,
		OsFamily:                             visitInfo.OsFamily,
		OsLanguage:                           "",
		OsTimezone:                           "",
		OsVersion:                            visitInfo.OsVersion,
		OsVersionString:                      visitInfo.OsVersionString,
		Platform:                             visitInfo.Platform,
		Province:                             visitInfo.Province,
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
