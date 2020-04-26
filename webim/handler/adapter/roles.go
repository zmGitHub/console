package adapter

import (
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type Role struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type RolesResp struct {
	Roles []*Role `json:"roles"`
}

type RolePerms struct {
	AccessConfig      []*Perm `json:"access_config"`
	BotConfig         []*Perm `json:"bot_config"`
	DataReport        []*Perm `json:"data_report"`
	Engage            []*Perm `json:"engage"`
	EntInfo           []*Perm `json:"ent_info"`
	HistoryConv       []*Perm `json:"history_conv"`
	OnlineAgentConfig []*Perm `json:"online_agent_config"`
	SalescloudConfig  []*Perm `json:"salescloud_config"`
	Ticket            []*Perm `json:"ticket"`
	VisitorAndConv    []*Perm `json:"visitor_and_conv"`
}

func ConvertRolesToAdapterRoles(roles []*models.Role) (resp *RolesResp) {
	resp = &RolesResp{
		Roles: make([]*Role, len(roles)),
	}

	for i, role := range roles {
		resp.Roles[i] = &Role{
			Name:  role.Name,
			Token: role.ID,
		}
	}
	return
}
