// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

// Role represents a row from 'custmchat.role'.
type Role struct {
	ID    string `json:"id"`     // id
	EntID string `json:"ent_id"` // ent_id
	Name  string `json:"name"`   // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Role exists in the database.
func (r *Role) Exists() bool {
	return r._exists
}

// Deleted provides information if the Role has been deleted from the database.
func (r *Role) Deleted() bool {
	return r._deleted
}

// Insert inserts the Role to the database.
func (r *Role) Insert(db XODB) error {
	var err error

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO custmchat.role (` +
		`id, ent_id, name` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, r.ID, r.EntID, r.Name)
	_, err = db.Exec(sqlstr, r.ID, r.EntID, r.Name)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the Role in the database.
func (r *Role) Update(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE custmchat.role SET ` +
		`ent_id = ?, name = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, r.EntID, r.Name, r.ID)
	_, err = db.Exec(sqlstr, r.EntID, r.Name, r.ID)
	return err
}

// Save saves the Role to the database.
func (r *Role) Save(db XODB) error {
	if r.Exists() {
		return r.Update(db)
	}

	return r.Insert(db)
}

// Delete deletes the Role from the database.
func (r *Role) Delete(db XODB) error {
	var err error

	// sql query
	const sqlstr = `DELETE FROM custmchat.role WHERE id = ?`

	// run query
	XOLog(sqlstr, r.ID)
	_, err = db.Exec(sqlstr, r.ID)
	if err != nil {
		return err
	}

	return nil
}

// RolesByEntID retrieves a row from 'custmchat.role' as a Role.
//
// Generated from index 'idx_ent'.
func RolesByEntID(db XODB, entID string) ([]*Role, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, name ` +
		`FROM custmchat.role ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Role{}
	for q.Next() {
		r := Role{
			_exists: true,
		}

		// scan
		err = q.Scan(&r.ID, &r.EntID, &r.Name)
		if err != nil {
			return nil, err
		}

		res = append(res, &r)
	}

	return res, nil
}

// RoleByID retrieves a row from 'custmchat.role' as a Role.
//
// Generated from index 'role_id_pkey'.
func RoleByID(db XODB, id string) (*Role, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, name ` +
		`FROM custmchat.role ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	r := Role{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&r.ID, &r.EntID, &r.Name)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
