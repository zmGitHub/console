// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// LeaveMessageConfig represents a row from 'custmchat.leave_message_config'.
type LeaveMessageConfig struct {
	EntID              string `json:"ent_id"`               // ent_id
	Introduction       string `json:"introduction"`         // introduction
	ShowVisitorName    bool   `json:"show_visitor_name"`    // show_visitor_name
	ShowTelephone      bool   `json:"show_telephone"`       // show_telephone
	ShowEmail          bool   `json:"show_email"`           // show_email
	ShowWechat         bool   `json:"show_wechat"`          // show_wechat
	ShowQq             bool   `json:"show_qq"`              // show_qq
	AutoCreateCategory bool   `json:"auto_create_category"` // auto_create_category
	FillContact        string `json:"fill_contact"`         // fill_contact
	UseDefaultContent  bool   `json:"use_default_content"`  // use_default_content
	DefaultContent     string `json:"default_content"`      // default_content

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the LeaveMessageConfig exists in the database.
func (lmc *LeaveMessageConfig) Exists() bool {
	return lmc._exists
}

// Deleted provides information if the LeaveMessageConfig has been deleted from the database.
func (lmc *LeaveMessageConfig) Deleted() bool {
	return lmc._deleted
}

// Insert inserts the LeaveMessageConfig to the database.
func (lmc *LeaveMessageConfig) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if lmc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.leave_message_config (` +
		`ent_id, introduction, show_visitor_name, show_telephone, show_email, show_wechat, show_qq, auto_create_category, fill_contact, use_default_content, default_content` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, lmc.EntID, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.ShowWechat, lmc.ShowQq, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent)
	_, err = db.Exec(sqlstr, lmc.EntID, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.ShowWechat, lmc.ShowQq, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent)
	if err != nil {
		return err
	}

	// set existence
	lmc._exists = true

	return nil
}

// Update updates the LeaveMessageConfig in the database.
func (lmc *LeaveMessageConfig) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !lmc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if lmc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE custmchat.leave_message_config SET ` +
		`introduction = ?, show_visitor_name = ?, show_telephone = ?, show_email = ?, show_wechat = ?, show_qq = ?, auto_create_category = ?, fill_contact = ?, use_default_content = ?, default_content = ?` +
		` WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.ShowWechat, lmc.ShowQq, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent, lmc.EntID)
	_, err = db.Exec(sqlstr, lmc.Introduction, lmc.ShowVisitorName, lmc.ShowTelephone, lmc.ShowEmail, lmc.ShowWechat, lmc.ShowQq, lmc.AutoCreateCategory, lmc.FillContact, lmc.UseDefaultContent, lmc.DefaultContent, lmc.EntID)
	return err
}

// Save saves the LeaveMessageConfig to the database.
func (lmc *LeaveMessageConfig) Save(db XODB) error {
	if lmc.Exists() {
		return lmc.Update(db)
	}

	return lmc.Insert(db)
}

// Delete deletes the LeaveMessageConfig from the database.
func (lmc *LeaveMessageConfig) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !lmc._exists {
		return nil
	}

	// if deleted, bail
	if lmc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.leave_message_config WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, lmc.EntID)
	_, err = db.Exec(sqlstr, lmc.EntID)
	if err != nil {
		return err
	}

	// set deleted
	lmc._deleted = true

	return nil
}

// LeaveMessageConfigByEntID retrieves a row from 'custmchat.leave_message_config' as a LeaveMessageConfig.
//
// Generated from index 'ent_id'.
func LeaveMessageConfigByEntID(db XODB, entID string) (*LeaveMessageConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, introduction, show_visitor_name, show_telephone, show_email, show_wechat, show_qq, auto_create_category, fill_contact, use_default_content, default_content ` +
		`FROM custmchat.leave_message_config ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	lmc := LeaveMessageConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, entID).Scan(&lmc.EntID, &lmc.Introduction, &lmc.ShowVisitorName, &lmc.ShowTelephone, &lmc.ShowEmail, &lmc.ShowWechat, &lmc.ShowQq, &lmc.AutoCreateCategory, &lmc.FillContact, &lmc.UseDefaultContent, &lmc.DefaultContent)
	if err != nil {
		return nil, err
	}

	return &lmc, nil
}