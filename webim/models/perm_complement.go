package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	db2 "bitbucket.org/forfd/custm-chat/webim/db"
)

func GetAgentPerms(db XODB, agentID string) (perms []*Perm, err error) {
	sqlStr := `SELECT P.id, P.ent_id, P.app_name, P.name, P.created_at, P.updated_at FROM ` +
		`custmchat.perm P ` +
		`INNER JOIN custmchat.role_perm RP ON RP.perm_id = P.id ` +
		`INNER JOIN custmchat.agent AG ON AG.role_id = RP.role_id ` +
		`WHERE AG.id = ?`

	rows, err := db.Query(sqlStr, agentID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log.Logger.Warn("close rows error: ", e)
		}
	}()

	var res []*Perm
	for rows.Next() {
		p := Perm{}

		if err = rows.Scan(&p.ID, &p.EntID, &p.AppName, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}

		res = append(res, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func GetAgentPermFromCache(mysql XODB, entID, agentID, appName, permName string) (hasPerm bool, err error) {
	key := fmt.Sprintf(common.AgentPerms, agentID)

	getPermFromDBAndSetCache := func() (bool, error) {
		hasPerm, err = AgentHasPerm(mysql, entID, agentID, appName, permName)
		if err != nil {
			return false, err
		}

		if _, err = db2.RedisClient.HSet(key, permName, hasPerm).Result(); err != nil {
			log.Logger.Warnf("RedisClient Set Perms(%s) error: %v", key, err)
		}

		return hasPerm, err
	}

	res, err := db2.RedisClient.HGet(key, permName).Result()
	if err != nil {
		log.Logger.Warnf("RedisClient Get Perms(%s) error: %v", key, err)
		hasPerm, err = getPermFromDBAndSetCache()
		if err != nil {
			return
		}

		return hasPerm, nil
	}
	if res == "" {
		return getPermFromDBAndSetCache()
	}

	return res == "1", nil
}

func AgentHasPerm(db XODB, entID, agentID, appName, permName string) (hasPerm bool, err error) {
	isAdmin, err := IsAdmin(db, agentID)
	if err == nil && isAdmin {
		return true, nil
	}

	sqlStr := `SELECT P.id FROM custmchat.perm P ` +
		`INNER JOIN custmchat.role_perm RP ON RP.perm_id = P.id ` +
		`INNER JOIN custmchat.agent AG ON AG.role_id = RP.role_id ` +
		`WHERE AG.id = ? AND P.ent_id=? AND P.app_name=? AND P.name = ?`

	var id string
	err = db.QueryRow(sqlStr, agentID, entID, appName, permName).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return id != "", nil
}

func BulkCreateEntApps(db XODB, entID string, apps []string) error {
	if len(apps) == 0 {
		return nil
	}

	insert := `INSERT INTO custmchat.ent_app(id, ent_id, app_name, create_at, update_at) ` +
		`VALUES %s`
	now := time.Now().UTC()
	var placeHolders []string
	var args []interface{}
	for _, app := range apps {
		placeHolders = append(placeHolders, "(?,?,?,?,?)")
		args = append(args, common.GenUniqueID())
		args = append(args, entID)
		args = append(args, app)
		args = append(args, now)
		args = append(args, now)
	}

	if _, err := db.Exec(fmt.Sprintf(insert, strings.Join(placeHolders, ", ")), args...); err != nil {
		return err
	}

	return nil
}

func BulkCreateEntPerms(db XODB, entID string, perms map[string][]string) ([]string, error) {
	if len(perms) == 0 {
		return nil, nil
	}

	var permIDs []string

	now := time.Now().UTC()
	insert := `INSERT INTO custmchat.perm(id, ent_id, app_name, name, created_at, updated_at) VALUES %s`
	var placeHolders []string
	var args []interface{}

	for appName, funcNames := range perms {
		for _, name := range funcNames {
			placeHolders = append(placeHolders, "(?, ?, ?, ?, ?, ?)")

			permID := common.GenUniqueID()
			permIDs = append(permIDs, permID)

			args = append(args, permID)
			args = append(args, entID)
			args = append(args, appName)
			args = append(args, name)
			args = append(args, now)
			args = append(args, now)
		}
	}

	if _, err := db.Exec(fmt.Sprintf(insert, strings.Join(placeHolders, ", ")), args...); err != nil {
		return nil, err
	}

	return permIDs, nil
}

func BulkAddAgentPermsRangeGroups(db XODB, agentID string, groups []string) error {
	if len(groups) == 0 {
		return nil
	}

	insert := `INSERT INTO custmchat.perms_range_groups(agent_id, group_id) VALUES %s`
	var placeHolders []string
	var args []interface{}

	for _, id := range groups {
		placeHolders = append(placeHolders, "(?, ?)")
		args = append(args, agentID)
		args = append(args, id)
	}

	if _, err := db.Exec(fmt.Sprintf(insert, strings.Join(placeHolders, ", ")), args...); err != nil {
		return err
	}

	return nil
}

func UpdateAgentPermGroups(db XODB, agentID string, newGroups []string) error {
	del := `DELETE FROM custmchat.perms_range_groups WHERE agent_id = ?`
	if _, err := db.Exec(del, agentID); err != nil {
		return err
	}

	return BulkAddAgentPermsRangeGroups(db, agentID, newGroups)
}

func PermsRangeGroupIDsByAgentID(db XODB, agentID string) (groups []string, err error) {
	const sqlstr = `SELECT ` +
		`group_id ` +
		`FROM custmchat.perms_range_groups ` +
		`WHERE agent_id = ?`

	rows, err := db.Query(sqlstr, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		groups = append(groups, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func AgentIDsByPermGroupIDs(db XODB, groupIDs []string) (agentIDs []string, err error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	query := `SELECT ` +
		`DISTINCT agent_id ` +
		`FROM custmchat.perms_range_groups ` +
		`WHERE group_id IN (%s)`

	var args []interface{}
	var placeHolders []string
	for _, id := range groupIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	query = fmt.Sprintf(query, strings.Join(placeHolders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		agentIDs = append(agentIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func PermsRangeGroupIDsByAgents(db XODB, agents []*Agent) (agentGroups map[string][]string, err error) {
	var agentIDs []string
	for _, agent := range agents {
		agentIDs = append(agentIDs, agent.ID)
	}

	return PermsRangeGroupIDsByAgentIDs(db, agentIDs)
}
func PermsRangeGroupIDsByAgentIDs(db XODB, agentIDs []string) (agentGroups map[string][]string, err error) {
	if len(agentIDs) == 0 {
		return map[string][]string{}, nil
	}

	query := `SELECT ` +
		`agent_id, group_id ` +
		`FROM custmchat.perms_range_groups ` +
		`WHERE agent_id IN (%s)`

	var args []interface{}
	var placeHolders []string
	for _, id := range agentIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	query = fmt.Sprintf(query, strings.Join(placeHolders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agentGroups = make(map[string][]string)
	for rows.Next() {
		var agentID, groupID string
		if err = rows.Scan(&agentID, &groupID); err != nil {
			return nil, err
		}

		_, ok := agentGroups[agentID]
		if ok {
			agentGroups[agentID] = append(agentGroups[agentID], groupID)
			continue
		}

		agentGroups[agentID] = []string{groupID}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func PermsByPermIDs(db XODB, permIDs []string) (perms []*Perm, err error) {
	if len(permIDs) == 0 {
		return []*Perm{}, nil
	}

	sqlStr := `SELECT P.id, P.ent_id, P.app_name, P.name, P.created_at, P.updated_at FROM custmchat.perm P WHERE P.id IN (%s)`
	var args []interface{}
	var placeHolders []string
	for _, id := range permIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	rows, err := db.Query(fmt.Sprintf(sqlStr, strings.Join(placeHolders, ",")), args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log.Logger.Warnf("close rows error: %v", e)
		}
	}()

	var res []*Perm
	for rows.Next() {
		p := Perm{}

		if err = rows.Scan(&p.ID, &p.EntID, &p.AppName, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}

		res = append(res, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
