package models

import (
	"fmt"
	"strings"

	"bitbucket.org/forfd/custm-chat/webim/common"
)

var (
	QuickReplyTextContentType    = "text"
	QuickReplyEmojiContentType   = "emoji"
	QuickReplyPictureContentType = "picture"

	QuickReplyAgentCreatorType      = "agent"
	QuickReplyEnterpriseCreatorType = "enterprise"
)

func DeleteQuickReplyItemsByGroupID(db XODB, groupID string) error {
	del := `DELETE FROM custmchat.quickreply_item WHERE quickreply_group_id=?`
	if _, err := db.Exec(del, groupID); err != nil {
		return err
	}

	return nil
}

func InsertQuickReplyGroups(db XODB, groups []*QuickreplyGroup) error {
	if len(groups) == 0 {
		return nil
	}

	insert := `INSERT INTO custmchat.quickreply_group(id, ent_id, title, rank, created_by, creator_type, created_at, updated_at) ` +
		`VALUES %s`

	var args []interface{}
	var placeHolders []string
	for _, group := range groups {
		placeHolders = append(placeHolders, fmt.Sprintf("(%s)", common.GenPlaceHolders(8)))
		args = append(args, group.ID)
		args = append(args, group.EntID)
		args = append(args, group.Title)
		args = append(args, group.Rank)
		args = append(args, group.CreatedBy)
		args = append(args, group.CreatorType)
		args = append(args, group.CreatedAt)
		args = append(args, group.UpdatedAt)
	}

	sqlStr := fmt.Sprintf(insert, strings.Join(placeHolders, ", "))
	_, err := db.Exec(sqlStr, args...)
	return err
}

func InsertQuickReplyItems(db XODB, items []*QuickreplyItem) error {
	if len(items) == 0 {
		return nil
	}

	insert := `INSERT INTO custmchat.quickreply_item(id, quickreply_group_id, title, content, content_type, rich_content, rank, hot_key, created_by, created_at, updated_at) ` +
		`VALUES %s`

	var args []interface{}
	var placeHolders []string
	for _, group := range items {
		placeHolders = append(placeHolders, "(?,?,?,?,?,?,?,?,?,?,?)")
		args = append(args, group.ID)
		args = append(args, group.QuickreplyGroupID)
		args = append(args, group.Title)
		args = append(args, group.Content)
		args = append(args, group.ContentType)
		args = append(args, group.RichContent)
		args = append(args, group.Rank)
		args = append(args, group.HotKey)
		args = append(args, group.CreatedBy)
		args = append(args, group.CreatedAt)
		args = append(args, group.UpdatedAt)
	}

	sqlStr := fmt.Sprintf(insert, strings.Join(placeHolders, ", "))
	_, err := db.Exec(sqlStr, args...)
	return err
}

func QuickReplyGroupsByEntID(db XODB, entID string) (groups []*QuickreplyGroup, err error) {
	const sqlstr = `SELECT ` +
		`id, ent_id, title, rank, created_by, creator_type, created_at, updated_at ` +
		`FROM custmchat.quickreply_group ` +
		`WHERE ent_id = ? ` +
		`ORDER BY rank ASC, created_at DESC`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	for q.Next() {
		qg := QuickreplyGroup{}

		// scan
		err = q.Scan(&qg.ID, &qg.EntID, &qg.Title, &qg.Rank, &qg.CreatedBy, &qg.CreatorType, &qg.CreatedAt, &qg.UpdatedAt)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &qg)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return
}

func QuickReplyItemsByGroupIDs(db XODB, groupIDs []string) (items []*QuickreplyItem, err error) {
	if len(groupIDs) == 0 {
		return []*QuickreplyItem{}, nil
	}

	sqlstr := `SELECT ` +
		`id, quickreply_group_id, title, content, content_type, rich_content, rank, hot_key, created_by, created_at, updated_at ` +
		`FROM custmchat.quickreply_item ` +
		`WHERE quickreply_group_id IN (%s)`

	placeHolders := strings.TrimRight(strings.Repeat("?,", len(groupIDs)), ",")
	sqlstr = fmt.Sprintf(sqlstr, placeHolders)

	var args []interface{}
	for _, id := range groupIDs {
		args = append(args, id)
	}

	q, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	for q.Next() {
		qi := QuickreplyItem{}

		// scan
		err = q.Scan(&qi.ID, &qi.QuickreplyGroupID, &qi.Title, &qi.Content, &qi.ContentType, &qi.RichContent, &qi.Rank, &qi.HotKey, &qi.CreatedBy, &qi.CreatedAt, &qi.UpdatedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, &qi)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return
}

func QuickGroups(db XODB, entID, groupType string) (groups []*QuickreplyGroup, err error) {
	query := `SELECT id, rank FROM custmchat.quickreply_group WHERE ent_id=? AND creator_type=? ORDER BY rank ASC, created_at DESC`

	rows, err := db.Query(query, entID, groupType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		qg := QuickreplyGroup{}

		err = rows.Scan(&qg.ID, &qg.Rank)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &qg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func LastQkGroupRank(db XODB, entID string) (rank int, err error) {
	query := "SELECT rank FROM custmchat.quickreply_group WHERE ent_id=? ORDER BY rank DESC LIMIT 1"
	err = db.QueryRow(query, entID).Scan(&rank)
	return
}

func LastQkReplyRank(db XODB, groupID string) (rank int, err error) {
	query := "SELECT rank FROM custmchat.quickreply_item WHERE quickreply_group_id=? ORDER BY rank DESC LIMIT 1"
	err = db.QueryRow(query, groupID).Scan(&rank)
	return
}
