// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"time"
)

// EntAllConfig represents a row from 'custmchat.ent_all_configs'.
type EntAllConfig struct {
	ID            string         `json:"id"`             // id
	EntID         string         `json:"ent_id"`         // ent_id
	ConfigContent sql.NullString `json:"config_content"` // config_content
	CreateAt      time.Time      `json:"create_at"`      // create_at
	UpdateAt      time.Time      `json:"update_at"`      // update_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the EntAllConfig exists in the database.
func (eac *EntAllConfig) Exists() bool {
	return eac._exists
}

// Deleted provides information if the EntAllConfig has been deleted from the database.
func (eac *EntAllConfig) Deleted() bool {
	return eac._deleted
}

// Insert inserts the EntAllConfig to the database.
func (eac *EntAllConfig) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if eac._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.ent_all_configs (` +
		`id, ent_id, config_content, create_at, update_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, eac.ID, eac.EntID, eac.ConfigContent, eac.CreateAt, eac.UpdateAt)
	_, err = db.Exec(sqlstr, eac.ID, eac.EntID, eac.ConfigContent, eac.CreateAt, eac.UpdateAt)
	if err != nil {
		return err
	}

	// set existence
	eac._exists = true

	return nil
}

// Update updates the EntAllConfig in the database.
func (eac *EntAllConfig) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !eac._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if eac._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE custmchat.ent_all_configs SET ` +
		`ent_id = ?, config_content = ?, create_at = ?, update_at = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, eac.EntID, eac.ConfigContent, eac.CreateAt, eac.UpdateAt, eac.ID)
	_, err = db.Exec(sqlstr, eac.EntID, eac.ConfigContent, eac.CreateAt, eac.UpdateAt, eac.ID)
	return err
}

// Save saves the EntAllConfig to the database.
func (eac *EntAllConfig) Save(db XODB) error {
	if eac.Exists() {
		return eac.Update(db)
	}

	return eac.Insert(db)
}

// Delete deletes the EntAllConfig from the database.
func (eac *EntAllConfig) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !eac._exists {
		return nil
	}

	// if deleted, bail
	if eac._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.ent_all_configs WHERE id = ?`

	// run query
	XOLog(sqlstr, eac.ID)
	_, err = db.Exec(sqlstr, eac.ID)
	if err != nil {
		return err
	}

	// set deleted
	eac._deleted = true

	return nil
}

// EntAllConfigByID retrieves a row from 'custmchat.ent_all_configs' as a EntAllConfig.
//
// Generated from index 'ent_all_configs_id_pkey'.
func EntAllConfigByID(db XODB, id string) (*EntAllConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, config_content, create_at, update_at ` +
		`FROM custmchat.ent_all_configs ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	eac := EntAllConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&eac.ID, &eac.EntID, &eac.ConfigContent, &eac.CreateAt, &eac.UpdateAt)
	if err != nil {
		return nil, err
	}

	return &eac, nil
}

// EntAllConfigByEntID retrieves a row from 'custmchat.ent_all_configs' as a EntAllConfig.
//
// Generated from index 'ent_id'.
func EntAllConfigByEntID(db XODB, entID string) (*EntAllConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, config_content, create_at, update_at ` +
		`FROM custmchat.ent_all_configs ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	eac := EntAllConfig{}

	err = db.QueryRow(sqlstr, entID).Scan(&eac.ID, &eac.EntID, &eac.ConfigContent, &eac.CreateAt, &eac.UpdateAt)
	if err != nil {
		return nil, err
	}

	return &eac, nil
}