package handler

import "database/sql"

func rollBackOrCommit(tx *sql.Tx, err error) error {
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
