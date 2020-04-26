package allocator

import (
	"sort"

	"bitbucket.org/forfd/custm-chat/webim/common"
)

// Priority 按照坐席排序顺序分配，当排在前面的坐席达到服务上限之后，分配下一个
type Priority struct {
	EntID     string
	AgentInfo AgentInfo
}

func (p *Priority) Allocate() (string, error) {
	agentRankings, err := p.AgentInfo.GetEntAgentIDs(p.EntID)
	if err != nil {
		return "", err
	}

	if len(agentRankings) == 0 {
		return "", common.NoOnlineAgents
	}

	var counts []*AgentConvCount

	for _, rk := range agentRankings {
		convNum, err := p.AgentInfo.AgentActiveConvNum(rk.AgentID)
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
	}

	sort.SliceStable(counts, func(i, j int) bool {
		return counts[i].Count <= counts[j].Count
	})

	return counts[0].ID, common.AgentServeLimitExceed
}
