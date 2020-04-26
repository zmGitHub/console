package models

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/imclient"
)

var (
	agentChannelSuffix = "_message"

	AgentAvailableStatus   = "on_duty"
	AgentUnavailableStatus = "off_duty"

	AgentAccountCreatedStatus = "created"
	AgentAccountValidStatus   = "valid"
	AgentAccountInValidStatus = "invalid"
	AgentAccountDeletedStatus = "deleted"

	AgentPermsRangePersonalType = "self"
	AgentPermsRangeAllType      = "all"
	AgentPermsRangePartType     = "part"

	AgentInvitationPendingStatus   = "pending"
	AgentInvitationSuccessStatus   = "success"
	AgentInvitationCancelledStatus = "cancelled"
	AgentInvitationExpiredStatus   = "expired"

	AgentPermsRangeType = map[string]string{
		AgentPermsRangePersonalType: "个人",
		AgentPermsRangeAllType:      "所有",
		AgentPermsRangePartType:     "部分",
	}

	defaultPingDuration = 30 * time.Second
)

type AgentRanking struct {
	AgentID    string
	Ranking    int
	ServeLimit int
}

func UpdateAgentAccountStatusValid(db XODB, id string) error {
	sqlStr := `UPDATE custmchat.agent SET account_status = ? WHERE id = ?`
	_, err := db.Exec(sqlStr, AgentAccountValidStatus, id)
	return err
}

