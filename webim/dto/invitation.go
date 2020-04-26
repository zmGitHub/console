package dto

import (
	"strings"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type AgentInvitation struct {
	Email          string      `json:"email"`
	GroupID        string      `json:"group_id"`
	Privilege      string      `json:"privilege"`       // role token
	PrivilegeRange interface{} `json:"privilege_range"` // self, admin, group_list: []
	Realname       string      `json:"realname"`
	ServingLimit   int         `json:"serving_limit"`
	WorkNum        string      `json:"work_num"`
}

type InvitationReq struct {
	Agents    []*AgentInvitation `json:"agents"`
	BrowserID string             `json:"browser_id"`
}

type Invitation struct {
	AcceptedOn     *string     `json:"accepted_on"`
	CreatedOn      string      `json:"created_on"`
	Email          string      `json:"email"`
	EnterpriseID   string      `json:"enterprise_id"`
	ExpiredOn      string      `json:"expired_on"`
	GroupID        string      `json:"group_id"`
	ID             string      `json:"id"`
	LastUpdated    string      `json:"last_updated"`
	Privilege      string      `json:"privilege"`
	PrivilegeRange interface{} `json:"privilege_range"`
	Realname       string      `json:"realname"`
	ServingLimit   int         `json:"serving_limit"`
	Status         string      `json:"status"`
	WorkNum        string      `json:"work_num"`
}

func ConvertInvitation(invt *models.AgentInvitation) (*Invitation, error) {
	v := &Invitation{
		ID:           invt.ID,
		CreatedOn:    *common.ConvertUTCToTimeString(invt.CreatedOn),
		Email:        invt.Email,
		EnterpriseID: invt.EnterpriseID,
		ExpiredOn:    *common.ConvertUTCToTimeString(invt.ExpiredOn),
		GroupID:      invt.GroupID,
		LastUpdated:  *common.ConvertUTCToTimeString(invt.LastUpdated),
		Privilege:    invt.Privilege,
		Realname:     invt.Realname,
		ServingLimit: invt.ServingLimit,
		Status:       invt.Status,
		WorkNum:      invt.WorkNum,
	}

	if invt.AcceptedOn.Valid {
		v.AcceptedOn = common.ConvertUTCToTimeString(invt.AcceptedOn.Time)
	}

	if strings.HasPrefix(invt.PrivilegeRange, "[") {
		var groups []string
		if err := common.Unmarshal(invt.PrivilegeRange, &groups); err != nil {
			return nil, err
		}
		v.PrivilegeRange = groups
	} else {
		v.PrivilegeRange = invt.PrivilegeRange
	}

	return v, nil
}
