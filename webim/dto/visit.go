package dto

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

type VisitPage struct {
	CreatedOn     string `json:"created_on"`
	Source        string `json:"source"`
	SourceKeyword string `json:"source_keyword"`
	SourceURL     string `json:"source_url"`
	Title         string `json:"title"`
	URL           string `json:"url"`
}

type VisitPages struct {
	Pages []*VisitPage `json:"pages"`
}

func ModelVisitInfoToVisitInfo(visitInfo *models.Visit, visitor *models.Visitor) *VisitInfo {
	var province = visitInfo.Province
	var city = visitInfo.City
	if common.IsMunicipality(visitInfo.City) {
		province = visitInfo.City
		city = ""
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
		City:                                 city,
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
