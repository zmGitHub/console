package models

import (
	"database/sql"
	"time"
)

type AgentOverallConversationStats struct {
	AgentID          string
	ConvCount        int
	AvgDurationInSec int
	DurationInSec    int
	MsgCnt           int
}
type AgentEvalConversationStats struct {
	AgentID       string
	EvaConvCnt    int
	GoodConvCnt   int
	MediumConvCnt int
	NoEvaConvCnt  int
	BadConvCnt    int
}

type AgentQualityConversationStats struct {
	AgentID        string
	GoldConvCnt    int
	SilverConvCnt  int
	BronzeConvCnt  int
	NogradeConvCnt int
}

func ConversationsByTimeRange(db XODB, entID string, start, end time.Time) (conversations []*Conversation, err error) {
	query := `SELECT id, created_at, msg_count, client_msg_count, agent_msg_count, duration, eval_level, quality_grade ` +
		`FROM custmchat.conversation WHERE ent_id=? AND created_at >= ? AND created_at < ? `

	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		conv := &Conversation{}
		if err = rows.Scan(
			&conv.ID, &conv.CreatedAt, &conv.MsgCount, &conv.ClientMsgCount, &conv.AgentMsgCount,
			&conv.Duration, &conv.EvalLevel, &conv.QualityGrade); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func AgentConversationsByTimeRange(db XODB, entID string, start, end time.Time) (conversations []*Conversation, err error) {
	query := `SELECT id, agent_id, eval_level, created_at FROM custmchat.conversation WHERE ent_id=? AND created_at >= ? AND created_at < ?`
	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		conv := &Conversation{}
		if err = rows.Scan(&conv.ID, &conv.AgentID, &conv.EvalLevel, &conv.CreatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func VisitsByTimeRange(db XODB, entID string, start, end time.Time) (visits []*Visit, err error) {
	query := `SELECT id, created_at FROM custmchat.visit WHERE ent_id=? AND created_at >= ? AND created_at < ?`
	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		visit := &Visit{}
		if err = rows.Scan(&visit.ID, &visit.CreatedAt); err != nil {
			return nil, err
		}

		visits = append(visits, visit)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func VisitorByTimeRange(db XODB, entID string, start, end time.Time) (visitors []*Visitor, err error) {
	query := `SELECT id, remark, created_at FROM custmchat.visitor WHERE ent_id=? AND created_at >= ? AND created_at < ?`
	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		visitor := &Visitor{}
		if err = rows.Scan(&visitor.ID, &visitor.Remark, &visitor.CreatedAt); err != nil {
			return nil, err
		}

		visitors = append(visitors, visitor)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func VisitPageByTimeRange(db XODB, entID string, start, end time.Time) (pages []*VisitPage, err error) {
	query := `SELECT id, created_at FROM custmchat.visit_page WHERE ent_id = ? AND created_at >= ? AND created_at < ?`
	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		page := &VisitPage{}
		if err = rows.Scan(&page.ID, &page.CreatedAt); err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func EvalConversationStatsByTimeRange(db XODB, agentIDs []string, start, end time.Time) (stats map[string]*AgentEvalConversationStats, err error) {
	query := `select count(agent_id) as conv_count, eval_level from conversation where agent_id=? and created_at >= ? and created_at < ? group by eval_level;`

	stats = map[string]*AgentEvalConversationStats{}
	for _, id := range agentIDs {
		rows, err := db.Query(query, id, start, end)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var stat = &AgentEvalConversationStats{AgentID: id}
		for rows.Next() {
			var convCount, level int
			if err = rows.Scan(&convCount, &level); err != nil {
				return nil, err
			}

			switch level {
			case NoEvalConversation:
				stat.NoEvaConvCnt = convCount
			case GoodConversation:
				stat.GoodConvCnt = convCount
			case MediumConversation:
				stat.MediumConvCnt = convCount
			case BadConversation:
				stat.BadConvCnt = convCount
			}
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		stats[id] = stat
	}

	return
}

func QualityConversationStatsByTimeRange(db XODB, agentIDs []string, start, end time.Time) (stats map[string]*AgentQualityConversationStats, err error) {
	query := `select count(agent_id) as conv_count, quality_grade from conversation where agent_id=? and created_at >= ? and created_at < ? group by quality_grade`

	stats = map[string]*AgentQualityConversationStats{}
	for _, id := range agentIDs {
		rows, err := db.Query(query, id, start, end)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var stat = &AgentQualityConversationStats{AgentID: id}
		for rows.Next() {
			var convCount int
			var qualityGrade sql.NullString
			if err = rows.Scan(&convCount, &qualityGrade); err != nil {
				return nil, err
			}

			switch qualityGrade.String {
			case FirstLevel:
				stat.GoldConvCnt = convCount
			case SecondLevel:
				stat.SilverConvCnt = convCount
			case ThirdLevel:
				stat.BronzeConvCnt = convCount
			default:
				stat.NogradeConvCnt = convCount
			}
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		stats[id] = stat
	}
	return
}

func OverallConversationStatsByTimeRange(db XODB, entID string, start, end time.Time) (stats map[string]*AgentOverallConversationStats, err error) {
	query := `select agent_id, count(agent_id) as conv_count, sum(msg_count) as mag_num, sum(duration) as total_duration ` +
		`from conversation where ent_id=? and created_at >= ? and created_at < ? group by agent_id;`
	rows, err := db.Query(query, entID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats = map[string]*AgentOverallConversationStats{}
	for rows.Next() {
		var stat = &AgentOverallConversationStats{}
		if err = rows.Scan(&stat.AgentID, &stat.ConvCount, &stat.MsgCnt, &stat.DurationInSec); err != nil {
			return
		}

		if stat.ConvCount > 0 {
			stat.AvgDurationInSec = stat.DurationInSec / stat.ConvCount
		}
		stats[stat.AgentID] = stat
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
