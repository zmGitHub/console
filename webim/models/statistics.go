package models

import (
	"database/sql"
	"time"
)

type ConversationStats struct {
	ConvCount          int64   `json:"conv_count"`
	EffectiveConvCount int64   `json:"effective_conv_count"`
	MsgCount           int64   `json:"msg_count"`
	RespDuration       float64 `json:"resp_duration"`
	ConvDuration       float64 `json:"conv_duration"`
}

func getDate() time.Time {
	now := time.Now().UTC()
	year, month, day := now.Date()
	hour := now.Hour()
	return time.Date(year, month, day, hour, 0, 0, 0, time.UTC)
}

func GetVisitStatsByDateRange(db XODB, entID string, start, end time.Time) (visitorCount, visitNum int64, err error) {
	query := `SELECT sum(visitor_count) AS visitor_count, sum(visit_num) AS visit_num` +
		`FROM visitor_statistic WHERE ent_id = ? AND DATE(created_at) BETWEEN ? AND ?`

	var visitorCountVal, visitNumVal, pageViewsVal sql.NullInt64
	if err = db.QueryRow(query, entID, start, end).Scan(&visitorCountVal, &visitNumVal, &pageViewsVal); err != nil {
		return 0, 0, err
	}

	if visitorCountVal.Valid {
		visitorCount = visitorCountVal.Int64
	}

	if visitNumVal.Valid {
		visitNum = visitNumVal.Int64
	}

	return
}

func GetConversationStatsByDateRange(db XODB, entID string, start, end time.Time) (stats *ConversationStats, err error) {
	query := `SELECT sum(conversation_count) AS conv_count, sum(effective_conversation_count) AS effect_conv_count, ` +
		`SUM(message_count) AS msg_count, ` +
		`SUM(avg_resp_duration) AS resp_duration, SUM(avg_conversation_duration) AS conversation_duration ` +
		`FROM custmchat.conversation_statistic WHERE ent_id = ? AND DATE(created_at) BETWEEN ? AND ?`

	stats = &ConversationStats{}
	var convCount, effectiveConvCount, msgCount sql.NullInt64
	var respDuration, convDuration sql.NullFloat64
	if err = db.QueryRow(query, entID, start, end).Scan(&convCount, &effectiveConvCount, &msgCount, &respDuration, &convDuration); err != nil {
		return nil, err
	}

	if convCount.Valid {
		stats.ConvCount = convCount.Int64
	}

	if effectiveConvCount.Valid {
		stats.EffectiveConvCount = effectiveConvCount.Int64
	}

	if msgCount.Valid {
		stats.MsgCount = msgCount.Int64
	}

	if respDuration.Valid {
		stats.RespDuration = respDuration.Float64
	}

	if convDuration.Valid {
		stats.ConvDuration = convDuration.Float64
	}

	return
}
