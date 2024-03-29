// Package new_models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"fmt"
	"time"
)

// VisitorTag represents a row from 'custmchat.visitor_tag'.
type VisitorTag struct {
	ID        string    `json:"id"`         // id
	EntID     string    `json:"ent_id"`     // ent_id
	Creator   string    `json:"creator"`    // creator
	Name      string    `json:"name"`       // name
	Color     string    `json:"color"`      // color
	UseCount  int       `json:"use_count"`  // use_count
	Rank      int       `json:"rank"`       // rank
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

func (vt *VisitorTag) Validate() error {
	if vt.Name == "" {
		return fmt.Errorf("tag name is empty")
	}

	if vt.Color == "" {
		return fmt.Errorf("color is empty")
	}

	return nil
}

// Exists determines if the VisitorTag exists in the database.
func (vt *VisitorTag) Exists() bool {
	return vt._exists
}

// Deleted provides information if the VisitorTag has been deleted from the database.
func (vt *VisitorTag) Deleted() bool {
	return vt._deleted
}

// Insert inserts the VisitorTag to the database.
func (vt *VisitorTag) Insert(db XODB) error {
	var err error

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.visitor_tag (` +
		`id, ent_id, creator, name, color, use_count, rank, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	_, err = db.Exec(sqlstr, vt.ID, vt.EntID, vt.Creator, vt.Name, vt.Color, vt.UseCount, vt.Rank, vt.CreatedAt, vt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the VisitorTag in the database.
func (vt *VisitorTag) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE custmchat.visitor_tag SET ` +
		`ent_id = ?, creator = ?, name = ?, color = ?, use_count = ?, rank=?, created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	_, err = db.Exec(sqlstr, vt.EntID, vt.Creator, vt.Name, vt.Color, vt.UseCount, vt.Rank, vt.CreatedAt, vt.UpdatedAt, vt.ID)
	return err
}

// Save saves the VisitorTag to the database.
func (vt *VisitorTag) Save(db XODB) error {
	if vt.Exists() {
		return vt.Update(db)
	}

	return vt.Insert(db)
}

// Delete deletes the VisitorTag from the database.
func (vt *VisitorTag) Delete(db XODB) error {
	var err error

	// sql query
	const sqlstr = `DELETE FROM custmchat.visitor_tag WHERE id = ?`

	// run query
	XOLog(sqlstr, vt.ID)
	_, err = db.Exec(sqlstr, vt.ID)
	if err != nil {
		return err
	}

	return nil
}

// VisitorTagsByEntID retrieves a row from 'custmchat.visitor_tag' as a VisitorTag.
//
// Generated from index 'idx_ent'.
func VisitorTagsByEntID(db XODB, entID string) ([]*VisitorTag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, creator, name, color, use_count, rank, created_at, updated_at ` +
		`FROM custmchat.visitor_tag ` +
		`WHERE ent_id = ? ` +
		`ORDER BY rank ASC, created_at DESC`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*VisitorTag
	for q.Next() {
		vt := VisitorTag{}

		// scan
		err = q.Scan(&vt.ID, &vt.EntID, &vt.Creator, &vt.Name, &vt.Color, &vt.UseCount, &vt.Rank, &vt.CreatedAt, &vt.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &vt)
	}

	return res, nil
}

// VisitorTagByID retrieves a row from 'custmchat.visitor_tag' as a VisitorTag.
//
// Generated from index 'visitor_tag_id_pkey'.
func VisitorTagByID(db XODB, id string) (*VisitorTag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, creator, name, color, use_count, rank, created_at, updated_at ` +
		`FROM custmchat.visitor_tag ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	vt := VisitorTag{}

	err = db.QueryRow(sqlstr, id).Scan(&vt.ID, &vt.EntID, &vt.Creator, &vt.Name, &vt.Color, &vt.UseCount, &vt.Rank, &vt.CreatedAt, &vt.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &vt, nil
}
