package models

func InsertOrUpdateQueueConfig(db XODB, entID string, queueConfig *QueueConfig) (err error) {
	const sqlstr = `INSERT INTO custmchat.queue_config (ent_id, queue_size, description, status) VALUES (?, ?, ?, ?) ` +
		`ON DUPLICATE KEY UPDATE queue_size=VALUES(queue_size), description=VALUES(description), status=VALUES(status)`

	_, err = db.Exec(sqlstr, entID, queueConfig.QueueSize, queueConfig.Description, queueConfig.Status)
	return
}

func InsertOrUpdateEndConversation(db XODB, entID string, endConv *EndingConversation) (err error) {
	const sqlstr = `INSERT INTO custmchat.ending_conversation (` +
		`ent_id, no_message_duration, offline_duration, status` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE no_message_duration=VALUES(no_message_duration), offline_duration=VALUES(offline_duration), status=VALUES(status)`
	_, err = db.Exec(sqlstr, entID, endConv.NoMessageDuration, endConv.OfflineDuration, endConv.Status)
	return
}

func InsertOrUpdateEndMessage(db XODB, entID string, endMsg *EndingMessage) (err error) {
	const sqlstr = `INSERT INTO custmchat.ending_message (` +
		`ent_id, platform, agent, system, status, prompt` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE agent=VALUES(agent), system=VALUES(system), status=VALUES(status), prompt=VALUES(prompt)`
	_, err = db.Exec(sqlstr, entID, endMsg.Platform, endMsg.Agent, endMsg.System, endMsg.Status, endMsg.Prompt)
	return
}

func InsertOrUpdateConversationTransfer(db XODB, entID string, convTransfer *ConversationTransfer) (err error) {
	const sqlstr = `INSERT INTO custmchat.conversation_transfer (` +
		`ent_id, duration, transfer_target, target_type, status` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE duration=VALUES(duration), transfer_target=VALUES(transfer_target), target_type=VALUES(target_type), status=VALUES(status)`
	_, err = db.Exec(sqlstr, entID, convTransfer.Duration, convTransfer.TransferTarget, convTransfer.TargetType, convTransfer.Status)
	return
}

func InsertOrUpdateConversationQuality(db XODB, entID string, convQuality *ConversationQuality) (err error) {
	const sqlstr = `INSERT INTO custmchat.conversation_quality (` +
		`ent_id, grade, visitor_msg_count, agent_msg_count, status` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`) ` +
		`ON DUPLICATE KEY UPDATE grade=VALUES(grade), visitor_msg_count=VALUES(visitor_msg_count), agent_msg_count=VALUES(agent_msg_count), status=VALUES(status)`
	_, err = db.Exec(sqlstr, entID, convQuality.Grade, convQuality.VisitorMsgCount, convQuality.AgentMsgCount, convQuality.Status)
	return
}
