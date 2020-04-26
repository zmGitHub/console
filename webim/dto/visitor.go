package dto

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

var SearchVisitorFields = []string{}

type LandingPage struct {
	LandingPageTitle string `json:"landing_page_title"`
	LandingPageURL   string `json:"landing_page_url"`
}

type Source struct {
	SourceDomain  string `json:"source_domain"`
	SourceKeyword string `json:"source_keyword"`
	SourceSource  string `json:"source_source"`
	SourceURL     string `json:"source_url"`
}

type SearchVisitor struct {
	ID            string       `json:"id"`
	Address       string       `json:"address"`
	Age           int          `json:"age"`
	BrowserFamily string       `json:"browser_family"`
	City          string       `json:"city"`
	Comment       string       `json:"comment"`
	Country       string       `json:"country"`
	CreatedOn     time.Time    `json:"created_on"`
	Custom        interface{}  `json:"custom"`
	Email         string       `json:"email"`
	EnterpriseID  string       `json:"enterprise_id"`
	Gender        interface{}  `json:"gender"`
	LandingPage   *LandingPage `json:"landing_page"`
	Name          string       `json:"name"`
	OsFamily      string       `json:"os_family"`
	Province      string       `json:"province"`
	Qq            string       `json:"qq"`
	Source        *Source      `json:"source"`
	Tag           string       `json:"tag"`
	Tel           string       `json:"tel"`
	TrackID       string       `json:"track_id"`
	UpdatedOn     time.Time    `json:"updated_on"`
	VisitID       string       `json:"visit_id"`
	Weibo         string       `json:"weibo"`
	Weixin        string       `json:"weixin"`
}

type VisitorFieldsSetting struct {
	Category string `json:"category"`
	Key      string `json:"key"`
	Type     string `json:"type"`
	Visible  string `json:"visible"`
}

type UserConfig struct {
	Setting      []*VisitorFieldsSetting `json:"setting"`
	SyncSettings interface{}             `json:"sync_settings"`
}

type OsCategory struct {
	Mobile []interface{} `json:"mobile"`
	Other  []interface{} `json:"other"`
	Pc     []interface{} `json:"pc"`
}

type SourceHash struct {
	DirectAccess []interface{} `json:"direct_access"`
	ExternalLink []interface{} `json:"external_link"`
}

type VisitFilter struct {
	Locals         map[string][]interface{} `json:"locals"`
	OsCategoryHash *OsCategory              `json:"os_category_hash"`
	SourceHash     *SourceHash              `json:"source_hash"`
}

type Visitor struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Tel     string `json:"tel"`    // tel
	QQ      string `json:"qq"`     // qq
	Weixin  string `json:"weixin"` // weixin
	Weibo   string `json:"weibo"`
	Address string `json:"address"`
	Email   string `json:"email"`   // email
	Comment string `json:"comment"` // comment
}

func ModelVisitorToVisitor(v *models.Visitor) *Visitor {
	return &Visitor{
		Name:    v.Name,
		Age:     v.Age,
		Gender:  v.Gender,
		Tel:     v.Mobile,
		QQ:      v.QqNum,
		Weixin:  v.Wechat,
		Weibo:   v.Weibo,
		Address: v.Address,
		Email:   v.Email,
		Comment: v.Remark,
	}
}
