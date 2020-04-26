package allocator

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriority_Allocate(t *testing.T) {
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
					"agent1": {1, 2},
					"agent2": {0, 2},
					"agent3": {0, 2},
					"agent4": {0, 2},
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
					"agent2": {0, 2},
					"agent3": {0, 2},
					"agent4": {0, 2},
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
					"agent1": {2, 2},
					"agent2": {2, 2},
					"agent3": {2, 2},
					"agent4": {0, 2},
				},
			},
			result: "agent4",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       onlineAgents[:2],
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {2, 2},
					"agent2": {2, 2},
					"agent3": {0, 2},
					"agent4": {0, 2},
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
					"agent2": {0, 2},
					"agent3": {0, 2},
					"agent4": {0, 2},
				},
			},
			result: "",
			err:    common.NoOnlineAgents,
		},
	}

	ast := assert.New(t)

	for _, c := range cases {
		p := &Priority{EntID: "", AgentInfo: c.mockAgent}
		result, err := p.Allocate()
		ast.Equal(c.result, result)
		ast.Equal(c.err, err)
	}
}
