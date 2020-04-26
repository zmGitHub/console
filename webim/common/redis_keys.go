package common

const (
	LastAllocationAgentID = "ent:%s:last:allocated:agent"
	EntAgents             = "ent:%s:agents"
	EntConfigs            = "ent:%s:configs"
	EntOnlineAgents       = "ent:%s:online:agents"
	EntOnlineAgentList    = "ent:%s:login:agent:list" // redis sorted set

	EntOnlineAgentSet = "ent:%s:online:agents:set"

	AgentConversationNum = "agent:%s:conversation:num"
	AgentServeLimit      = "agent:%s:serve:limit"
	AgentLoginCount      = "agent:%s:login:count"
	AgentTokenList       = "agent:%s:token:list"

	AgentPerms = "agent:%s:perms"
)

var (
	EntConfigsAgentNum = "agent_num"
	EntConfigsContent  = "configs"
)
