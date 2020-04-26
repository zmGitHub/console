package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/go-sql-driver/mysql"
)

const (
	FirstLevel  = "first_level"
	SecondLevel = "second_level"
	ThirdLevel  = "third_level"
)

var (
	ConversationHumanAgentType = "human"

	GoodConversation   = 2
	MediumConversation = 1
	BadConversation    = 0
	NoEvalConversation = -1

	EvalLevels = mapset.NewSetFromSlice([]interface{}{0, 1, 2})

	ConversationFields = `id, ent_id, trace_id, agent_id, agent_msg_count, agent_type, msg_count, ` +
		`title, client_first_send_time, client_msg_count, duration, first_msg_created_at, first_response_wait_time, last_msg_content, last_msg_created_at, ` +
		`quality_grade, summary, created_at, update_at, ended_at, ended_by, ` +
		`agent_effective_msg_count, client_last_send_time, first_msg_create_time, eval_content, eval_level, has_summary `
)

func UpdateConversationAgentID(db XODB, id, targetID string) error {
	sqlStr := `UPDATE custmchat.conversation SET ` +
		`agent_id = ? WHERE id = ?`

	_, err := db.Exec(sqlStr, targetID, id)
	return err
}

func UpdateConversationSummary(db XODB, id, summary string) error {
	sqlStr := `UPDATE custmchat.conversation SET ` +
		`summary = ?, has_summary = ? WHERE id = ?`

	_, err := db.Exec(sqlStr, summary, true, id)
	return err
}

