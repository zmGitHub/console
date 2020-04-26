package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func IncrMessageCount(db XODB, convID, fromType string) error {
	var update string
	switch fromType {
	case MessageFromAgentType:
		update = `UPDATE custmchat.conversation SET agent_msg_count = agent_msg_count + 1, msg_count = msg_count + 1 ` +
			`WHERE id = ?`
	case MessageFromVisitorType:
		update = `UPDATE custmchat.conversation SET client_msg_count = client_msg_count + 1, msg_count = msg_count + 1 ` +
			`WHERE id = ?`
	case MessageMsgSystemType:
		update = `UPDATE custmchat.conversation SET msg_count = msg_count + 1 WHERE id = ?`
	default:
		return nil
	}

	if _, err := db.Exec(update, convID); err != nil {
		return err
	}

	return nil
}

func LastMessageByConversationID(db XODB, conversationID string) (*Message, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, agent_id, conversation_id, from_type, content_type, created_at, read_time, content, msg_type, extra ` +
		`FROM custmchat.message ` +
		`WHERE conversation_id = ? AND from_type IN (?,?) AND msg_type <> ? ` +
		`ORDER BY created_at DESC ` +
		`LIMIT 1`

	m := Message{}
	args := []interface{}{conversationID, MessageFromVisitorType, MessageFromAgentType, MessageMsgSystemType}
	err = db.QueryRow(sqlstr, args...).Scan(&m.ID, &m.EntID, &m.TraceID, &m.AgentID, &m.ConversationID, &m.FromType, &m.ContentType, &m.CreatedAt, &m.ReadTime, &m.Content, &m.MsgType, &m.Extra)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func FirstClientMessageSendTime(db XODB, conversationID string) (t time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE conversation_id=? AND from_type=? ORDER BY created_at ASC LIMIT 1`
	if err = db.QueryRow(query, conversationID, MessageFromVisitorType).Scan(&t); err != nil {
		return
	}

	return
}

func LastClientMessageSendTime(db XODB, conversationID string) (t time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE conversation_id=? AND from_type=? ORDER BY created_at DESC LIMIT 1`
	if err = db.QueryRow(query, conversationID, MessageFromVisitorType).Scan(&t); err != nil {
		return
	}

	return
}

func FirstAgentMsgCreateTime(db XODB, conversationID string) (t time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE conversation_id=? AND from_type=? AND msg_type=? ORDER BY created_at LIMIT 1`
	if err = db.QueryRow(query, conversationID, MessageFromAgentType, MessageMsgPublicType).Scan(&t); err != nil {
		return
	}

	return
}

