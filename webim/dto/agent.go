package dto

import "time"

type Agent struct {
	ID              string     `json:"id"`
	Avatar          string     `json:"avatar"`
	Cellphone       string     `json:"cellphone"`
	CreatedOn       *string    `json:"created_on"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Email           string     `json:"email"`
	EmailActivated  bool       `json:"email_activated"`
	EnterpriseID    string     `json:"enterprise_id"`
	GroupID         string     `json:"group_id"`
	IsOnline        bool       `json:"is_online"`
	Nickname        string     `json:"nickname"`
	Privilege       string     `json:"privilege"`
	PrivilegeRange  string     `json:"privilege_range"`
	PublicCellphone string     `json:"public_cellphone"`
	PublicEmail     string     `json:"public_email"`
	Qq              string     `json:"qq"`
	Rank            int        `json:"rank"`
	ReadFeatureID   int        `json:"read_feature_id"`
	Realname        string     `json:"realname"`
	ServingLimit    int        `json:"serving_limit"`
	Signature       string     `json:"signature"`
	Status          string     `json:"status"`
	Telephone       string     `json:"telephone"`
	Token           string     `json:"token"`
	Weixin          string     `json:"weixin"`
	WorkNum         string     `json:"work_num"`
}
