package models

import (
	"fmt"
	"strings"
)

var (
	LeaveMessageHandledStatus   = "close"
	LeaveMessageUnHandledStatus = "open"

	LeaveMessageConfigFillSingleContact = "single"
	LeaveMessageConfigFillMultiContact  = "multi"

	LeaveMessageFields = `id, ent_id, track_id, last_option_agent, mobile, content, status, created_at, updated_at `
)

func UpdateLeaveMessageStatus(db XODB, id, status, agentID string) error {
	update := `UPDATE custmchat.leave_message SET last_option_agent=?, status = ? WHERE id = ?`
	if _, err := db.Exec(update, agentID, status, id); err != nil {
		return err
	}

	return nil
}

func LeaveMessagesByTrackIDs(db XODB, trackIDs []string, offset, limit int) (total int64, lms []*LeaveMessage, err error) {
	if len(trackIDs) == 0 {
		return
	}

	var trackIDIs []interface{}
	var ps []string
	for _, id := range trackIDs {
		trackIDIs = append(trackIDIs, id)
		ps = append(ps, "?")
	}

	var args []interface{}
	args = append(args, trackIDIs...)
	args = append(args, offset)
	args = append(args, limit)

	query := `SELECT ` + LeaveMessageFields +
		` FROM custmchat.leave_message ` +
		`WHERE track_id IN (%s) ` +
		`ORDER BY created_at DESC ` +
		`LIMIT ?,?`
	query = fmt.Sprintf(query, strings.Join(ps, ","))

	q, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer q.Close()

	for q.Next() {
		lm := LeaveMessage{}

		err = q.Scan(&lm.ID, &lm.EntID, &lm.TrackID, &lm.LastOptionAgent, &lm.Mobile, &lm.Content, &lm.Status, &lm.CreatedAt, &lm.UpdatedAt)
		if err != nil {
			return
		}

		lms = append(lms, &lm)
	}
	if err = q.Err(); err != nil {
		return
	}

	query = `SELECT COUNT(id) FROM custmchat.leave_message WHERE track_id IN (%s)`
	if err = db.QueryRow(fmt.Sprintf(query, strings.Join(ps, ",")), trackIDIs...).Scan(&total); err != nil {
		return
	}

	return
}

func EntLeaveMessagesByStatus(db XODB, entID string, offset, limit int) (total int64, lms []*LeaveMessage, err error) {
	query := `SELECT COUNT(id) FROM custmchat.leave_message WHERE ent_id=?`
	if err := db.QueryRow(query, entID).Scan(&total); err != nil {
		return 0, nil, err
	}

	sqlstr := `SELECT ` + LeaveMessageFields +
		`FROM custmchat.leave_message ` +
		`WHERE ent_id = ? ` +
		`ORDER BY created_at DESC ` +
		`LIMIT ?, ?`

	q, err := db.Query(sqlstr, entID, offset, limit)
	if err != nil {
		return
	}
	defer q.Close()

	for q.Next() {
		lm := LeaveMessage{}

		// scan
		err = q.Scan(&lm.ID, &lm.EntID, &lm.TrackID, &lm.LastOptionAgent, &lm.Mobile, &lm.Content, &lm.Status, &lm.CreatedAt, &lm.UpdatedAt)
		if err != nil {
			return
		}

		lms = append(lms, &lm)
	}
	if err = q.Err(); err != nil {
		return
	}

	return
}

func CreateOrUpdateLeaveMessageConfig(db XODB, lmc *LeaveMessageConfig) error {
	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.leave_message_config (` +
		`ent_id, introduction, show_visitor_name, show_telephone, show_email, auto_create_category, fill_contact, use_default_content, default_content` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE ` +
		`introduction = VALUES(introduction), show_visitor_name = VALUES(show_visitor_name), ` +
		`show_telephone = VALUES(show_telephone), show_email = VALUES(show_email), ` +
		`show_wechat = VALUES(show_wechat), show_qq = VALUES(show_qq), ` +
		`auto_create_category = VALUES(auto_create_category), fill_contact = VALUES(fill_contact), ` +
		`use_default_content = VALUES(use_default_content), default_content = VALUES(default_content)`

	// run query
	XOLog(sqlstr, lmc.EntID, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent)
	_, err := db.Exec(sqlstr, lmc.EntID, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent)
	if err != nil {
		return err
	}

	return nil
}