func FirstConversationMsgCreateTime(db XODB, conversationID string) (t time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE conversation_id=? AND msg_type=? ORDER BY created_at LIMIT 1`
	if err = db.QueryRow(query, conversationID, MessageMsgPublicType).Scan(&t); err != nil {
		return
	}

	return
}

func CreateOrUpdateMessageBeep(db XODB, mb *MessageBeep) error {
	const sqlstr = `INSERT INTO custmchat.message_beep (` +
		`agent_id, client_type, beep_type, new_conversation, new_message, conversation_transfer_in, conversation_transfer_out, colleague_conversation` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE ` +
		`client_type = VALUES(client_type), beep_type = VALUES(beep_type), new_conversation = VALUES(new_conversation), ` +
		`new_message = VALUES(new_message), conversation_transfer_in = VALUES(conversation_transfer_in), ` +
		`conversation_transfer_out = VALUES(conversation_transfer_out), colleague_conversation = VALUES(colleague_conversation)`

	_, err := db.Exec(sqlstr, mb.AgentID, mb.ClientType, mb.BeepType, mb.NewConversation, mb.NewMessage, mb.ConversationTransferIn, mb.ConversationTransferOut, mb.ColleagueConversation)
	return err
}

func LastMessageCreateTimeOfSender(db XODB, convID string, fromType string) (createdAt time.Time, err error) {
	// sql query
	const sqlstr = `SELECT created_at ` +
		`FROM custmchat.message ` +
		`WHERE conversation_id = ? AND from_type = ? ` +
		`ORDER BY created_at DESC ` +
		`LIMIT 1`

	XOLog(sqlstr, convID, fromType)
	q, err := db.Query(sqlstr, convID, fromType)
	if err != nil {
		return
	}
	defer q.Close()

	// load results
	var res []*Message
	for q.Next() {
		m := Message{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.ID, &m.EntID, &m.TraceID, &m.AgentID, &m.ConversationID, &m.FromType, &m.ContentType, &m.CreatedAt, &m.ReadTime, &m.Content, &m.MsgType, &m.Extra)
		if err != nil {
			return
		}

		res = append(res, &m)
	}
	if err = q.Err(); err != nil {
		return
	}

	if len(res) > 0 {
		return res[0].CreatedAt, nil
	}

	return time.Time{}, nil
}

func LastMessageByFromType(db XODB, convID, fromType string) (lastMsgCreateTime time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE conversation_id = ? AND from_type = ? ORDER BY created_at DESC LIMIT 1`

	if err = db.QueryRow(query, convID, fromType).Scan(&lastMsgCreateTime); err != nil {
		return time.Time{}, err
	}

	return
}

func LastMessagebyTraceID(db XODB, traceID string) (msg *Message, err error) {
	query := `SELECT ` +
		`id, ent_id, trace_id, agent_id, conversation_id, from_type, content_type, created_at, read_time, content, msg_type, extra ` +
		`FROM custmchat.message ` +
		`WHERE trace_id = ? ORDER BY created_at DESC LIMIT 1`

	m := &Message{}
	err = db.QueryRow(query, traceID).Scan(
		&m.ID, &m.EntID, &m.TraceID, &m.AgentID, &m.ConversationID, &m.FromType, &m.ContentType, &m.CreatedAt, &m.ReadTime, &m.Content, &m.MsgType, &m.Extra,
	)
	if err != nil {
		return nil, err
	}

	msg = m
	return
}

func VisitorLastMessage(db XODB, traceID string) (lastMsgCreateTime time.Time, err error) {
	query := `SELECT created_at FROM custmchat.message WHERE trace_id = ? AND from_type = ? ORDER BY created_at DESC LIMIT 1`

	if err = db.QueryRow(query, traceID, MessageFromVisitorType).Scan(&lastMsgCreateTime); err != nil {
		return time.Time{}, err
	}

	return
}

func AutomaticMessageByEndIDMsgType(db XODB, entID, msgType string) (msg *AutomaticMessage, err error) {
	const sqlstr = `SELECT ` +
		`id, ent_id, channel_type, msg_type, msg_content, after_seconds, enabled, created_at ` +
		`FROM custmchat.automatic_message ` +
		`WHERE ent_id = ? AND msg_type = ? AND enabled = ?`

	am := AutomaticMessage{}
	err = db.QueryRow(sqlstr, entID, msgType, true).Scan(&am.ID, &am.EntID, &am.ChannelType, &am.MsgType, &am.MsgContent, &am.AfterSeconds, &am.Enabled, &am.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &am, nil
}

func MessagesByConversationIDs(db XODB, conversationIDs []string) (msgs []*Message, err error) {
	if len(conversationIDs) == 0 {
		return []*Message{}, nil
	}

	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, agent_id, conversation_id, from_type, content_type, created_at, read_time, content, msg_type, extra ` +
		`FROM custmchat.message ` +
		`WHERE conversation_id IN (%s)`

	var args []interface{}
	var placeHolders []string
	for _, convID := range conversationIDs {
		args = append(args, convID)
		placeHolders = append(placeHolders, "?")
	}

	rows, err := db.Query(fmt.Sprintf(sqlstr, strings.Join(placeHolders, ",")), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := Message{}

		err = rows.Scan(&m.ID, &m.EntID, &m.TraceID, &m.AgentID, &m.ConversationID, &m.FromType, &m.ContentType, &m.CreatedAt, &m.ReadTime, &m.Content, &m.MsgType, &m.Extra)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

// AgentsByConversationID 参与对话的所有agent
func AgentsByConversationID(db XODB, convID string) (agentIDs []string, err error) {
	query := `SELECT DISTINCT agent_id FROM custmchat.message WHERE conversation_id = ?`
	rows, err := db.Query(query, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return
		}

		agentIDs = append(agentIDs, id)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func MessageCountsByConversationID(db XODB, convID string) (counts map[string]int, err error) {
	query := `SELECT from_type, count(id) AS count FROM  message WHERE conversation_id=? GROUP BY from_type`
	rows, err := db.Query(query, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts = map[string]int{}
	for rows.Next() {
		var fromType string
		var count int

		err = rows.Scan(&fromType, &count)
		if err != nil {
			return
		}

		counts[fromType] = count
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func TimeLimeMessages(db XODB, entID, trackID string, start time.Time) (messages []*Message, err error) {
	query := `SELECT ` +
		`id, ent_id, trace_id, agent_id, conversation_id, from_type, content_type, created_at, read_time, content, msg_type, extra ` +
		`FROM custmchat.message ` +
		`WHERE ent_id = ? AND trace_id=? AND created_at >= ?`

	rows, err := db.Query(query, entID, trackID, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m := Message{}

		err = rows.Scan(&m.ID, &m.EntID, &m.TraceID, &m.AgentID, &m.ConversationID, &m.FromType, &m.ContentType, &m.CreatedAt, &m.ReadTime, &m.Content, &m.MsgType, &m.Extra)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}
