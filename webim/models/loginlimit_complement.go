package models

import (
	"database/sql"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/db"
)

type EntLoginLimit struct {
	EntID         string   `json:"ent_id"`
	Status        bool     `json:"status"`
	GroupIDs      []string `json:"role_ids"`
	CityList      []string `json:"city_list"`
	AllowedIPList []string `json:"allowed_ip_list"`
}

// buildLoginLimitModel ...
func buildLoginLimitModel(limit *EntLoginLimit) (ll *LoginLimit, err error) {
	ll = &LoginLimit{
		EntID:  limit.EntID,
		Status: limit.Status,
	}

	var groupIDs = `[]`
	if len(limit.GroupIDs) > 0 {
		groupIDs, err = common.Marshal(limit.GroupIDs)
		if err != nil {
			return
		}
	}

	var cityList = `[]`
	if len(limit.CityList) > 0 {
		cityList, err = common.Marshal(limit.CityList)
		if err != nil {
			return
		}
	}

	var ipList = `[]`
	if len(limit.AllowedIPList) > 0 {
		cityList, err = common.Marshal(limit.AllowedIPList)
		if err != nil {
			return
		}
	}

	ll.GroupIds = sql.NullString{String: groupIDs, Valid: true}
	ll.CityList = sql.NullString{String: cityList, Valid: true}
	ll.AllowedIPList = sql.NullString{String: ipList, Valid: true}
	return
}

func CreateLoginLimit(db XODB, limit *EntLoginLimit) error {
	ll, err := buildLoginLimitModel(limit)
	if err != nil {
		return err
	}

	return ll.Insert(db)
}

func UpdateEntLoginLimit(db XODB, limit *EntLoginLimit) error {
	ll, err := buildLoginLimitModel(limit)
	if err != nil {
		return err
	}

	return ll.Update(db)
}

func EntLoginLimitByEntID(entID string) (limit *EntLoginLimit, err error) {
	limit = &EntLoginLimit{
		EntID:         entID,
		Status:        false,
		GroupIDs:      []string{},
		CityList:      []string{},
		AllowedIPList: []string{},
	}

	ll, err := LoginLimitByEntID(db.Mysql, entID)
	if err != nil {
		if err == sql.ErrNoRows {
			return limit, nil
		}

		return nil, err
	}

	limit.Status = ll.Status
	err = common.Unmarshal(ll.GroupIds.String, limit.GroupIDs)
	if err != nil {
		return nil, err
	}

	err = common.Unmarshal(ll.GroupIds.String, limit.GroupIDs)
	if err != nil {
		return nil, err
	}

	err = common.Unmarshal(ll.GroupIds.String, limit.GroupIDs)
	if err != nil {
		return nil, err
	}

	return
}
