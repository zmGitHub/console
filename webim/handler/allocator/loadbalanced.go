package allocator

import (
	"math"
	"sort"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
)

type AgentConvCount struct {
	ID    string
	Count int
}

// LoadBalanced
type LoadBalanced struct {
	EntID     string
	AgentInfo AgentInfo
}

func (lb *LoadBalanced) Allocate() (string, error) {
	agentRankings, err := lb.AgentInfo.GetEntAgentIDs(lb.EntID)
	if err != nil {
		log.Logger.Warnf("Get Ent Agents: %v", err)
		return "", err
	}

	if len(agentRankings) == 0 {
		return "", common.NoOnlineAgents
	}

	var agentID string
	var minConvNum = math.MaxInt16
	var counts []*AgentConvCount

	for _, rk := range agentRankings {
		convNum, err := lb.AgentInfo.AgentActiveConvNum(rk.AgentID)
		if err != nil {
			return "", err
		}

		counts = append(counts, &AgentConvCount{
			ID:    rk.AgentID,
			Count: convNum,
		})

		if convNum < rk.ServeLimit && convNum < minConvNum {
			minConvNum = convNum
			agentID = rk.AgentID
		}
	}

	if agentID != "" {
		return agentID, nil
	}

	sort.SliceStable(counts, func(i, j int) bool {
		return counts[i].Count <= counts[j].Count
	})

	return counts[0].ID, common.AgentServeLimitExceed
}
