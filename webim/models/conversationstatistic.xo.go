// Package models contains the types for schema 'custmchat'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"time"
)

// ConversationStatistic represents a row from 'custmchat.conversation_statistic'.
type ConversationStatistic struct {
	EntID                      string    `json:"ent_id"`                       // ent_id
	ConversationCount          uint      `json:"conversation_count"`           // conversation_count
	EffectiveConversationCount uint      `json:"effective_conversation_count"` // effective_conversation_count
	MessageCount               uint      `json:"message_count"`                // message_count
	AvgRespDuration            float32   `json:"avg_resp_duration"`            // avg_resp_duration
	AvgConversationDuration    float32   `json:"avg_conversation_duration"`    // avg_conversation_duration
	CreatedAt                  time.Time `json:"created_at"`                   // created_at
	UpdatedAt                  time.Time `json:"updated_at"`                   // updated_at
}

// ConversationStatisticsByEntID retrieves a row from 'custmchat.conversation_statistic' as a ConversationStatistic.
//
// Generated from index 'ent_id'.
func ConversationStatisticsByEntID(db XODB, entID string) ([]*ConversationStatistic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`ent_id, conversation_count, effective_conversation_count, message_count, avg_resp_duration, avg_conversation_duration, created_at, updated_at ` +
		`FROM custmchat.conversation_statistic ` +
		`WHERE ent_id = ?`

	// run query
	XOLog(sqlstr, entID)
	q, err := db.Query(sqlstr, entID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ConversationStatistic{}
	for q.Next() {
		cs := ConversationStatistic{}

		// scan
		err = q.Scan(&cs.EntID, &cs.ConversationCount, &cs.EffectiveConversationCount, &cs.MessageCount, &cs.AvgRespDuration, &cs.AvgConversationDuration, &cs.CreatedAt, &cs.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &cs)
	}

	return res, nil
}
