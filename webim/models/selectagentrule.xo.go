// Package new_models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// SelectAgentRule represents a row from 'custmchat.select_agent_rule'.
type SelectAgentRule struct {
	ID          string         `json:"id"`           // id
	EntID       string         `json:"ent_id"`       // ent_id
	Rank        sql.NullInt64  `json:"rank"`         // rank
	Type        string         `json:"type"`         // type
	URLType     sql.NullString `json:"url_type"`     // url_type
	MatchType   sql.NullString `json:"match_type"`   // match_type
	MatchString sql.NullString `json:"match_string"` // match_string
	MatchRules  sql.NullString `json:"match_rules"`  // match_rules
	SourceRules sql.NullString `json:"source_rules"` // source_rules
	Targets     sql.NullString `json:"targets"`      // targets
	Inverted    sql.NullBool   `json:"inverted"`     // inverted

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the SelectAgentRule exists in the database.
func (sar *SelectAgentRule) Exists() bool {
	return sar._exists
}

// Deleted provides information if the SelectAgentRule has been deleted from the database.
func (sar *SelectAgentRule) Deleted() bool {
	return sar._deleted
}

// Insert inserts the SelectAgentRule to the database.
func (sar *SelectAgentRule) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if sar._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.select_agent_rule (` +
		`id, ent_id, rank, type, url_type, match_type, match_string, match_rules, source_rules, targets, inverted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, sar.ID, sar.EntID, sar.Rank, sar.Type, sar.URLType, sar.MatchType, sar.MatchString, sar.MatchRules, sar.SourceRules, sar.Targets, sar.Inverted)
	_, err = db.Exec(sqlstr, sar.ID, sar.EntID, sar.Rank, sar.Type, sar.URLType, sar.MatchType, sar.MatchString, sar.MatchRules, sar.SourceRules, sar.Targets, sar.Inverted)
	if err != nil {
		return err
	}

	// set existence
	sar._exists = true

	return nil
}

// Update updates the SelectAgentRule in the database.
func (sar *SelectAgentRule) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !sar._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if sar._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE custmchat.select_agent_rule SET ` +
		`ent_id = ?, rank = ?, type = ?, url_type = ?, match_type = ?, match_string = ?, match_rules = ?, source_rules = ?, targets = ?, inverted = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, sar.EntID, sar.Rank, sar.Type, sar.URLType, sar.MatchType, sar.MatchString, sar.MatchRules, sar.SourceRules, sar.Targets, sar.Inverted, sar.ID)
	_, err = db.Exec(sqlstr, sar.EntID, sar.Rank, sar.Type, sar.URLType, sar.MatchType, sar.MatchString, sar.MatchRules, sar.SourceRules, sar.Targets, sar.Inverted, sar.ID)
	return err
}

// Save saves the SelectAgentRule to the database.
func (sar *SelectAgentRule) Save(db XODB) error {
	if sar.Exists() {
		return sar.Update(db)
	}

	return sar.Insert(db)
}

// Delete deletes the SelectAgentRule from the database.
func (sar *SelectAgentRule) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !sar._exists {
		return nil
	}

	// if deleted, bail
	if sar._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.select_agent_rule WHERE id = ?`

	// run query
	XOLog(sqlstr, sar.ID)
	_, err = db.Exec(sqlstr, sar.ID)
	if err != nil {
		return err
	}

	// set deleted
	sar._deleted = true

	return nil
}

// SelectAgentRulesByEntID retrieves a row from 'custmchat.select_agent_rule' as a SelectAgentRule.
//
// Generated from index 'ent_id'.
func SelectAgentRulesByEntID(db XODB, entID string) ([]*SelectAgentRule, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, rank, type, url_type, match_type, match_string, match_rules, source_rules, targets, inverted ` +
		`FROM custmchat.select_agent_rule ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*SelectAgentRule{}
	for q.Next() {
		sar := SelectAgentRule{
			_exists: true,
		}

		// scan
		err = q.Scan(&sar.ID, &sar.EntID, &sar.Rank, &sar.Type, &sar.URLType, &sar.MatchType, &sar.MatchString, &sar.MatchRules, &sar.SourceRules, &sar.Targets, &sar.Inverted)
		if err != nil {
			return nil, err
		}

		res = append(res, &sar)
	}

	return res, nil
}

// SelectAgentRuleByID retrieves a row from 'custmchat.select_agent_rule' as a SelectAgentRule.
//
// Generated from index 'select_agent_rule_id_pkey'.
func SelectAgentRuleByID(db XODB, id string) (*SelectAgentRule, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, rank, type, url_type, match_type, match_string, match_rules, source_rules, targets, inverted ` +
		`FROM custmchat.select_agent_rule ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	sar := SelectAgentRule{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&sar.ID, &sar.EntID, &sar.Rank, &sar.Type, &sar.URLType, &sar.MatchType, &sar.MatchString, &sar.MatchRules, &sar.SourceRules, &sar.Targets, &sar.Inverted)
	if err != nil {
		return nil, err
	}

	return &sar, nil
}
