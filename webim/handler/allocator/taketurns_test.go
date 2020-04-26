package allocator

import (
	"bitbucket.org/forfd/custm-chat/webim/common"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

type mockAgentInfo struct {
	OnlineAgents       []*models.AgentRanking
	LastAllocatedAgent string
	AgentConvNumMap    map[string][]int
}

func (ag *mockAgentInfo) GetEntAgentIDs(entID string) ([]*models.AgentRanking, error) {
	return ag.OnlineAgents, nil
}

func (ag *mockAgentInfo) LastAllocatedAgentID(entID string) (string, error) {
	return ag.LastAllocatedAgent, nil
}

func (ag *mockAgentInfo) AgentActiveConvNum(agentID string) (int, error) {
	vals := ag.AgentConvNumMap[agentID]
	return vals[0], nil
}

func (ag *mockAgentInfo) SetLastAllocatedAgent(entID, agentID string) error {
	ag.LastAllocatedAgent = agentID
	return nil
}

func TestTakeTurnsByAgentIDs(t *testing.T) {
	cases := []struct {
		mockAgent *mockAgentInfo
		result    string
		err       error
	}{
		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {1},
					"agent3": {1},
					"agent4": {1},
				},
			},
			result: "agent1",
			err:    nil,
		},
		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "agent1",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {2},
					"agent3": {1},
					"agent4": {1},
				},
			},
			result: "agent1",
			err:    nil,
		},
		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "agent4",
				AgentConvNumMap: map[string][]int{
					"agent1": {1},
					"agent2": {2},
					"agent3": {1},
					"agent4": {1},
				},
			},
			result: "agent1",
			err:    nil,
		},
		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "agent4",
				AgentConvNumMap: map[string][]int{
					"agent1": {2},
					"agent2": {2},
					"agent3": {2},
					"agent4": {2},
				},
			},
			result: "",
			err:    common.AgentServeLimitExceed,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "agent3",
				AgentConvNumMap: map[string][]int{
					"agent1": {1, 2},
					"agent2": {1, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "agent1",
			err:    nil,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents:       []*models.AgentRanking{},
				LastAllocatedAgent: "agent3",
				AgentConvNumMap: map[string][]int{
					"agent1": {1, 2},
					"agent2": {1, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "",
			err:    common.NoOnlineAgents,
		},

		{
			mockAgent: &mockAgentInfo{
				OnlineAgents: []*models.AgentRanking{
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
				},
				LastAllocatedAgent: "agent2",
				AgentConvNumMap: map[string][]int{
					"agent1": {1, 2},
					"agent2": {1, 2},
					"agent3": {2, 2},
					"agent4": {2, 2},
				},
			},
			result: "agent1",
			err:    nil,
		},
	}

	ast := assert.New(t)

	for _, c := range cases {
		result, err := TakeTurnsByAgentIDs(c.mockAgent, "", c.mockAgent.OnlineAgents)
		ast.Equal(c.result, result)
		ast.Equal(c.err, err)
	}
}