func CreateConversationEvaluation(db XODB, eval *Evaluation) (err error) {
	const sqlstr = `INSERT INTO custmchat.evaluation (` +
		`ent_id, agent_id, conv_id, level, content, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE ` +
		`level = VALUES(level), content = VALUES(content)`

	XOLog(sqlstr, eval.EntID, eval.AgentID, eval.ConvID, eval.Level, eval.Content, eval.CreatedAt)
	_, err = db.Exec(sqlstr, eval.EntID, eval.AgentID, eval.ConvID, eval.Level, eval.Content, eval.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func CreateConversationEvalV1(db XODB, id string, evalContent string, evalLEvel int) error {
	update := `UPDATE custmchat.conversation SET eval_content=?, eval_level=? WHERE id = ?`
	if _, err := db.Exec(update, evalContent, evalLEvel, id); err != nil {
		return err
	}

	return nil
}

func ConversationNumOfAgent(db XODB, agentID string) (count int, err error) {
	query := `SELECT COUNT(id) FROM custmchat.conversation WHERE agent_id=? AND ended_at IS NULL`
	if err = db.QueryRow(query, agentID).Scan(&count); err != nil {
		return 0, err
	}

	return
}

func ConversationsByEntIDTraceID(db XODB, entID, traceID string, offset, limit int) ([]*Conversation, error) {
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, agent_id, agent_msg_count, agent_type, msg_count, title, client_first_send_time, client_msg_count, duration, first_msg_created_at, first_response_wait_time, last_msg_content, last_msg_created_at, quality_grade, summary, created_at, update_at, ended_at, ended_by, agent_effective_msg_count, client_last_send_time, first_msg_create_time, eval_content, eval_level, has_summary ` +
		`FROM custmchat.conversation ` +
		`WHERE ent_id = ? AND trace_id = ? ` +
		`ORDER BY created_at DESC ` +
		`LIMIT ?,?`

	// run query
	XOLog(sqlstr, entID, traceID, offset, limit)
	q, err := db.Query(sqlstr, entID, traceID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var res []*Conversation
	for q.Next() {
		c := Conversation{}

		// scan
		err = q.Scan(&c.ID, &c.EntID, &c.TraceID, &c.AgentID, &c.AgentMsgCount, &c.AgentType, &c.MsgCount, &c.Title, &c.ClientFirstSendTime, &c.ClientMsgCount, &c.Duration, &c.FirstMsgCreatedAt, &c.FirstResponseWaitTime, &c.LastMsgContent, &c.LastMsgCreatedAt, &c.QualityGrade, &c.Summary, &c.CreatedAt, &c.UpdateAt, &c.EndedAt, &c.EndedBy, &c.AgentEffectiveMsgCount, &c.ClientLastSendTime, &c.FirstMsgCreateTime, &c.EvalContent, &c.EvalLevel, &c.HasSummary)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func ConversationsByTraceID(db XODB, traceID string, offset, limit int, active bool) ([]*Conversation, error) {
	sqlstr := `SELECT ` + ConversationFields +
		` FROM custmchat.conversation ` +
		`WHERE trace_id = ? %s ` +
		`ORDER BY created_at DESC ` +
		`LIMIT ?,?`

	if active {
		sqlstr = fmt.Sprintf(sqlstr, " AND ended_at IS NULL ")
	} else {
		sqlstr = fmt.Sprintf(sqlstr, "")
	}

	q, err := db.Query(sqlstr, traceID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var res []*Conversation
	for q.Next() {
		c := Conversation{}

		// scan
		err = q.Scan(
			&c.ID, &c.EntID, &c.TraceID, &c.AgentID, &c.AgentMsgCount, &c.AgentType, &c.MsgCount,
			&c.Title, &c.ClientFirstSendTime, &c.ClientMsgCount, &c.Duration, &c.FirstMsgCreatedAt,
			&c.FirstResponseWaitTime, &c.LastMsgContent, &c.LastMsgCreatedAt, &c.QualityGrade, &c.Summary,
			&c.CreatedAt, &c.UpdateAt, &c.EndedAt, &c.EndedBy, &c.AgentEffectiveMsgCount, &c.ClientLastSendTime, &c.FirstMsgCreateTime,
			&c.EvalContent, &c.EvalLevel, &c.HasSummary)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func IsEffectiveConversation(db XODB, convID string) (effective bool, err error) {
	sqlStr := `SELECT COUNT(id) FROM custmchat.message WHERE conversation_id = ? AND from_type = ?`

	var count int64
	if err = db.QueryRow(sqlStr, convID, MessageFromVisitorType).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return count >= 1, nil
}

func ConvNumByAgentID(db XODB, agentID string) (int, error) {
	sqlStr := `SELECT COUNT(id) AS conv_num FROM custmchat.conversation WHERE agent_id=? AND ended_at IS NULL`

	var convCount int
	if err := db.QueryRow(sqlStr, agentID).Scan(&convCount); err != nil {
		return -1, err
	}
	return convCount, nil
}

func ActiveConversationNum(db XODB, agentID string) (count int, err error) {
	query := `SELECT COUNT(id) FROM custmchat.conversation WHERE agent_id = ? AND ended_at IS NULL`
	err = db.QueryRow(query, agentID).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

func ActiveConversationNumOfAgents(db XODB, agentIDs []string) (counts map[string]int, err error) {
	if len(agentIDs) == 0 {
		return
	}

	var args []interface{}
	var ps []string
	for _, id := range agentIDs {
		args = append(args, id)
		ps = append(ps, "?")
	}

	query := `SELECT agent_id, COUNT(id) AS conv_count FROM custmchat.conversation ` +
		`WHERE agent_id IN (%s)  AND ended_at IS NULL GROUP BY agent_id HAVING conv_count > 0`
	query = fmt.Sprintf(query, strings.Join(ps, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts = map[string]int{}
	for rows.Next() {
		var agentID string
		var count int
		if err = rows.Scan(&agentID, &count); err != nil {
			return
		}

		counts[agentID] = count
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func ActiveConversationsByAgentID(db XODB, agentID string) ([]*Conversation, error) {
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, agent_id, agent_msg_count, agent_type, msg_count, title, client_first_send_time, client_msg_count, duration, first_msg_created_at, first_response_wait_time, last_msg_content, last_msg_created_at, quality_grade, summary, created_at, update_at, ended_at, ended_by ` +
		`FROM custmchat.conversation ` +
		`WHERE agent_id = ? AND ended_at IS NULL ` +
		`ORDER BY created_at DESC`

	// run query
	XOLog(sqlstr, agentID)
	q, err := db.Query(sqlstr, agentID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Conversation
	for q.Next() {
		c := Conversation{}

		// scan
		err = q.Scan(&c.ID, &c.EntID, &c.TraceID, &c.AgentID, &c.AgentMsgCount, &c.AgentType, &c.MsgCount, &c.Title, &c.ClientFirstSendTime, &c.ClientMsgCount, &c.Duration, &c.FirstMsgCreatedAt, &c.FirstResponseWaitTime, &c.LastMsgContent, &c.LastMsgCreatedAt, &c.QualityGrade, &c.Summary, &c.CreatedAt, &c.UpdateAt, &c.EndedAt, &c.EndedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func EndConversation(
	db XODB, id, endBy, lastMsgContent string, lastMsgTime time.Time,
	duration, respWaitTime int64, qualityLevel string,
	firstMsgCreateTime, firstClientSendMsgTime, clientLastSendTime time.Time) error {
	fields := map[string]interface{}{
		"duration":                 duration,
		"ended_at":                 time.Now().UTC(),
		"ended_by":                 endBy,
		"last_msg_content":         lastMsgContent,
		"first_response_wait_time": respWaitTime,
	}

	if !lastMsgTime.IsZero() {
		fields["last_msg_created_at"] = lastMsgTime
	}

	if !firstMsgCreateTime.IsZero() {
		fields["first_msg_created_at"] = firstMsgCreateTime
		fields["first_msg_create_time"] = firstMsgCreateTime
	}

	if !firstClientSendMsgTime.IsZero() {
		fields["client_first_send_time"] = firstClientSendMsgTime
	}

	if !clientLastSendTime.IsZero() {
		fields["client_last_send_time"] = clientLastSendTime
	}

	if qualityLevel != "" {
		fields["quality_grade"] = qualityLevel
	}

	update := `UPDATE custmchat.conversation SET %s WHERE id=?`

	var placeHolders []string
	var args []interface{}
	for name, value := range fields {
		placeHolders = append(placeHolders, fmt.Sprintf("%s=?", name))
		args = append(args, value)
	}
	args = append(args, id)

	_, err := db.Exec(fmt.Sprintf(update, strings.Join(placeHolders, ",")), args...)
	return err
}

func IsConversationEnd(db XODB, id string) (end bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM custmchat.conversation WHERE id=? AND ended_at IS NOT NULL)`
	if err = db.QueryRow(query, id).Scan(&end); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return
	}

	return end, nil
}

func ConversationCreatedAtByID(db XODB, id string) (createTime time.Time, endTime mysql.NullTime, err error) {
	query := `SELECT created_at, ended_at FROM custmchat.conversation WHERE id = ?`
	if err = db.QueryRow(query, id).Scan(&createTime, &endTime); err != nil {
		return time.Time{}, mysql.NullTime{}, err
	}

	return
}

func TrackIDsByAgentID(db XODB, agentID string) (trackIDs []string, err error) {
	query := `SELECT DISTINCT(trace_id) FROM custmchat.conversation WHERE agent_id=?`
	rows, err := db.Query(query, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		trackIDs = append(trackIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func TrackIDsByAgentIDs(db XODB, agentIDs []string) (trackIDs []string, err error) {
	if len(agentIDs) == 0 {
		return []string{}, nil
	}

	query := `SELECT DISTINCT(trace_id) FROM custmchat.conversation WHERE agent_id IN (%s)`

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

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		trackIDs = append(trackIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func AllTrackIDs(db XODB, entID string) (trackIDs []string, err error) {
	query := `SELECT DISTINCT(trace_id) FROM custmchat.conversation WHERE ent_id=?`
	rows, err := db.Query(query, entID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		trackIDs = append(trackIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func ColleagueConversationsByEntIDAgentID(db XODB, entID, agentID string, offset, limit int) ([]*Conversation, bool, error) {
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, agent_id, agent_msg_count, agent_type, msg_count, title, client_first_send_time, client_msg_count, duration, first_msg_created_at, first_response_wait_time, last_msg_content, last_msg_created_at, quality_grade, summary, created_at, update_at, ended_at, ended_by ` +
		`FROM custmchat.conversation ` +
		`WHERE ent_id = ? AND agent_id <> ? AND ended_at IS NULL ` +
		`ORDER BY created_at DESC ` +
		`LIMIT ?,?`

	// run query
	XOLog(sqlstr, entID, agentID, offset, limit)
	q, err := db.Query(sqlstr, entID, agentID, offset, limit+1)
	if err != nil {
		return nil, false, err
	}
	defer q.Close()

	// load results
	var res []*Conversation
	for q.Next() {
		c := Conversation{}

		// scan
		err = q.Scan(&c.ID, &c.EntID, &c.TraceID, &c.AgentID, &c.AgentMsgCount, &c.AgentType, &c.MsgCount, &c.Title, &c.ClientFirstSendTime, &c.ClientMsgCount, &c.Duration, &c.FirstMsgCreatedAt, &c.FirstResponseWaitTime, &c.LastMsgContent, &c.LastMsgCreatedAt, &c.QualityGrade, &c.Summary, &c.CreatedAt, &c.UpdateAt, &c.EndedAt, &c.EndedBy)
		if err != nil {
			return nil, false, err
		}

		res = append(res, &c)
	}
	if err = q.Err(); err != nil {
		return nil, false, err
	}

	var hasNext bool
	if len(res) > limit {
		hasNext = true
		res = res[:len(res)-1]
	}

	return res, hasNext, nil
}

func ConvNumDurationByEntID(db XODB, entID string) (num int64, duration int64, err error) {
	query := `SELECT COUNT(id) FROM custmchat.conversation WHERE ent_id = ? AND ended_at IS NOT NULL`
	if err = db.QueryRow(query, entID).Scan(&num); err != nil {
		return
	}

	query = `SELECT SUM(duration) FROM custmchat.conversation WHERE ent_id = ? AND ended_at IS NOT NULL`
	if err = db.QueryRow(query, entID).Scan(&duration); err != nil {
		return
	}

	return
}

func ConvNumDurationByAgentID(db XODB, agentID string) (num int64, duration, waitRespDuration int64, err error) {
	query := `SELECT COUNT(id) FROM custmchat.conversation WHERE agent_id = ? AND ended_at IS NOT NULL`
	if err = db.QueryRow(query, agentID).Scan(&num); err != nil {
		return
	}

	query = `SELECT SUM(duration) FROM custmchat.conversation WHERE agent_id = ? AND ended_at IS NOT NULL`
	if err = db.QueryRow(query, agentID).Scan(&duration); err != nil {
		return
	}

	query = `SELECT SUM(first_response_wait_time) FROM custmchat.conversation WHERE agent_id = ? AND ended_at IS NOT NULL`
	if err = db.QueryRow(query, agentID).Scan(&waitRespDuration); err != nil {
		return
	}

	return
}

func GetAgentIDByConversationID(db XODB, convID string) (agentID string, err error) {
	query := `SELECT agent_id FROM custmchat.conversation where id=?`
	if err = db.QueryRow(query, convID).Scan(&agentID); err != nil {
		return
	}
	return
}
