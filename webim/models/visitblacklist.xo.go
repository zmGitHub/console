// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// VisitBlacklist represents a row from 'custmchat.visit_blacklist'.
type VisitBlacklist struct {
	ID        string    `json:"id"`         // id
	EntID     string    `json:"ent_id"`     // ent_id
	TraceID   string    `json:"trace_id"`   // trace_id
	VisitID   string    `json:"visit_id"`   // visit_id
	AgentID   string    `json:"agent_id"`   // agent_id
	ConvID    string    `json:"conv_id"`    // conv_id
	CreatedAt time.Time `json:"created_at"` // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the VisitBlacklist exists in the database.
func (vb *VisitBlacklist) Exists() bool {
	return vb._exists
}

// Deleted provides information if the VisitBlacklist has been deleted from the database.
func (vb *VisitBlacklist) Deleted() bool {
	return vb._deleted
}

// Insert inserts the VisitBlacklist to the database.
func (vb *VisitBlacklist) Insert(db XODB) error {
	var err error

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.visit_blacklist (` +
		`id, ent_id, trace_id, visit_id, agent_id, conv_id, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, vb.ID, vb.EntID, vb.TraceID, vb.VisitID, vb.AgentID, vb.ConvID, vb.CreatedAt)
	_, err = db.Exec(sqlstr, vb.ID, vb.EntID, vb.TraceID, vb.VisitID, vb.AgentID, vb.ConvID, vb.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the VisitBlacklist in the database.
func (vb *VisitBlacklist) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE custmchat.visit_blacklist SET ` +
		`ent_id = ?, trace_id = ?, visit_id = ?, agent_id = ?, conv_id = ?, created_at = ?` +
		` WHERE id = ?`

	_, err = db.Exec(sqlstr, vb.EntID, vb.TraceID, vb.VisitID, vb.AgentID, vb.ConvID, vb.CreatedAt, vb.ID)
	return err
}

// Save saves the VisitBlacklist to the database.
func (vb *VisitBlacklist) Save(db XODB) error {
	if vb.Exists() {
		return vb.Update(db)
	}

	return vb.Insert(db)
}

// Delete deletes the VisitBlacklist from the database.
func (vb *VisitBlacklist) Delete(db XODB) error {
	var err error

	// sql query
	const sqlstr = `DELETE FROM custmchat.visit_blacklist WHERE id = ?`

	_, err = db.Exec(sqlstr, vb.ID)
	if err != nil {
		return err
	}

	return nil
}

// VisitBlacklistByEntIDTraceID retrieves a row from 'custmchat.visit_blacklist' as a VisitBlacklist.
//
// Generated from index 'ent_id'.
func VisitBlacklistByEntIDTraceID(db XODB, entID string, traceID string) (*VisitBlacklist, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, visit_id, agent_id, conv_id, created_at ` +
		`FROM custmchat.visit_blacklist ` +
		`WHERE ent_id = ? AND trace_id = ?`

	vb := VisitBlacklist{}
	err = db.QueryRow(sqlstr, entID, traceID).Scan(&vb.ID, &vb.EntID, &vb.TraceID, &vb.VisitID, &vb.AgentID, &vb.ConvID, &vb.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &vb, nil
}

// VisitBlacklistByID retrieves a row from 'custmchat.visit_blacklist' as a VisitBlacklist.
//
// Generated from index 'visit_blacklist_id_pkey'.
func VisitBlacklistByID(db XODB, id string) (*VisitBlacklist, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, visit_id, agent_id, conv_id, created_at ` +
		`FROM custmchat.visit_blacklist ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	vb := VisitBlacklist{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&vb.ID, &vb.EntID, &vb.TraceID, &vb.VisitID, &vb.AgentID, &vb.ConvID, &vb.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &vb, nil
}