// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

var (
	EndingMessageWebPlatform = "web"
)

// EndingMessage represents a row from 'custmchat.ending_message'.
type EndingMessage struct {
	EntID    string `json:"ent_id"`   // ent_id
	Platform string `json:"platform"` // platform
	Agent    string `json:"agent"`    // agent
	System   string `json:"system"`   // system
	Status   bool   `json:"status"`   // status
	Prompt   bool   `json:"prompt"`   // prompt

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the EndingMessage exists in the database.
func (em *EndingMessage) Exists() bool {
	return em._exists
}

// Deleted provides information if the EndingMessage has been deleted from the database.
func (em *EndingMessage) Deleted() bool {
	return em._deleted
}

// Insert inserts the EndingMessage to the database.
func (em *EndingMessage) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if em._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.ending_message (` +
		`ent_id, platform, agent, system, status, prompt` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, em.EntID, em.Platform, em.Agent, em.System, em.Status, em.Prompt)
	_, err = db.Exec(sqlstr, em.EntID, em.Platform, em.Agent, em.System, em.Status, em.Prompt)
	if err != nil {
		return err
	}

	// set existence
	em._exists = true

	return nil
}

// Update updates the EndingMessage in the database.
func (em *EndingMessage) Update(db XODB) error {
	var err error

	// sql query with composite primary key
	const sqlstr = `UPDATE custmchat.ending_message SET ` +
		`agent = ?, system = ?, status = ?, prompt = ?` +
		` WHERE ent_id = ? AND platform = ?`

	// run query
	XOLog(sqlstr, em.Agent, em.System, em.Status, em.Prompt, em.EntID, em.Platform)
	_, err = db.Exec(sqlstr, em.Agent, em.System, em.Status, em.Prompt, em.EntID, em.Platform)
	return err
}

// Save saves the EndingMessage to the database.
func (em *EndingMessage) Save(db XODB) error {
	if em.Exists() {
		return em.Update(db)
	}

	return em.Insert(db)
}

// Delete deletes the EndingMessage from the database.
func (em *EndingMessage) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !em._exists {
		return nil
	}

	// if deleted, bail
	if em._deleted {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `DELETE FROM custmchat.ending_message WHERE ent_id = ? AND platform = ?`

	// run query
	XOLog(sqlstr, em.EntID, em.Platform)
	_, err = db.Exec(sqlstr, em.EntID, em.Platform)
	if err != nil {
		return err
	}

	// set deleted
	em._deleted = true

	return nil
}

// EndingMessageByPlatform retrieves a row from 'custmchat.ending_message' as a EndingMessage.
//
// Generated from index 'ending_message_platform_pkey'.
func EndingMessageByPlatform(db XODB, platform string) (*EndingMessage, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, platform, agent, system, status, prompt ` +
		`FROM custmchat.ending_message ` +
		`WHERE platform = ?`

	// run query
	XOLog(sqlstr, platform)
	em := EndingMessage{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, platform).Scan(&em.EntID, &em.Platform, &em.Agent, &em.System, &em.Status, &em.Prompt)
	if err != nil {
		return nil, err
	}

	return &em, nil
}

// EndingMessageByEntIDPlatform retrieves a row from 'custmchat.ending_message' as a EndingMessage.
//
// Generated from index 'ent_id'.
func EndingMessageByEntIDPlatform(db XODB, entID string, platform string) (*EndingMessage, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, platform, agent, system, status, prompt ` +
		`FROM custmchat.ending_message ` +
		`WHERE ent_id = ? AND platform = ?`

	// run query
	XOLog(sqlstr, entID, platform)
	em := EndingMessage{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, entID, platform).Scan(&em.EntID, &em.Platform, &em.Agent, &em.System, &em.Status, &em.Prompt)
	if err != nil {
		return nil, err
	}

	return &em, nil
}