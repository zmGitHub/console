package adapter

import (
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type Agent struct {
	ID              string      `json:"id"`
	Avatar          string      `json:"avatar"`
	Cellphone       string      `json:"cellphone"`
	CreatedOn       *string     `json:"created_on"`
	DeletedAt       *time.Time  `json:"deleted_at"`
	Email           string      `json:"email"`
	EmailActivated  bool        `json:"email_activated"`
	EnterpriseID    string      `json:"enterprise_id"`
	GroupID         string      `json:"group_id"`
	IsOnline        bool        `json:"is_online"`
	Nickname        string      `json:"nickname"`
	Privilege       string      `json:"privilege"`
	PrivilegeRange  interface{} `json:"privilege_range"`
	PublicCellphone string      `json:"public_cellphone"`
	PublicEmail     string      `json:"public_email"`
	Qq              string      `json:"qq"`
	Rank            int         `json:"rank"`
	ReadFeatureID   int         `json:"read_feature_id"`
	Realname        string      `json:"realname"`
	ServingLimit    int         `json:"serving_limit"`
	Signature       string      `json:"signature"`
	Status          string      `json:"status"`
	Telephone       string      `json:"telephone"`
	Token           string      `json:"token"`
	Weixin          string      `json:"weixin"`
	WorkNum         string      `json:"work_num"`
}

type AgentResp map[string]*Agent

type AgentGroup struct {
	EnterpriseID string `json:"enterprise_id"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Token        string `json:"token"`
}

type AgentGroupResp struct {
	AgentGroups []*AgentGroup `json:"agent_groups"`
}

func ConvertAgentToAgentInfo(agent *models.Agent, queryGroups bool) *Agent {
	if agent == nil {
		return nil
	}

	var deleteAt *time.Time
	if agent.DeletedAt.Valid {
		deleteAt = &agent.DeletedAt.Time
	}

	var emailActivated bool
	if agent.AccountStatus == models.AgentAccountValidStatus {
		emailActivated = true
	}

	var permsRange interface{}
	if agent.PermsRangeType == models.AgentPermsRangePersonalType || agent.PermsRangeType == models.AgentPermsRangeAllType {
		permsRange = agent.PermsRangeType
	} else {
		if queryGroups {
			groups, err := models.PermsRangeGroupIDsByAgentID(db.Mysql, agent.ID)
			if err != nil {
				log.Logger.Warnf("[PermsRangeGroupIDsByAgentID] agentID: %s, error: %v", agent.ID, err)
			}
			permsRange = groups
		}
	}

	return &Agent{
		ID:              agent.ID,
		Avatar:          agent.Avatar,
		Cellphone:       agent.Mobile,
		CreatedOn:       common.ConvertUTCToTimeString(agent.CreateAt),
		DeletedAt:       deleteAt,
		Email:           agent.Email,
		EmailActivated:  emailActivated,
		EnterpriseID:    agent.EntID,
		GroupID:         agent.GroupID,
		IsOnline:        false,
		Nickname:        agent.NickName,
		Privilege:       agent.RoleID,
		PrivilegeRange:  permsRange,
		PublicCellphone: agent.PublicTelephone,
		PublicEmail:     agent.PublicEmail,
		Qq:              agent.QqNum,
		Rank:            agent.Ranking,
		ReadFeatureID:   0,
		Realname:        agent.RealName,
		ServingLimit:    agent.ServeLimit,
		Signature:       agent.Signature,
		Status:          agent.Status,
		Telephone:       agent.Mobile,
		Token:           agent.ID,
		Weixin:          agent.Wechat,
		WorkNum:         agent.JobNumber,
	}
}

func ConvertAgentsToAdapterAgents(agents []*models.Agent) AgentResp {
	resp := make(AgentResp, len(agents))
	for _, agent := range agents {
		resp[agent.ID] = ConvertAgentToAgentInfo(agent, false)
	}

	return resp
}

func ConvertAgentsToAdapterAgentsV1(agents []*models.Agent, agentGroups map[string][]string) []*Agent {
	resp := []*Agent{}
	for _, agent := range agents {
		v := ConvertAgentToAgentInfo(agent, false)
		if v.PrivilegeRange == nil {
			if groups, ok := agentGroups[agent.ID]; ok {
				v.PrivilegeRange = groups
			}
		}

		resp = append(resp, v)
	}

	return resp
}

func ConvertAgentGroupsToAdapterGroups(groups []*models.AgentGroup) *AgentGroupResp {
	resp := &AgentGroupResp{}
	var resultGroups = make([]*AgentGroup, len(groups))

	for i, group := range groups {
		resultGroups[i] = &AgentGroup{
			EnterpriseID: group.EntID,
			ID:           group.ID,
			Name:         group.Name,
			Token:        group.ID,
		}
	}
	resp.AgentGroups = resultGroups

	return resp
}
