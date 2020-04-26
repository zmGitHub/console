package common

import "fmt"

var (
	NoOnlineAgents        = fmt.Errorf("No Online Agents")
	AgentServeLimitExceed = fmt.Errorf("Agent Serve Limit Exceed")
)

const (
	Ok                       = 0
	InvalidParameterErr      = 1000
	DBErr                    = 1001
	EntExistErr              = 1002
	UserNotExistErr          = 1003
	UserPasswordErr          = 1004
	RedisErr                 = 1005
	EncodeJSONErr            = 1006
	DecodeJSONErr            = 1007
	EntNotExistErr           = 1008
	GenTokenErr              = 1009
	UserLoginCountExceedErr  = 1010
	SendingMsgErr            = 1011
	AgentNumExceedErr        = 1012
	AgentNotOnlineErr        = 1013
	AgentAllocateErr         = 1014
	AgentAlreadyExistsErr    = 1015
	MessageExistsTooLong     = 1016
	UserNotActivatedErr      = 1017
	UploadFileErr            = 1018
	AgentServeLimitExceedErr = 1019
	ConversationEndedErr     = 1020
	PermissionLimited        = 1021
	ExportFileErr            = 1022
	ParseQuickReplyErr       = 1023
	UserExistsErr            = 1024

	InternalServerErr = 5000
)
