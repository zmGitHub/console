package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	log "bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
)

// 版本类型定义
const (
	// EditionFree 体验版
	EditionFree                = 1
	EditionFreeAgentServeLimit = 2

	EditionStandard                = 2
	EditionStandardAgentServeLimit = 99

	EditionEnterprise                = 3
	EditionEnterpriseAgentServeLimit = 160
)

// 试用状态定义
const (
	// TrialNone 没有试用
	TrialNone = 1
	// TrialIn 试用中
	TrialIn = 2
	// TrialEnd 已经结束试用
	TrialEnd = 3

	TrialAgentNum = 2
)

var EntFields = `id, name, admin_id, allocation_rule, full_name, nick_name, ` +
	`province, city, avatar, industry, location, ` +
	`address, website, email, mobile, description, created_at, ` +
	`owner, plan, agent_num, trial_status, is_activated, expiration_time, last_activated_at, ` +
	`contact_mobile, contact_email, contact_qq, contact_wechat, contact_signature, contact_name `

func GetEntByEmailORName(db XODB, email, fullName string) (*Enterprise, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, email, mobile ` +
		`FROM custmchat.enterprise ` +
		`WHERE email = ? OR full_name = ?`

	e := Enterprise{}
	err = db.QueryRow(sqlstr, email, fullName).Scan(&e.ID, &e.Name, &e.Email, &e.Mobile)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func UpdateEntTrialStatus(db XODB, id string, trialStatus int) error {
	sqlStr := `UPDATE custmchat.enterprise SET trial_status = ?, agent_num = 1 WHERE id = ?`
	if _, err := db.Exec(sqlStr, trialStatus, id); err != nil {
		return err
	}

	return nil
}

func UpdateEntActivated(db XODB, id string) (err error) {
	update := `UPDATE custmchat.enterprise SET is_activated = ? WHERE id=?`
	_, err = db.Exec(update, true, id)
	return err
}

func SetEntTrial(db XODB, id string, plan, trialStatus, agentNum int, exp time.Time) error {
	sqlStr := `UPDATE custmchat.enterprise SET plan=?, trial_status=?, agent_num=?, expiration_time=? WHERE id = ?`
	if _, err := db.Exec(sqlStr, plan, trialStatus, agentNum, exp, id); err != nil {
		return err
	}

	return nil
}

func GetEntAgentNum(dber XODB, entID string) (agentNum int, err error) {
	entKey := fmt.Sprintf(common.EntConfigs, entID)
	result, err := db.RedisClient.HGet(entKey, common.EntConfigsAgentNum).Result()
	if err != nil || result == "" {
		agentNum, err = getEntAgentNumFromDB(dber, entID)
		if err != nil {
			return
		}

		db.RedisClient.HSet(entKey, common.EntConfigsAgentNum, agentNum)
		return
	}

	return strconv.Atoi(result)
}

func getEntAgentNumFromDB(dber XODB, entID string) (agentNum int, err error) {
	query := "SELECT agent_num FROM custmchat.enterprise WHERE id = ?"
	if err = dber.QueryRow(query, entID).Scan(&agentNum); err != nil {
		return
	}

	return
}

func UpdateEntInfo(db XODB, entID string, values map[string]interface{}) error {
	if len(values) == 0 {
		return nil
	}

	update := `UPDATE custmchat.enterprise SET %s WHERE id = ?`

	var args []interface{}
	var placeHolders []string
	for name, value := range values {
		placeHolders = append(placeHolders, fmt.Sprintf("%s=?", name))
		args = append(args, value)
	}

	args = append(args, entID)
	update = fmt.Sprintf(update, strings.Join(placeHolders, ","))
	if _, err := db.Exec(update, args...); err != nil {
		return err
	}

	return nil
}

func QueryEntTotalCount(db XODB, conds map[string]interface{}) (total int, err error) {
	query := `SELECT COUNT(id) FROM custmchat.enterprise `
	var args []interface{}
	var where []string

	if len(conds) > 0 {
		query += " WHERE %s"
		for name, val := range conds {
			where = append(where, fmt.Sprintf("%s=?", name))
			args = append(args, val)
		}
		query = fmt.Sprintf(query, strings.Join(where, " AND "))
	}

	if err = db.QueryRow(query, args...).Scan(&total); err != nil {
		return
	}

	return
}

func QueryEnterprises(db XODB, offset, limit int, entType, entEmail string) (total int64, result []*Enterprise, err error) {
	countQuery := `SELECT COUNT(id) FROM custmchat.enterprise `
	query := `SELECT ` + EntFields + ` FROM custmchat.enterprise `

	var args []interface{}
	var where string

	switch entType {
	case "paid":
		where = " WHERE plan IN (?,?) AND trial_status=? "
		args = []interface{}{EditionStandard, EditionEnterprise, TrialNone}
	case "unpaid":
		where = " WHERE trial_status IN (?,?) OR plan=? "
		args = []interface{}{TrialIn, TrialEnd, EditionFree}
	default:
	}

	if entEmail != "" {
		if len(args) > 0 {
			where += " AND email = ? "
		} else {
			where += " WHERE email = ? "
		}
		args = append(args, entEmail)
	}

	query += where + ` LIMIT ?,?`
	countQuery += where + ` LIMIT ?,?`
	args = append(args, offset)
	args = append(args, limit)

	log.Logger.Infof("QueryEnterprises Query SQL: %s, args: %v", query, args)
	if err = db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		e := &Enterprise{}
		err = rows.Scan(
			&e.ID, &e.Name, &e.AdminID, &e.AllocationRule, &e.FullName, &e.NickName,
			&e.Province, &e.City, &e.Avatar, &e.Industry, &e.Location, &e.Address, &e.Website,
			&e.Email, &e.Mobile, &e.Description, &e.CreatedAt, &e.Owner, &e.Plan,
			&e.AgentNum, &e.TrialStatus, &e.IsActivated, &e.ExpirationTime, &e.LastActivatedAt,
			&e.ContactMobile, &e.ContactEmail, &e.ContactQq, &e.ContactWechat, &e.ContactSignature, &e.ContactName)
		if err != nil {
			return
		}

		result = append(result, e)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}
