// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// EntPlan represents a row from 'custmchat.ent_plan'.
type EntPlan struct {
	ID              string    `json:"id"`                // id
	EntID           string    `json:"ent_id"`            // ent_id
	PlanType        int8      `json:"plan_type"`         // plan_type
	TrialStatus     int       `json:"trial_status"`      // trial_status
	AgentServeLimit int       `json:"agent_serve_limit"` // agent_serve_limit
	LoginAgentLimit int       `json:"login_agent_limit"` // login_agent_limit
	AgentNum        int       `json:"agent_num"`         // agent_num
	PayAmount       int       `json:"pay_amount"`        // pay_amount
	ExpirationTime  time.Time `json:"expiration_time"`   // expiration_time
	CreateAt        time.Time `json:"create_at"`         // create_at
	UpdateAt        time.Time `json:"update_at"`         // update_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the EntPlan exists in the database.
func (ep *EntPlan) Exists() bool {
	return ep._exists
}

// Deleted provides information if the EntPlan has been deleted from the database.
func (ep *EntPlan) Deleted() bool {
	return ep._deleted
}

// Insert inserts the EntPlan to the database.
func (ep *EntPlan) Insert(db XODB) error {
	var err error

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.ent_plan (` +
		`id, ent_id, plan_type, trial_status, agent_serve_limit, login_agent_limit, agent_num, pay_amount, expiration_time, create_at, update_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, ep.ID, ep.EntID, ep.PlanType, ep.TrialStatus, ep.AgentServeLimit, ep.LoginAgentLimit, ep.AgentNum, ep.PayAmount, ep.ExpirationTime, ep.CreateAt, ep.UpdateAt)
	_, err = db.Exec(sqlstr, ep.ID, ep.EntID, ep.PlanType, ep.TrialStatus, ep.AgentServeLimit, ep.LoginAgentLimit, ep.AgentNum, ep.PayAmount, ep.ExpirationTime, ep.CreateAt, ep.UpdateAt)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the EntPlan in the database.
func (ep *EntPlan) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE custmchat.ent_plan SET ` +
		`ent_id = ?, plan_type = ?, trial_status = ?, agent_serve_limit = ?, login_agent_limit = ?, agent_num = ?, pay_amount = ?, expiration_time = ?, create_at = ?, update_at = ?` +
		` WHERE id = ?`

	_, err = db.Exec(sqlstr, ep.EntID, ep.PlanType, ep.TrialStatus, ep.AgentServeLimit, ep.LoginAgentLimit, ep.AgentNum, ep.PayAmount, ep.ExpirationTime, ep.CreateAt, ep.UpdateAt, ep.ID)
	return err
}

// Save saves the EntPlan to the database.
func (ep *EntPlan) Save(db XODB) error {
	if ep.Exists() {
		return ep.Update(db)
	}

	return ep.Insert(db)
}

// Delete deletes the EntPlan from the database.
func (ep *EntPlan) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ep._exists {
		return nil
	}

	// if deleted, bail
	if ep._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM custmchat.ent_plan WHERE id = ?`

	// run query
	XOLog(sqlstr, ep.ID)
	_, err = db.Exec(sqlstr, ep.ID)
	if err != nil {
		return err
	}

	// set deleted
	ep._deleted = true

	return nil
}

// EntPlanByEntID retrieves a row from 'custmchat.ent_plan' as a EntPlan.
//
// Generated from index 'ent_id'.
func EntPlanByEntID(db XODB, entID string) (*EntPlan, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, plan_type, trial_status, agent_serve_limit, login_agent_limit, agent_num, pay_amount, expiration_time, create_at, update_at ` +
		`FROM custmchat.ent_plan ` +
		`WHERE ent_id = ?`

	ep := EntPlan{}

	err = db.QueryRow(sqlstr, entID).Scan(&ep.ID, &ep.EntID, &ep.PlanType, &ep.TrialStatus, &ep.AgentServeLimit, &ep.LoginAgentLimit, &ep.AgentNum, &ep.PayAmount, &ep.ExpirationTime, &ep.CreateAt, &ep.UpdateAt)
	if err != nil {
		return nil, err
	}

	return &ep, nil
}

// EntPlanByID retrieves a row from 'custmchat.ent_plan' as a EntPlan.
//
// Generated from index 'ent_plan_id_pkey'.
func EntPlanByID(db XODB, id string) (*EntPlan, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, plan_type, trial_status, agent_serve_limit, login_agent_limit, agent_num, pay_amount, expiration_time, create_at, update_at ` +
		`FROM custmchat.ent_plan ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	ep := EntPlan{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&ep.ID, &ep.EntID, &ep.PlanType, &ep.TrialStatus, &ep.AgentServeLimit, &ep.LoginAgentLimit, &ep.AgentNum, &ep.PayAmount, &ep.ExpirationTime, &ep.CreateAt, &ep.UpdateAt)
	if err != nil {
		return nil, err
	}

	return &ep, nil
}