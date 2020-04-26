// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// AgentStatistic represents a row from 'custmchat.agent_statistic'.
type AgentStatistic struct {
	EntID             string    `json:"ent_id"`              // ent_id
	AgentID           string    `json:"agent_id"`            // agent_id
	ConversationCount uint      `json:"conversation_count"`  // conversation_count
	GoodCount         uint      `json:"good_count"`          // good_count
	MediumCount       uint      `json:"medium_count"`        // medium_count
	BadCount          uint      `json:"bad_count"`           // bad_count
	MessageCount      uint      `json:"message_count"`       // message_count
	Duration          int       `json:"duration"`            // duration
	FirstRespDuration int       `json:"first_resp_duration"` // first_resp_duration
	CreatedAt         time.Time `json:"created_at"`          // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the AgentStatistic exists in the database.
func (as *AgentStatistic) Exists() bool {
	return as._exists
}

// Deleted provides information if the AgentStatistic has been deleted from the database.
func (as *AgentStatistic) Deleted() bool {
	return as._deleted
}

// Insert inserts the AgentStatistic to the database.
func (as *AgentStatistic) Insert(db XODB, updates []string) (err error) {
	sqlstr := `INSERT INTO custmchat.agent_statistic (` +
		`ent_id, agent_id, conversation_count, good_count, medium_count, bad_count, message_count, duration, first_resp_duration, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE %s`

	var placeHolders []string
	var args = []interface{}{
		as.EntID, as.AgentID, as.ConversationCount, as.GoodCount, as.MediumCount, as.BadCount, as.MessageCount, as.Duration, as.FirstRespDuration, as.CreatedAt,
	}
	for _, name := range updates {
		placeHolders = append(placeHolders, fmt.Sprintf("%s = %s + VALUES(%s)", name, name, name))
	}

	sqlstr = fmt.Sprintf(sqlstr, strings.Join(placeHolders, ","))
	_, err = db.Exec(sqlstr, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the AgentStatistic in the database.
func (as *AgentStatistic) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !as._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if as._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query with composite primary key
	const sqlstr = `UPDATE custmchat.agent_statistic SET ` +
		`ent_id = ?, conversation_count = ?, good_count = ?, medium_count = ?, bad_count = ?, message_count = ?, duration = ?, first_resp_duration = ?` +
		` WHERE agent_id = ? AND created_at = ?`

	// run query
	XOLog(sqlstr, as.EntID, as.ConversationCount, as.GoodCount, as.MediumCount, as.BadCount, as.MessageCount, as.Duration, as.FirstRespDuration, as.AgentID, as.CreatedAt)
	_, err = db.Exec(sqlstr, as.EntID, as.ConversationCount, as.GoodCount, as.MediumCount, as.BadCount, as.MessageCount, as.Duration, as.FirstRespDuration, as.AgentID, as.CreatedAt)
	return err
}

// Delete deletes the AgentStatistic from the database.
func (as *AgentStatistic) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !as._exists {
		return nil
	}

	// if deleted, bail
	if as._deleted {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `DELETE FROM custmchat.agent_statistic WHERE agent_id = ? AND created_at = ?`

	// run query
	XOLog(sqlstr, as.AgentID, as.CreatedAt)
	_, err = db.Exec(sqlstr, as.AgentID, as.CreatedAt)
	if err != nil {
		return err
	}

	// set deleted
	as._deleted = true

	return nil
}

// AgentStatisticByAgentIDCreatedAt retrieves a row from 'custmchat.agent_statistic' as a AgentStatistic.
//
// Generated from index 'agent_id'.
func AgentStatisticByAgentIDCreatedAt(db XODB, agentID string, createdAt time.Time) (*AgentStatistic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, agent_id, conversation_count, good_count, medium_count, bad_count, message_count, duration, first_resp_duration, created_at ` +
		`FROM custmchat.agent_statistic ` +
		`WHERE agent_id = ? AND created_at = ?`

	// run query
	XOLog(sqlstr, agentID, createdAt)
	as := AgentStatistic{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, agentID, createdAt).Scan(&as.EntID, &as.AgentID, &as.ConversationCount, &as.GoodCount, &as.MediumCount, &as.BadCount, &as.MessageCount, &as.Duration, &as.FirstRespDuration, &as.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &as, nil
}

// AgentStatisticByCreatedAt retrieves a row from 'custmchat.agent_statistic' as a AgentStatistic.
//
// Generated from index 'agent_statistic_created_at_pkey'.
func AgentStatisticByCreatedAt(db XODB, createdAt time.Time) (*AgentStatistic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, agent_id, conversation_count, good_count, medium_count, bad_count, message_count, duration, first_resp_duration, created_at ` +
		`FROM custmchat.agent_statistic ` +
		`WHERE created_at = ?`

	// run query
	XOLog(sqlstr, createdAt)
	as := AgentStatistic{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, createdAt).Scan(&as.EntID, &as.AgentID, &as.ConversationCount, &as.GoodCount, &as.MediumCount, &as.BadCount, &as.MessageCount, &as.Duration, &as.FirstRespDuration, &as.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &as, nil
}

// AgentStatisticsByEntID retrieves a row from 'custmchat.agent_statistic' as a AgentStatistic.
//
// Generated from index 'ent_id'.
func AgentStatisticsByEntID(db XODB, entID string, start, end time.Time) ([]*AgentStatistic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, agent_id, conversation_count, good_count, medium_count, bad_count, message_count, duration, first_resp_duration, created_at ` +
		`FROM custmchat.agent_statistic ` +
		`WHERE ent_id = ? AND created_at >= ? AND created_at <= ?`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AgentStatistic{}
	for q.Next() {
		as := AgentStatistic{}

		// scan
		err = q.Scan(&as.EntID, &as.AgentID, &as.ConversationCount, &as.GoodCount, &as.MediumCount, &as.BadCount, &as.MessageCount, &as.Duration, &as.FirstRespDuration, &as.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &as)
	}

	return res, nil
}