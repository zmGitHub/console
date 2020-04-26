package adapter

import (
	"bitbucket.org/forfd/custm-chat/webim/models"
)

type Perm struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type Perms map[string][]*Perm

type PermsRangeResp struct {
	Perms  Perms       `json:"perms"`
	Ranges interface{} `json:"ranges"`
}

type EntPerms map[string][]*Perm

func ConvertAgentPermsToPermsRangesResp(permsRange interface{}, perms []*models.Perm) (resp *PermsRangeResp) {
	resp = &PermsRangeResp{
		Ranges: permsRange,
		Perms:  make(Perms),
	}

	for _, perm := range perms {
		ps, ok := resp.Perms[perm.AppName]
		if !ok {
			resp.Perms[perm.AppName] = []*Perm{
				{
					ID:  perm.ID,
					Key: perm.Name,
				},
			}
			continue
		}

		resp.Perms[perm.AppName] = append(ps, &Perm{ID: perm.ID, Key: perm.Name})
	}
	return
}

func ConvertEntPermsToPerms(perms map[string][]*models.Perm) EntPerms {
	result := make(EntPerms, len(perms))
	for appName, ps := range perms {
		var resp []*Perm
		for _, p := range ps {
			resp = append(resp, &Perm{
				ID:  p.ID,
				Key: p.Name,
			})
		}

		result[appName] = resp
	}
	return result
}
