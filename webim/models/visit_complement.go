package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/forfd/custm-chat/webim/conf"
)

var (
	visitFields = `id, ent_id, trace_id, visit_page_cnt, residence_time_sec, browser_family, browser_version_string, browser_version, os_category, os_family, os_version_string, os_version, platform, ua_string, ip, country, province, city, isp, first_page_source, first_page_source_keyword, first_page_source_domain, first_page_source_url, first_page_title, first_page_domain, first_page_url, latest_title, latest_url, created_at, updated_at`

	maxOnlineDuration     = 60 * time.Second
	defaultOnlineDuration = 25 * time.Second
)

type VisitorTrack struct {
	TrackID string
	EntID   string
}

func IsVisitNotExists(db XODB, entID, trackID string) bool {
	query := `SELECT id FROM custmchat.visit WHERE ent_id=? and trace_id=?`
	var id string
	if err := db.QueryRow(query, entID, trackID).Scan(&id); err != nil {
		return err == sql.ErrNoRows
	}

	return false
}

func IncrVisitResidenceTimeSec(db XODB, traceID string, residenceTimeSec int64) error {
	sqlStr := `UPDATE custmchat.visit SET residence_time_sec = residence_time_sec + ? WHERE trace_id = ?`
	if _, err := db.Exec(sqlStr, residenceTimeSec, traceID); err != nil {
		return err
	}

	return nil
}

func IncrVisitorResidenceTimeSec(db XODB, traceID string, residenceTimeSec int64) error {
	sqlStr := `UPDATE custmchat.visitor SET residence_time_sec = residence_time_sec + ? WHERE trace_id = ?`
	if _, err := db.Exec(sqlStr, residenceTimeSec, traceID); err != nil {
		return err
	}

	return nil
}

func IsClientOnline(db XODB, entID, traceID string) (bool, error) {
	query := `SELECT trace_id FROM online_visitors WHERE ent_id = ? AND trace_id = ? AND updated_at >= ?`
	now := time.Now().UTC()
	t := now.Add(-1 * conf.DefaultPingInterval)
	interval := conf.IMConf.CentrifugoConf.PingInterval.Duration
	if interval >= conf.MinPingInterval && interval <= conf.MaxPingInterval {
		t = now.Add(-1 * interval)
	}

	var v string
	if err := db.QueryRow(query, entID, traceID, t).Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return v == traceID, nil
}

func IncreaseVisitPageCount(db XODB, visitID string, pageCount int64) error {
	sqlStr := `UPDATE custmchat.visit SET visit_page_cnt = visit_page_cnt + ? WHERE id = ?`
	if _, err := db.Exec(sqlStr, pageCount, visitID); err != nil {
		return err
	}

	return nil
}

func IncreaseVisitorPageCount(db XODB, visitorID string, pageCount int64, visitedAt time.Time, visitCnt int64, lastVisitID string) error {
	sqlStr := `UPDATE custmchat.visitor SET visit_page_cnt = visit_page_cnt + ?, visit_cnt = visit_cnt + ?, visited_at = ?, %s WHERE id = ?`

	var args = []interface{}{pageCount, visitCnt, visitedAt}
	var s string
	if lastVisitID != "" {
		args = append(args, lastVisitID)
		s = "last_visit_id = ? "
	}

	sqlStr = fmt.Sprintf(sqlStr, s)
	args = append(args, visitorID)
	if _, err := db.Exec(sqlStr, args...); err != nil {
		return err
	}

	return nil
}

func UpdateOnlineVisitor(db XODB, traceID, entID string) error {
	q := `INSERT INTO online_visitors (ent_id, trace_id, created_at, updated_at) VALUES (?,?,?,?) ` +
		`ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)`
	now := time.Now().UTC()
	_, err := db.Exec(q, entID, traceID, now, now)
	return err
}

func BulkUpdateOnlineVisitor(db XODB, tracks []*VisitorTrack) error {
	if len(tracks) == 0 {
		return nil
	}

	now := time.Now().UTC()
	insert := `INSERT INTO online_visitors (ent_id, trace_id, created_at, updated_at) VALUES %s ` +
		`ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)`

	var args []interface{}
	var placeHolders []string
	for _, track := range tracks {
		placeHolders = append(placeHolders, "(?,?,?,?)")
		args = append(args, track.EntID)
		args = append(args, track.TrackID)
		args = append(args, now)
		args = append(args, now)
	}

	insert = fmt.Sprintf(insert, strings.Join(placeHolders, ","))
	_, err := db.Exec(insert, args...)
	return err
}

