package allocator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

var onlineAgents = []*models.AgentRanking{
	{
		AgentID:    "agent1",
		Ranking:    0,
		ServeLimit: 2,
	},
	{
		AgentID:    "agent2",
		Ranking:    0,
		ServeLimit: 2,
	},
	{
		AgentID:    "agent3",
		Ranking:    0,
		ServeLimit: 2,
	},
	{
		AgentID:    "agent4",
		Ranking:    0,
		ServeLimit: 2,
	},
}

func TestLoadBalanced_Allocate(t *testing.T) {
	cases := []struct {
		mockAgent *mockAgentInfo
		result    string
		err       error
	}{
		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents,
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {0},
					"agent3": {0},
					"agent4": {0},
				},
			},
			result: "agent2",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents,
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {1},
					"agent3": {0},
					"agent4": {0},
				},
			},
			result: "agent3",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents,
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {1},
					"agent3": {1},
					"agent4": {0},
				},
			},
			result: "agent4",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents,
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {1, 2},
					"agent2": {2, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "agent1",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents,
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {2, 2},
					"agent2": {2, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "",
			err:    common.AgentServeLimitExceed,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       []*models.AgentRanking{},
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {2, 2},
					"agent2": {2, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "",
			err:    common.NoOnlineAgents,
		},
	}

	ast := assert.New(t)

	for _, c := range cases {
		loadBalanced := &LoadBalanced{EntID: "", AgentInfo: c.mockAgent}
		result, err := loadBalanced.Allocate()
		ast.Equal(c.result, result)
		ast.Equal(c.err, err)
	}
}
