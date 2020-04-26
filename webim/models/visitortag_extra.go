package models

func GetClientTagRanks(db XODB, entID string) (tags []*VisitorTag, err error) {
	// sql query
	const sqlstr = `SELECT id, rank ` +
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
	for q.Next() {
		vt := VisitorTag{}

		// scan
		err = q.Scan(&vt.ID, &vt.Rank)
		if err != nil {
			return nil, err
		}

		tags = append(tags, &vt)
	}
	if err = q.Err(); err != nil {
		return
	}

	return tags, nil
}

func LastTagRank(db XODB, entID string) (rank int, err error) {
	query := "SELECT rank FROM custmchat.visitor_tag WHERE ent_id=? ORDER BY rank DESC LIMIT 1"
	err = db.QueryRow(query, entID).Scan(&rank)
	return
}