func GetAgentPermsRange(db XODB, id string) (permsRange string, err error) {
	query := `SELECT perms_range_type FROM custmchat.agent WHERE id=?`
	if err = db.QueryRow(query, id).Scan(&permsRange); err != nil {
		return
	}

	return
}
func OnDutyAgents(db XODB, entID string) (rankings []*AgentRanking, err error) {
	sqlStr := `SELECT id, ranking, serve_limit FROM custmchat.agent WHERE ent_id=? AND status=? AND deleted_at IS NULL ORDER BY ranking`

	rows, err := db.Query(sqlStr, entID, AgentAvailableStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ranking = &AgentRanking{}
		if err = rows.Scan(&ranking.AgentID, &ranking.Ranking, &ranking.ServeLimit); err != nil {
			return nil, err
		}

		rankings = append(rankings, ranking)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rankings, nil
}

func AgentRankingsByEntID(db XODB, entID string) (rankings []*AgentRanking, err error) {
	sqlStr := `SELECT id, ranking FROM custmchat.agent WHERE ent_id=? AND account_status=? AND deleted_at IS NULL ORDER BY ranking`

	rows, err := db.Query(sqlStr, entID, AgentAccountValidStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var ranking int
		if err = rows.Scan(&id, &ranking); err != nil {
			return nil, err
		}

		rankings = append(rankings, &AgentRanking{AgentID: id, Ranking: ranking})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rankings, nil
}

func AgentNumOfEnt(db XODB, entID string) (agentNum int, err error) {
	sqlStr := `SELECT count(id) FROM custmchat.agent WHERE ent_id=? AND deleted_at IS NULL`

	if err = db.QueryRow(sqlStr, entID).Scan(&agentNum); err != nil {
		return -1, err
	}

	return
}

func AgentServeLimitByID(db XODB, id string) (limit int, err error) {
	// sql query
	const sqlstr = `SELECT serve_limit FROM custmchat.agent WHERE id = ?`
	err = db.QueryRow(sqlstr, id).Scan(&limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return -1, err
	}

	return
}

func ServeLimitByAgentIDs(db XODB, ids []string) (limits map[string]int, err error) {
	if len(ids) == 0 {
		return
	}

	var args []interface{}
	var ps []string
	for _, id := range ids {
		args = append(args, id)
		ps = append(ps, "?")
	}

	query := `SELECT id, serve_limit FROM custmchat.agent WHERE id IN (%s)`
	query = fmt.Sprintf(query, strings.Join(ps, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	limits = map[string]int{}
	for rows.Next() {
		var agentID string
		var serveLimit int
		if err = rows.Scan(&agentID, &serveLimit); err != nil {
			return
		}

		limits[agentID] = serveLimit
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func AgentsByMsgs(db XODB, msgs []*Message) (agents []*Agent, err error) {
	if len(msgs) == 0 {
		return
	}

	var agentIDs []string
	for _, msg := range msgs {
		agentIDs = append(agentIDs, msg.AgentID)
	}

	return AgentsByAgentIDs(db, agentIDs)
}

func AgentsByAgentIDs(db XODB, ids []string) (agents []*Agent, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// sql query
	const sqlstr = `SELECT ` +
		`id, ent_id, group_id, role_id, avatar, username, real_name, nick_name, hashed_password, job_number, serve_limit, is_online, ranking, email, mobile, public_email, public_telephone, qq_num, signature, status, wechat, is_admin, perms_range_type, account_status, create_at, update_at, deleted_at ` +
		`FROM custmchat.agent ` +
		`WHERE id in (%s)`

	var args []interface{}
	var placeHolders []string
	for _, id := range ids {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}

	q, err := db.Query(fmt.Sprintf(sqlstr, strings.Join(placeHolders, ",")), args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var res []*Agent
	for q.Next() {
		a := Agent{}

		err = q.Scan(&a.ID, &a.EntID, &a.GroupID, &a.RoleID, &a.Avatar, &a.Username, &a.RealName, &a.NickName, &a.HashedPassword, &a.JobNumber, &a.ServeLimit, &a.IsOnline, &a.Ranking, &a.Email, &a.Mobile, &a.PublicEmail, &a.PublicTelephone, &a.QqNum, &a.Signature, &a.Status, &a.Wechat, &a.IsAdmin, &a.PermsRangeType, &a.AccountStatus, &a.CreateAt, &a.UpdateAt, &a.DeletedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}
	if err := q.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateAgentStatus(db XODB, id, status string) error {
	sqlStr := `UPDATE custmchat.agent SET status = ? WHERE id = ?`
	_, err := db.Exec(sqlStr, status, id)
	return err
}

func IsAgentOnDuty(db XODB, id string) (bool, error) {
	var status string
	sqlStr := `SELECT status FROM custmchat.agent WHERE id = ?`
	err := db.QueryRow(sqlStr, id).Scan(&status)
	if err != nil {
		return false, err
	}

	return status == AgentAvailableStatus, nil
}

func UpdateAgentsRanking(db XODB, rankings map[string]int) error {
	if len(rankings) == 0 {
		return nil
	}
	const sqlstr = `INSERT INTO custmchat.agent (id, ranking) VALUES %s ON DUPLICATE KEY UPDATE ranking=VALUES(ranking)`

	var args []interface{}
	var holders []string
	for agentID, ranking := range rankings {
		holders = append(holders, "(?, ?)")
		args = append(args, agentID)
		args = append(args, ranking)
	}

	_, err := db.Exec(fmt.Sprintf(sqlstr, strings.Join(holders, ",")), args...)
	return err
}

func UpdateAgentPassword(db XODB, id, password string) error {
	sqlStr := `UPDATE custmchat.agent SET hashed_password = ? WHERE id = ?`
	_, err := db.Exec(sqlStr, password, id)
	return err
}

func UpdateAgentPasswordByEmail(db XODB, password, email string) error {
	sqlStr := `UPDATE custmchat.agent SET hashed_password = ? WHERE email = ?`
	_, err := db.Exec(sqlStr, password, email)
	return err
}

func IsAgentExistsByPassword(db XODB, id, password string) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM custmchat.agent WHERE id=? AND hashed_password=?)`

	if err := db.QueryRow(query, id, password).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func UpdateAgentEmail(db XODB, id, email string) error {
	sqlStr := `UPDATE custmchat.agent SET email = ? WHERE id = ?`
	_, err := db.Exec(sqlStr, email, id)
	return err
}

func DeleteAgent(db XODB, id string) error {
	update := `UPDATE custmchat.agent SET deleted_at = ? WHERE id = ?`
	if _, err := db.Exec(update, time.Now().UTC(), id); err != nil {
		return err
	}
	return nil
}

func IsAdmin(db XODB, agentID string) (bool, error) {
	var isAdmin int
	query := "SELECT is_admin FROM custmchat.agent WHERE id = ?"
	err := db.QueryRow(query, agentID).Scan(&isAdmin)
	if err != nil {
		return false, err
	}

	return isAdmin == 1, nil
}

func ChangeAdmin(db XODB, agentID string, newAdminID string) error {
	update := `UPDATE custmchat.agent SET is_admin = ? WHERE id=?`

	if _, err := db.Exec(update, 0, agentID); err != nil {
		return err
	}
	if _, err := db.Exec(update, 1, newAdminID); err != nil {
		return err
	}

	return nil
}

func AgentIDsByGroupIDs(db XODB, groupIDs []string) (agentIDs []string, err error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	query := `SELECT id FROM custmchat.agent WHERE group_id IN (%s) AND deleted_at IS NULL`

	var args []interface{}
	var ps []string
	for _, groupID := range groupIDs {
		args = append(args, groupID)
		ps = append(ps, "?")
	}

	query = fmt.Sprintf(query, strings.Join(ps, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Logger.Errorf("db rows close error: %v", closeErr)
		}
	}()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return
		}

		agentIDs = append(agentIDs, id)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func AgentIDsByEntID(db XODB, entID string) (agentIDs []string, err error) {
	query := `SELECT id from custmchat.agent WHERE ent_id=? AND deleted_at IS NULL`

	rows, err := db.Query(query, entID)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Logger.Errorf("db rows close error: %v", closeErr)
		}
	}()

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return
		}

		agentIDs = append(agentIDs, id)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func AgentIDsByRoleID(db XODB, roleID string) (ids []string, err error) {
	const sqlstr = `SELECT id FROM custmchat.agent WHERE role_id = ?`
	q, err := db.Query(sqlstr, roleID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	for q.Next() {
		var id string
		err = q.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return
}

func IncrOnlineAgentCount(db XODB, entID, agentID string, ping bool) error {
	q := `INSERT INTO online_agents (ent_id, agent_id, online_count, updated_at) VALUES (?,?,?,?) ` +
		`ON DUPLICATE KEY UPDATE %s updated_at = VALUES(updated_at)`
	if !ping {
		q = fmt.Sprintf(q, "online_count = online_count + 1, ")
	} else {
		q = fmt.Sprintf(q, "")
	}

	_, err := db.Exec(q, entID, agentID, 1, time.Now().UTC())
	return err
}

func DecrOnlineAgentCount(db XODB, entID, agentID string) error {
	q := `INSERT INTO online_agents (ent_id, agent_id, online_count, updated_at) VALUES (?,?,?,?) ` +
		`ON DUPLICATE KEY UPDATE online_count = online_count - 1`
	_, err := db.Exec(q, entID, agentID, 0, time.Now().UTC())
	return err
}

type AgentInfo struct {
	Mysql XODB
}

// GetEntAgentIDs ...
func (agent *AgentInfo) GetEntAgentIDs(entID string) (rks []*AgentRanking, err error) {
	agentRankings, err := OnDutyAgents(agent.Mysql, entID)
	if err != nil {
		return nil, err
	}

	if len(agentRankings) == 0 {
		return []*AgentRanking{}, nil
	}

	var ids []string
	for _, rk := range agentRankings {
		ids = append(ids, rk.AgentID)
	}

	agentIDs := FilterOnline(ids)

Loop:
	for _, agentID := range agentIDs {
		for _, rk := range agentRankings {
			if rk.AgentID == agentID {
				rks = append(rks, rk)
				continue Loop
			}
		}

	}

	sort.SliceStable(rks, func(i, j int) bool {
		return rks[i].Ranking <= rks[j].Ranking
	})

	return rks, nil
}

func (agent *AgentInfo) LastAllocatedAgentID(entID string) (string, error) {
	key := fmt.Sprintf(common.LastAllocationAgentID, entID)
	res, err := db.RedisClient.Get(key).Result()
	if err != nil && err == redis.Nil {
		return "", nil
	}

	return res, err
}

func (agent *AgentInfo) SetLastAllocatedAgent(entID, agentID string) error {
	key := fmt.Sprintf(common.LastAllocationAgentID, entID)
	return db.RedisClient.Set(key, agentID, 0).Err()
}

func (agent *AgentInfo) AgentActiveConvNum(agentID string) (convNum int, err error) {
	convNum, err = ConversationNumOfAgent(agent.Mysql, agentID)
	if err != nil {
		return
	}

	return
}

func (agent *AgentInfo) getAgentServeLimit(agentID string) (int, error) {
	limit, err := AgentServeLimitByID(agent.Mysql, agentID)
	if err != nil {
		return -1, err
	}
	return limit, nil
}

func FilterOnline(agentIDs []string) (res []string) {
	onlineUsers := GetEntOnlineAgents()
	users := common.NonRepeatElements(onlineUsers)

	for _, id := range agentIDs {
		for _, user := range users {
			if user == id {
				res = append(res, id)
				break
			}
		}
	}
	return
}

func FilterOnlineV1(entID string, onDutyAgents []string) (res []string) {
	onlineUsers := GetEntOnlineAgentsV1(entID)
	users := common.NonRepeatElements(onlineUsers)

Loop:
	for _, id := range onDutyAgents {
		for _, user := range users {
			if user == id {
				res = append(res, id)
				continue Loop
			}
		}
	}
	return
}

func GetEntOnlineAgents() []string {
	chans, err := imclient.CentriClient.ActiveChannels(context.Background())
	if err != nil {
		log.Logger.Warnf("imCli.ActiveChannels error: %v", err)
	}

	if len(chans) == 0 {
		return []string{}
	}

	var agents []string
	for _, channel := range chans {
		if strings.HasSuffix(channel, agentChannelSuffix) {
			agents = append(agents, strings.TrimSuffix(channel, agentChannelSuffix))
		}
	}

	return agents
}

func GetEntOnlineAgentsV1(entID string) (agentIDs []string) {
	agents, err := OnlineAgentsByEntID(db.Mysql, entID)
	if err != nil {
		log.Logger.Warnf("[GetEntOnlineAgentsV1] error: %v", err)
		return []string{}
	}

	for _, agent := range agents {
		agentIDs = append(agentIDs, agent.AgentID)
	}
	return
}

func UpdateAgentInvitationStatus(db XODB, id, status string) error {
	update := `UPDATE agent_invitation SET status = ? WHERE id = ?`
	if _, err := db.Exec(update, status, id); err != nil {
		return err
	}

	return nil
}

func IsAgentOnlineV1(agentID string) bool {
	res := FilterOnline([]string{agentID})
	return len(res) == 1 && res[0] == agentID
}
