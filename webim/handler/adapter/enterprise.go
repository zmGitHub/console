package adapter

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

/*
allocation_rule: "roundrobin"
avatar: "https://s3-qcloud.meiqia.com/pics.meiqia.bucket/avatars/20180209/73c6647092456d1632c91be8cd57e45d.jpg"
city: "null"
contact_email: "85196508438@qq.com"
contact_name: "一二三四五六七八九八七六五四三二一"
contact_telephone: "12345"
created_on: "2016-01-13T04:07:55.493572"
fullname: "1"
id: 5869
industry: "互联网/软件"
last_active_time: null
location: null
mailing_address: "高新区软件园e5区4楼"
name: "美洽公司"
province: "北京市"
public_cellphone: "11381381388"
public_email: "team@meiqia.com"
public_nickname: "美洽客服团队"
public_qq: "1234567"
public_signature: "您想了解的美洽产品"
public_telephone: "88888888"
public_weixin: "7654321"
returning_customer_enabled: true
selecting_rule_enabled: true
special_serving_limit: 0
superadmin_id: 32824
telephone: null
token: "fa417e72be909a80fa891c6f66cd6752"
type: "platform"
website: null
*/
type EnterpriseResp struct {
	AllocationRule           string    `json:"allocation_rule"`
	Avatar                   string    `json:"avatar"`
	City                     string    `json:"city"`
	ContactEmail             string    `json:"contact_email"`
	ContactName              string    `json:"contact_name"`
	ContactTelephone         string    `json:"contact_telephone"`
	CreatedOn                time.Time `json:"created_on"`
	FullName                 string    `json:"fullname"`
	ID                       string    `json:"id"`
	Industry                 string    `json:"industry"`
	LastActiveTime           time.Time `json:"last_active_time"`
	Location                 string    `json:"location"`
	MailingAddress           string    `json:"mailing_address"`
	Name                     string    `json:"name"`
	Province                 string    `json:"province"`
	PublicCellphone          string    `json:"public_cellphone"`
	PublicEmail              string    `json:"public_email"`
	PublicNickname           string    `json:"public_nickname"`
	PublicQQ                 string    `json:"public_qq"`
	PublicSignature          string    `json:"public_signature"`
	PublicTelephone          string    `json:"public_telephone"`
	PublicWeixin             string    `json:"public_weixin"`
	ReturningCustomerEnabled bool      `json:"returning_customer_enabled"`
	SelectingRuleEnabled     bool      `json:"selecting_rule_enabled"`
	SpecialServingLimit      int       `json:"special_serving_limit"`
	SuperadminID             *string   `json:"superadmin_id"`
	Telephone                string    `json:"telephone"`
	Token                    string    `json:"token"`
	Type                     int       `json:"type"`
	Website                  string    `json:"website"`
}

func ConvertEntToAdapterEnt(ent *models.Enterprise) *EnterpriseResp {
	return &EnterpriseResp{
		AllocationRule:           ent.AllocationRule,
		Avatar:                   ent.Avatar,
		City:                     ent.City,
		ContactEmail:             ent.ContactEmail,
		ContactName:              ent.ContactName,
		ContactTelephone:         ent.ContactMobile,
		CreatedOn:                ent.CreatedAt,
		FullName:                 ent.FullName,
		ID:                       ent.ID,
		Industry:                 ent.Industry,
		LastActiveTime:           ent.LastActivatedAt,
		Location:                 ent.Location,
		MailingAddress:           ent.Address,
		Name:                     ent.Name,
		Province:                 ent.Province,
		PublicCellphone:          ent.Mobile,
		PublicEmail:              ent.Email,
		PublicNickname:           ent.Name,
		PublicQQ:                 ent.ContactQq,
		PublicSignature:          ent.ContactSignature,
		PublicTelephone:          ent.Mobile,
		PublicWeixin:             ent.ContactWechat,
		ReturningCustomerEnabled: false,
		SelectingRuleEnabled:     ent.AllocationRule != "",
		SpecialServingLimit:      0,
		SuperadminID:             &ent.AdminID,
		Telephone:                ent.Mobile,
		Token:                    "",
		Type:                     ent.Plan,
		Website:                  ent.Website,
	}
}
