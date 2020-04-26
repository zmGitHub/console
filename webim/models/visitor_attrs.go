package models

import (
	"database/sql"
	"fmt"

	"bitbucket.org/forfd/custm-chat/webim/common"
)

// visitor_attrs
// `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
//    `trace_id` CHAR(20) NOT NULL COMMENT 'trace id',
//    `attrs` TEXT NOT NULL COMMENT 'attrs',
type VisitorAttrs struct {
	EntID   string         `json:"ent_id"`
	TraceID string         `json:"trace_id"`
	Attrs   sql.NullString `json:"attrs"`
}

func GetVisitorAttrs(db XODB, entID, trackID string) (attrs map[string]interface{}, err error) {
	query := `SELECT attrs FROM custmchat.visitor_attrs WHERE ent_id=? AND trace_id=?`

	var attrsValue sql.NullString
	err = db.QueryRow(query, entID, trackID).Scan(&attrsValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return map[string]interface{}{}, nil
		}

		return nil, err
	}

	if !attrsValue.Valid {
		return nil, fmt.Errorf("record not found")
	}

	err = common.Unmarshal(attrsValue.String, &attrs)
	return
}

func UpdateVisitorAttrs(db XODB, entID, trackID string, attrs map[string]interface{}) error {
	update := `INSERT INTO custmchat.visitor_attrs (ent_id, trace_id, attrs) VALUES (?, ?, ?)  ON DUPLICATE KEY UPDATE attrs=VALUES(attrs)`
	attrsValue, err := common.Marshal(attrs)
	if err != nil {
		return err
	}
	_, err = db.Exec(update, entID, trackID, attrsValue)
	return err
}
