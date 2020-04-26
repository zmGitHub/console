// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// InvitationMechsRule represents a row from 'custmchat.invitation_mechs_rules'.
type InvitationMechsRule struct {
	ID        string    `json:"id"`         // id
	EntID     string    `json:"ent_id"`     // ent_id
	Rule      string    `json:"rule"`       // rule
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the InvitationMechsRule exists in the database.
func (imr *InvitationMechsRule) Exists() bool {
	return imr._exists
}

// Deleted provides information if the InvitationMechsRule has been deleted from the database.
func (imr *InvitationMechsRule) Deleted() bool {
	return imr._deleted
}

// Insert inserts the InvitationMechsRule to the database.
func (imr *InvitationMechsRule) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if imr._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.invitation_mechs_rules (` +
		`id, ent_id, rule, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, imr.ID, imr.EntID, imr.Rule, imr.CreatedAt, imr.UpdatedAt)
	_, err = db.Exec(sqlstr, imr.ID, imr.EntID, imr.Rule, imr.CreatedAt, imr.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	imr._exists = true

	return nil
}

// Update updates the InvitationMechsRule in the database.
func (imr *InvitationMechsRule) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !imr._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if imr._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE custmchat.invitation_mechs_rules SET ` +
		`ent_id = ?, rule = ?, created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, imr.EntID, imr.Rule, imr.CreatedAt, imr.UpdatedAt, imr.ID)
	_, err = db.Exec(sqlstr, imr.EntID, imr.Rule, imr.CreatedAt, imr.UpdatedAt, imr.ID)
	return err
}

// Save saves the InvitationMechsRule to the database.
func (imr *InvitationMechsRule) Save(db XODB) error {
	if imr.Exists() {
		return imr.Update(db)
	}

	return imr.Insert(db)
}

// Delete deletes the InvitationMechsRule from the database.
func (imr *InvitationMechsRule) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !imr._exists {
		return nil
	}

	// if deleted, bail
	if imr._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.invitation_mechs_rules WHERE id = ?`

	// run query
	XOLog(sqlstr, imr.ID)
	_, err = db.Exec(sqlstr, imr.ID)
	if err != nil {
		return err
	}

	// set deleted
	imr._deleted = true

	return nil
}

// InvitationMechsRulesByEntID retrieves a row from 'custmchat.invitation_mechs_rules' as a InvitationMechsRule.
//
// Generated from index 'idx_enterprise_id'.
func InvitationMechsRulesByEntID(db XODB, entID string) ([]*InvitationMechsRule, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, rule, created_at, updated_at ` +
		`FROM custmchat.invitation_mechs_rules ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*InvitationMechsRule{}
	for q.Next() {
		imr := InvitationMechsRule{
			_exists: true,
		}

		// scan
		err = q.Scan(&imr.ID, &imr.EntID, &imr.Rule, &imr.CreatedAt, &imr.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &imr)
	}

	return res, nil
}

// InvitationMechsRuleByID retrieves a row from 'custmchat.invitation_mechs_rules' as a InvitationMechsRule.
//
// Generated from index 'invitation_mechs_rules_id_pkey'.
func InvitationMechsRuleByID(db XODB, id string) (*InvitationMechsRule, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, rule, created_at, updated_at ` +
		`FROM custmchat.invitation_mechs_rules ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	imr := InvitationMechsRule{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&imr.ID, &imr.EntID, &imr.Rule, &imr.CreatedAt, &imr.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &imr, nil
}
