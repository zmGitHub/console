// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
)

// PersonalConfig represents a row from 'custmchat.personal_config'.
type PersonalConfig struct {
	AgentID       string         `json:"agent_id"`       // agent_id
	ConfigContent sql.NullString `json:"config_content"` // config_content

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the PersonalConfig exists in the database.
func (pc *PersonalConfig) Exists() bool {
	return pc._exists
}

// Deleted provides information if the PersonalConfig has been deleted from the database.
func (pc *PersonalConfig) Deleted() bool {
	return pc._deleted
}

// Insert inserts the PersonalConfig to the database.
func (pc *PersonalConfig) Insert(db XODB) error {
	var err error

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.personal_config (` +
		`agent_id, config_content` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, pc.AgentID, pc.ConfigContent)
	_, err = db.Exec(sqlstr, pc.AgentID, pc.ConfigContent)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the PersonalConfig in the database.
func (pc *PersonalConfig) Update(db XODB) error {
	var err error
	// sql query
	const sqlstr = `UPDATE custmchat.personal_config SET ` +
		`config_content = ?` +
		` WHERE agent_id = ?`

	// run query
	XOLog(sqlstr, pc.ConfigContent, pc.AgentID)
	_, err = db.Exec(sqlstr, pc.ConfigContent, pc.AgentID)
	return err
}

// Save saves the PersonalConfig to the database.
func (pc *PersonalConfig) Save(db XODB) error {
	if pc.Exists() {
		return pc.Update(db)
	}

	return pc.Insert(db)
}

// Delete deletes the PersonalConfig from the database.
func (pc *PersonalConfig) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !pc._exists {
		return nil
	}

	// if deleted, bail
	if pc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.personal_config WHERE agent_id = ?`

	// run query
	XOLog(sqlstr, pc.AgentID)
	_, err = db.Exec(sqlstr, pc.AgentID)
	if err != nil {
		return err
	}

	// set deleted
	pc._deleted = true

	return nil
}

// PersonalConfigByAgentID retrieves a row from 'custmchat.personal_config' as a PersonalConfig.
//
// Generated from index 'agent_id'.
func PersonalConfigByAgentID(db XODB, agentID string) (*PersonalConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`agent_id, config_content ` +
		`FROM custmchat.personal_config ` +
		`WHERE agent_id = ?`

	// run query
	XOLog(sqlstr, agentID)
	pc := PersonalConfig{}

	err = db.QueryRow(sqlstr, agentID).Scan(&pc.AgentID, &pc.ConfigContent)
	if err != nil {
		return nil, err
	}

	return &pc, nil
}