func DeleteOnlineVisitor(db XODB, traceID, entID string) error {
	d := `DELETE FROM online_visitors WHERE ent_id=? AND trace_id = ?`
	_, err := db.Exec(d, entID, traceID)
	return err
}

func OnlineVisitorsByEntID(db XODB, entID string) (visitors []*OnlineVisitor, err error) {
	const sqlstr = `SELECT ` +
		`ent_id, trace_id, created_at, updated_at ` +
		`FROM custmchat.online_visitors ` +
		`WHERE ent_id = ? AND updated_at >= ? AND updated_at <= ?`

	now := time.Now().UTC()
	d := conf.IMConf.CentrifugoConf.PingInterval.Duration
	if d <= defaultOnlineDuration || d >= maxOnlineDuration {
		d = defaultOnlineDuration
	}
	start := now.Add(-1 * d)
	end := now.Add(d)

	q, err := db.Query(sqlstr, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	for q.Next() {
		ov := OnlineVisitor{}

		// scan
		err = q.Scan(&ov.EntID, &ov.TraceID, &ov.CreatedAt, &ov.UpdatedAt)
		if err != nil {
			return nil, err
		}

		visitors = append(visitors, &ov)
	}

	return
}

func (v *Visitor) UpdateVisitorCount(db XODB) error {
	updateSQL := `UPDATE custmchat.visitor SET ` +
		`visit_cnt = visit_cnt + 1, visit_page_cnt = visit_page_cnt + 1, ` +
		`last_visit_id = ?, visited_at = ? WHERE ent_id = ? AND trace_id = ?`

	_, err := db.Exec(updateSQL, v.LastVisitID, v.VisitedAt, v.EntID, v.TraceID)
	return err
}

func (v *Visitor) UpdateVisitorResidenceTimeSec(db XODB, residenceTimeSec int64) (err error) {
	updateSQL := `UPDATE custmchat.visitor SET ` +
		`residence_time_sec = residence_time_sec + ? WHERE ent_id = ? AND trace_id = ?`

	_, err = db.Exec(updateSQL, residenceTimeSec, v.EntID, v.TraceID)
	return
}

func (v *Visitor) UpdateVisitorName(db XODB) error {
	updateSQL := `UPDATE custmchat.visitor SET ` +
		`name = ? WHERE id = ?`

	_, err := db.Exec(updateSQL, v.Name, v.ID)
	return err
}

func VisitsByTraceIDs(db XODB, traceIDs []string) ([]*Visit, error) {
	if len(traceIDs) == 0 {
		return nil, nil
	}

	sqlstr := `SELECT %s FROM custmchat.visit WHERE trace_id IN (%s)`

	var args []interface{}
	var placeHolders []string
	for _, id := range traceIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	sqlstr = fmt.Sprintf(sqlstr, visitFields, strings.Join(placeHolders, ","))
	return visitsByQuery(db, sqlstr, args...)
}

func TraceVisitsByConds(db XODB, traceIDs []string, province string, limit int) ([]*Visit, error) {
	if len(traceIDs) == 0 {
		return []*Visit{}, nil
	}

	query := `SELECT %s FROM custmchat.visit WHERE trace_id IN (%s) %s %s `

	var args []interface{}
	var placeHolders []string
	for _, id := range traceIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	var provinceQuery = ""
	if province != "" {
		provinceQuery = ` AND province = ? `
		args = append(args, province)
	}

	var limitQuery = ""
	if limit > 0 {
		limitQuery = ` LIMIT 0, ?`
		args = append(args, limit)
	}

	query = fmt.Sprintf(query, visitFields, strings.Join(placeHolders, ","), provinceQuery, limitQuery)
	return visitsByQuery(db, query, args...)
}

func visitsByQuery(db XODB, query string, args ...interface{}) (visits []*Visit, err error) {
	q, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	for q.Next() {
		v := Visit{}

		// scan
		err = q.Scan(&v.ID, &v.EntID, &v.TraceID, &v.VisitPageCnt, &v.ResidenceTimeSec, &v.BrowserFamily, &v.BrowserVersionString, &v.BrowserVersion, &v.OsCategory, &v.OsFamily, &v.OsVersionString, &v.OsVersion, &v.Platform, &v.UaString, &v.IP, &v.Country, &v.Province, &v.City, &v.Isp, &v.FirstPageSource, &v.FirstPageSourceKeyword, &v.FirstPageSourceDomain, &v.FirstPageSourceURL, &v.FirstPageTitle, &v.FirstPageDomain, &v.FirstPageURL, &v.LatestTitle, &v.LatestURL, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}

		visits = append(visits, &v)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}
	return
}

func VisitorTagsByVisitorID(db XODB, visitorID string) (tags []string, err error) {
	q := `SELECT VT.name FROM custmchat.visitor_tag VT ` +
		`INNER JOIN custmchat.visitor_tag_relation VTR ON VT.id = VTR.tag_id ` +
		`WHERE VTR.visitor_id = ?`

	rows, err := db.Query(q, visitorID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}

		tags = append(tags, name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func VisitorTagRelationsByVisitors(db XODB, visitors []*Visitor) (tags []*VisitorTagRelation, err error) {
	if len(visitors) == 0 {
		return
	}

	var visitorIDs []string
	for _, visitor := range visitors {
		visitorIDs = append(visitorIDs, visitor.ID)
	}

	return VisitorTagRelationsByVisitorIDs(db, visitorIDs)
}

func VisitorTagRelationsByVisitorIDs(db XODB, visitorIDs []string) (tags []*VisitorTagRelation, err error) {
	if len(visitorIDs) == 0 {
		return
	}

	query := `SELECT visitor_id,tag_id FROM custmchat.visitor_tag_relation WHERE visitor_id IN (%s)`

	var args []interface{}
	var ps []string
	for _, id := range visitorIDs {
		args = append(args, id)
		ps = append(ps, "?")
	}

	query = fmt.Sprintf(query, strings.Join(ps, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vt = &VisitorTagRelation{}
		if err = rows.Scan(&vt.VisitorID, &vt.TagID); err != nil {
			return nil, err
		}

		tags = append(tags, vt)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

func VisitorsByTraceIDs(db XODB, traceIDs []string) ([]*Visitor, error) {
	var err error

	if len(traceIDs) == 0 {
		return nil, nil
	}

	// sql query
	sqlstr := `SELECT ` +
		`id, ent_id, trace_id, name, age, gender, avatar, mobile, weibo, wechat, email, qq_num, address, remark, visit_cnt, visit_page_cnt, residence_time_sec, last_visit_id, visited_at, created_at, updated_at ` +
		`FROM custmchat.visitor ` +
		`WHERE trace_id IN (%s)`

	var args []interface{}
	var placeHolders []string
	for _, id := range traceIDs {
		args = append(args, id)
		placeHolders = append(placeHolders, "?")
	}

	sqlstr = fmt.Sprintf(sqlstr, strings.Join(placeHolders, ","))
	q, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Visitor{}
	for q.Next() {
		v := Visitor{}

		// scan
		err = q.Scan(&v.ID, &v.EntID, &v.TraceID, &v.Name, &v.Age, &v.Gender, &v.Avatar, &v.Mobile, &v.Weibo, &v.Wechat, &v.Email, &v.QqNum, &v.Address, &v.Remark, &v.VisitCnt, &v.VisitPageCnt, &v.ResidenceTimeSec, &v.LastVisitID, &v.VisitedAt, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &v)
	}
	if err = q.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func IncrVisitorTagUseCount(db XODB, tagID string) (err error) {
	update := `UPDATE custmchat.visitor_tag SET use_count = use_count + 1 WHERE id=?`
	_, err = db.Exec(update, tagID)
	return
}

func VisitorBlacklistsByEntID(db XODB, entID string, offset, limit int) (list []*VisitBlacklist, err error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, trace_id, visit_id, agent_id, conv_id, created_at ` +
		`FROM custmchat.visit_blacklist ` +
		`WHERE ent_id = ? ` +
		`LIMIT ?, ?`

	rows, err := db.Query(sqlstr, entID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vb := VisitBlacklist{}
		err = rows.Scan(&vb.ID, &vb.EntID, &vb.TraceID, &vb.VisitID, &vb.AgentID, &vb.ConvID, &vb.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &vb)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
