package models

import (
	"fmt"
	"strings"
	"time"
)

func UpdateRolePerms(db XODB, roleID string, permIDs []string) error {
	sqlStr := `DELETE FROM custmchat.role_perm WHERE role_id = ?`
	if _, err := db.Exec(sqlStr, roleID); err != nil {
		return err
	}

	if len(permIDs) == 0 {
		return nil
	}

	insert := `INSERT INTO custmchat.role_perm(role_id, perm_id, created_at, updated_at) ` +
		`VALUES %s`

	var placeHolders []string
	var args []interface{}

	now := time.Now().UTC()
	for _, id := range permIDs {
		placeHolders = append(placeHolders, "(?, ?, ?, ?)")
		args = append(args, roleID)
		args = append(args, id)
		args = append(args, now)
		args = append(args, now)
	}

	insert = fmt.Sprintf(insert, strings.Join(placeHolders, ", "))
	if _, err := db.Exec(insert, args...); err != nil {
		return err
	}

	return nil
}

// RoleNamesByRoleIDs
// return map[id]=name
func RoleNamesByRoleIDs(db XODB, roleIDs []string) (names map[string]string, err error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}

	query := `SELECT id, name FROM custmchat.role WHERE id IN (%s)`
	placeHolders := strings.TrimRight(strings.Repeat("?, ", len(roleIDs)), ", ")
	var args []interface{}
	for _, v := range roleIDs {
		args = append(args, v)
	}

	q, err := db.Query(fmt.Sprintf(query, placeHolders), args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	names = make(map[string]string)

	for q.Next() {
		var id, name string
		err = q.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		names[id] = name
	}
	if err = q.Err(); err != nil {
		return
	}

	return
}
