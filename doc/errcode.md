# error code

| code         | 意义   | 
| ------------ | ------ |
| 0     | OK | 
| 1000   | InvalidParameterErr(参数不合法) |
| 1001  | DBErr(数据库错误) |
| 1002  | EntExistErr(租户已存在) |
| 1003  | UserNotExistErr(用户不存在) |
| 1004  | UserPasswordErr(用户密码错误) |
| 1005  | RedisErr(Redis 错误) |
| 1006  | EncodeJSONErr(json编码错误) |
| 1007  | DecodeJSONErr(json解码错误) |
| 1008  | EntNotExistErr(租户不存在) |
| 1009  | GenTokenErr(生成jwt错误) |
| 1010  | UserLoginCountExceedErr(在线用户数超过限制) |
| 1011  | SendingMsgErr(发送消息错误) |
| 1012  | AgentNumExceedErr(添加坐席数超过购买数) |
| 1013  | AgentNotOnlineErr(被转发的坐席不在线) |
| 1014  | AgentAllocateErr(坐席分配错误) |
| 1015  | AgentAlreadyExistsErr(被添加的坐席账户已存在(邮箱重复)) |
| 1016  | MessageExistsTooLong(超过1分钟的消息不能撤回) |
| 1017  | UserNotActivatedErr(用户未验证) |
| 1018  | UploadFileErr(上传文件错误) |
| 1019  | AgentServeLimitExceedErr(坐席最大会话数超过限制) |
| 1020  | ConversationEndedErr(会话已结束) |
| 1021  | PermissionLimited(没有权限) |
| 1022  | ExportFileErr(导出文件错误) |
| 1023  | ParseQuickReplyErr(解析快捷回复错误) |
| 5000  | InternalServerErr(服务内部错误) |
 