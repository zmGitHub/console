package models

import "database/sql"

func GetAgentLimitByEntID(db XODB, entID string) (int, error) {
	sqlStr := `SELECT agent_num FROM custmchat.enterprise WHERE id=?`

	var agentLimit int
	if err := db.QueryRow(sqlStr, entID).Scan(&agentLimit); err != nil {
		return -1, err
	}

	return agentLimit, nil
}

func GetAgentServeLimitByEntID(db XODB, entID string) (serveLimit int, err error) {
	sqlStr := `SELECT agent_serve_limit FROM custmchat.ent_plan WHERE ent_id=?`

	if err := db.QueryRow(sqlStr, entID).Scan(&serveLimit); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return -1, err
	}

	return
}
