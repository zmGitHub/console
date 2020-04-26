package allocator

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"fmt"
	"sort"
)

var (
	NoAvailableAgentsErr = fmt.Errorf("no onlines agents")
	AgentNotFoundErr     = fmt.Errorf("no agent found, allocation failed")
)

// TakeTurns 按坐席顺序轮流分配
type TakeTurns struct {
	EntID     string
	AgentInfo AgentInfo
}

func (tt *TakeTurns) Allocate() (string, error) {
	agentRankings, err := tt.AgentInfo.GetEntAgentIDs(tt.EntID)
	if err != nil {
		return "", err
	}

	return TakeTurnsByAgentIDs(tt.AgentInfo, tt.EntID, agentRankings)
}

func TakeTurnsByAgentIDs(agentInfo AgentInfo, entID string, onlineAgents []*models.AgentRanking) (string, error) {
	agentsLen := len(onlineAgents)
	if agentsLen == 0 {
		return "", common.NoOnlineAgents
	}

	lastAllocatedAgentID, err := agentInfo.LastAllocatedAgentID(entID)
	if err != nil {
		return "", err
	}

	var pos = 0
	if lastAllocatedAgentID != "" {
		for i, rk := range onlineAgents {
			if rk.AgentID == lastAllocatedAgentID {
				pos = i + 1
				break
			}
		}
	}

	var counts []*AgentConvCount

	result := onlineAgents[pos%agentsLen]
	convNum, err := agentInfo.AgentActiveConvNum(result.AgentID)
	if err != nil {
		return "", err
	}

	counts = append(counts, &AgentConvCount{
		ID:    result.AgentID,
		Count: convNum,
	})
	if convNum < result.ServeLimit {
		return result.AgentID, nil
	}

	pos++
	for i := 1; i <= len(onlineAgents)-1; i++ {
		rk := onlineAgents[pos%agentsLen]
		convNum, err = agentInfo.AgentActiveConvNum(rk.AgentID)
		if err != nil {
			return "", err
		}

		counts = append(counts, &AgentConvCount{
			ID:    rk.AgentID,
			Count: convNum,
		})

		if convNum < rk.ServeLimit {
			return rk.AgentID, nil
		}

		pos++
	}

	sort.SliceStable(counts, func(i, j int) bool {
		return counts[i].Count <= counts[j].Count
	})

	return counts[0].ID, common.AgentServeLimitExceed
}
