package allocator

import "bitbucket.org/forfd/custm-chat/webim/models"

// Allocator agent allocator
type Allocator interface {
	Allocate() (agentID string, err error)
}

type AgentInfo interface {
	GetEntAgentIDs(entID string) ([]*models.AgentRanking, error)
	LastAllocatedAgentID(entID string) (string, error)
	AgentActiveConvNum(agentID string) (int, error)
	SetLastAllocatedAgent(entID, agentID string) error
}
