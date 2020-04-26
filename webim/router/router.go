package router

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/handler"
	imid "bitbucket.org/forfd/custm-chat/webim/middleware"
)

func SetRouter(e *echo.Echo, service *handler.IMService) {
	dumpConfig := middleware.BodyDumpConfig{
		Skipper: func(context echo.Context) bool {
			req := context.Request()
			p := req.URL.Path
			if p == "" {
				p = "/"
			}
			if strings.HasPrefix(p, "/api/v1/upload") {
				return true
			}

			if strings.HasPrefix(p, "/upload") {
				return true
			}
			if strings.HasPrefix(p, "/api/v1/upload_img") {
				return true
			}

			if strings.HasPrefix(p, "/captcha_images") {
				return true
			}

			if strings.HasPrefix(p, "/api/v1/status/sync") {
				return true
			}

			if strings.HasPrefix(p, "/metrics") {
				return true
			}

			return false
		},
		Handler: func(context echo.Context, req []byte, resp []byte) {
			log.Logger.WithField("req", string(req)).WithField("resp", string(resp)).Info()
		},
	}
	e.Use(middleware.Recover())
	e.Use(imid.Log())
	e.Use(middleware.BodyDumpWithConfig(dumpConfig))

	e.Use(imid.NewMetric())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	setSuperAdminEndpoints(e, service)
	setClientEndpoints(e, service)
	setFrontendAdapterRouter(e, service)

	e.POST("/signin", service.SignIn)
	e.POST("/reset", service.ResetPassword)
	e.POST("/forget_password", service.ForgetPassword)
	e.POST("/captchas", service.CreateCaptchas)
	e.GET("/register_redirect", service.VerifyRegister)
	e.POST("/register_direct", service.RegisterEnterprise)

	api := e.Group("/api/v1")
	api.POST("/reset_password", service.ResetForgetPassword)
	api.GET("/allocate_agent", service.AllocateAgent)
	api.GET("/activate", service.ActivateUser)
	api.GET("/activate_email", service.ActivateAgentEmail)
	api.GET("/activate_ent", service.ActivateEnterprise)
	api.GET("/connection_token", service.GenVisitorConnectionToken)
	api.POST("/upload", service.Upload)
	api.POST("/upload_img", service.UploadImg)
	api.PUT("/status/sync", service.SyncStatus)

	sysGroup := api.Group("/system")
	sysGroup.POST("/end_conversation", service.SysEndConversation)
	sysGroup.POST("/offline_end_conversation", service.SysOfflineEndConversation)
	sysGroup.POST("/send_ent_message", service.SendEntMessage)
	sysGroup.POST("/send_no_resp_message", service.SendNoRespMessage)

	entGroup := api.Group("/enterprises/:ent_id")
	entGroup.GET("/visitor_allowed", service.CheckVisitorAllowed)
	entGroup.POST("/visits", service.InitVisit)
	entGroup.POST("/visits/:visit_id/update_residence", service.UpdateVisitResidenceTimeSec)
	entGroup.POST("/leave_messages", service.CreateLeaveMessage)
	entGroup.GET("/chat_forms", service.GetEntForms)

	entConversationGroup := entGroup.Group("/conversations")
	entConversationGroup.POST("/:conversation_id/evaluations", service.CreateEvaluation)
	entConversationGroup.GET("/history", service.GetHistoryConvs)

	adminAPI := e.Group("/admin/api/v1")
	adminAPI.Use(imid.AgentAuth(conf.IMConf.JWTKey))
	adminAPI.Use(imid.RefreshOnline(service.RefreshAgentOnline))

	adminAPI.POST("/reset_password", service.ResetPassword)
	adminAPI.GET("/check_online_count", service.CheckOnlineAgentCount)

	adminAPI.POST("/allocation_rules", service.CreateOrUpdateAllocationRule)
	adminAPI.GET("/active_conversations", service.GetActiveConversations)
	adminAPI.GET("/colleague_conversations", service.GetColleagueConversations)

	egroup := adminAPI.Group("/enterprise")
	egroup.GET("/configs", service.GetEntConfigs)
	egroup.PUT("/info", service.UpdateEntInfo)
	egroup.GET("/agents", service.GetEntAgents)
	egroup.GET("/agent_groups", service.GetEntAgentGroups)
	egroup.POST("/agent_groups", service.CreateAgentGroup)
	egroup.GET("/visits", service.GetEntVisits)
	egroup.PUT("/conv_configs", service.UpdateEntConvConfigs)
	egroup.GET("/security_configs", service.SecurityConfigByEntID)
	egroup.POST("/security_configs", service.AddSecurityConfig)
	egroup.PUT("/security_configs", service.UpdateSecurityConfig)
	egroup.GET("/reports/visit", service.VisitReport)
	egroup.GET("/reports/conversation", service.ConversationReport)

	agentGroup := adminAPI.Group("/agents")
	agentGroup.POST("/", service.AddAgent)
	agentGroup.GET("/export", service.ExportAgents)
	agentGroup.GET("/info", service.GetCurrentAgentInfo)
	agentGroup.GET("/:agent_id", service.GetAgentByID)
	agentGroup.PUT("/:agent_id", service.UpdateAgentInfo)
	agentGroup.DELETE("/:agent_id", service.DeleteAgent)
	agentGroup.GET("/connection_token", service.GenConnectionToken)
	agentGroup.POST("/message_beep", service.CreateOrUpdateMessageBeep) // 消息提醒
	agentGroup.GET("/message_beep", service.GetMessageBeepConfig)
	agentGroup.GET("/:agent_id/perms", service.GetAgentPerms)
	agentGroup.PUT("/status", service.UpdateAgentStatus)
	agentGroup.PUT("/rankings", service.UpdateAgentsRanking)

	visitorGroup := adminAPI.Group("/visitors")
	visitorGroup.GET("/", service.OnlineVisitors)
	visitorGroup.PUT("/", service.UpdateVisitorInfo)
	visitorGroup.POST("/tags", service.AddTagToVisitor)
	visitorGroup.POST("/search", service.SearchVisitors)

	quickReplyGroup := adminAPI.Group("/quickreplies")
	quickReplyGroup.GET("/", service.GetQuickReplies)
	quickReplyGroup.POST("/", service.CreateQuickReplyGroup)
	quickReplyGroup.GET("/export", service.ExportQuickReply)
	quickReplyGroup.POST("/import", service.ImportQuickReply)
	quickReplyGroup.PUT("/:group_id", service.UpdateQuickReplyGroup)
	quickReplyGroup.DELETE("/:group_id", service.DeleteQuickReplyGroup)
	quickReplyGroup.POST("/:group_id/items", service.CreateQuickReplyItem)
	quickReplyGroup.PUT("/items/:item_id", service.UpdateQuickReplyItem)
	quickReplyGroup.DELETE("/items/:reply_id", service.DeleteQuickReplyItem)

	autoMsgGroup := adminAPI.Group("/automessages")
	autoMsgGroup.Use(middleware.RemoveTrailingSlash())
	autoMsgGroup.GET("/", service.GetEntAutoMessages)
	autoMsgGroup.POST("/", service.CreateAutoMessage)
	autoMsgGroup.PUT("/:msg_id", service.UpdateAutoMessage)
	autoMsgGroup.DELETE("/:msg_id", service.DeleteAutoMessage)

	preFormGroup := adminAPI.Group("/prechat_forms")
	preFormGroup.POST("/", service.AddForm)
	preFormGroup.PUT("/:form_id", service.UpdateForm)

	visitorTag := adminAPI.Group("/visitor_tags")
	visitorTag.POST("/", service.CreateVisitorTag)
	visitorTag.GET("/", service.GetEntVisitorTags)
	visitorTag.PUT("/:tag_id", service.UpdateVisitorTag)
	visitorTag.DELETE("/:tag_id", service.DeleteVisitorTag)

	blackListGroup := adminAPI.Group("/blacklists")
	blackListGroup.POST("/", service.AddVisitorToBlacklist)

	conversationGroup := adminAPI.Group("/conversations")
	conversationGroup.PUT("/:conversation_id/transfer", service.TransferConversation)
	conversationGroup.PUT("/:conversation_id/summary", service.AddSummary)

	msgGroup := conversationGroup.Group("/:conversation_id/messages")
	msgGroup.GET("/", service.MessagesByConversationID)

	roleGroup := adminAPI.Group("/roles")
	roleGroup.POST("/", service.AddRole)
	roleGroup.POST("/:role_id/perms", service.AddOrUpdatePermsToRole)
	roleGroup.GET("/:role_id/perms", service.GetRolePerms) // get role perm ids

	permGroup := adminAPI.Group("/perms")
	permGroup.GET("/", service.GetPerms)

	adminAPI.POST("/leave_message_config", service.CreateLeaveMessageConfig)

	lmGroup := adminAPI.Group("/leave_messages")
	lmGroup.GET("/", service.LeaveMessagesByEnt)
	lmGroup.PUT("/:msg_id", service.UpdateLeaveMessageStatus)
}

func setSuperAdminEndpoints(e *echo.Echo, service *handler.IMService) {
	group := e.Group("/super_admin")
	group.Use(imid.SuperAdminAuth(conf.IMConf.APIToken))
	group.GET("/enterprises", service.QueryAllEnterprise)
	group.POST("/enterprises", service.RegisterEnterprise)
	group.POST("/trials", service.StartTrial)
	group.DELETE("/enterprises/tests/:ent_id", service.DeleteTestEnterprise)
	group.PUT("/enterprises/upgrade", service.UpgradeEnterprise)
}

func setClientEndpoints(e *echo.Echo, service *handler.IMService) {
	e.GET("/visit/init", service.InitVisitV1)
	e.POST("/visit/:ent_id/:track_id/reject", service.RejectVisit)
	e.GET("/api/client/custom_field", service.GetClientCustomFields)
	e.GET("/api/client/attr_order", service.GetClientAttrOrder)
	e.POST("/scheduler", service.Scheduler)
	e.GET("/client/timeline", service.TimeLine)
	e.POST("/client/send_msg", service.ClientSendMsg)
	e.GET("/client/forms", service.GetClientForms) // client/forms
	e.POST("/client/forms", service.FillForms)
	e.GET("/client/history_conversation", service.GetHistoryConversation)
	e.GET("/client/queue/position", service.GetVisitorPosition)
	e.GET("/api/conversation/:track_id/streams", service.GetStreamConversations)
	e.POST("/conversation/:conversation_id/evaluation", service.EvaluateConversation)
	e.POST("/api/leave_messages", service.CreateLeaveMessage)
}

func setFrontendAdapterRouter(e *echo.Echo, service *handler.IMService) {
	e.POST("/bulletin-pull/gpul", service.GPUL)
	e.POST("/agent/accept_invitations", service.AcceptInvitation)
	e.GET("/captcha_images/:captcha_id", service.GetCaptcha)
	e.POST("/api/agent/confirm_email_change", service.ConfirmAgentEmail)
	e.POST("/api/request_reset", service.RequestResetPwd)
	e.POST("/upload", service.AdminUpload)
	e.GET("/qrcode", service.QrCode)

	e.GET("/api/agent/quick_replies/export", service.ExportQuickReply)

	adapterAPI := e.Group("/api")
	adapterAPI.Use(imid.AgentAuth(conf.IMConf.JWTKey))
	adapterAPI.Use(imid.RefreshOnline(service.RefreshAgentOnline))

	adapterAPI.POST("/resend_activate_account_email", service.ResendActivateEmail)
	adapterAPI.POST("/logout", service.Logout)
	adapterAPI.POST("/captchas", service.CreateCaptchas)
	adapterAPI.GET("/leave_messages", service.LeaveMessagesByEnt)
	adapterAPI.PUT("/leave_messages/:msg_id", service.UpdateLeaveMessageStatus)

	adapterAPI.GET("/enterprise", service.GetEnterpriseInfo)
	adapterAPI.PUT("/enterprise", service.UpdateEntInfo)
	adapterAPI.GET("/enterprise/plan", service.GetEnterprisePlan)
	adapterAPI.GET("/enterprise/features", service.GetEnterpriseFeatures)
	adapterAPI.PUT("/enterprise/allocation_rule", service.CreateOrUpdateAllocationRule)
	adapterAPI.GET("/enterprise/promotion_msgs", service.GetPromotionMessages)
	adapterAPI.POST("/enterprise/promotion_msgs", service.SetPromotionMessages)
	adapterAPI.PUT("/enterprise/promotion_msgs/:msg_id", service.UpdatePromotionMessages)
	adapterAPI.DELETE("/enterprise/promotion_msgs/:msg_id", service.DeletePromotionMessage)

	adapterAPI.GET("/perm/roles", service.GetEntRoles)
	adapterAPI.POST("/perm/roles", service.AddRole)
	adapterAPI.PUT("/perm/roles/:role_id", service.UpdateRole)
	adapterAPI.DELETE("/perm/roles/:role_id", service.DeleteRole)
	adapterAPI.GET("/perm/roles/:role_id/perms", service.GetRolePerms)
	adapterAPI.PUT("/perm/roles/:role_id/perms", service.AddOrUpdatePermsToRole)
	adapterAPI.GET("/perm/perms", service.GetEnterprisePerms)

	adapterAPI.POST("/conversation/new", service.NewConversation)
	adapterAPI.GET("/conversation/search/v2", service.SearchConversations)
	adapterAPI.GET("/conversation/agent/active", service.GetActiveConversations)
	adapterAPI.GET("/conversation/agent/colleagues", service.GetColleagueConversations)
	adapterAPI.PUT("/conversation/:conversation_id/summary", service.AddSummary)
	adapterAPI.POST("/conversation/:conversation_id/invite_evaluation", service.InviteEval)

	adapterAPI.GET("/stats/agents", service.GetAgentsStats)
	adapterAPI.GET("/feature", service.GetEnterpriseFeatures)

	adapterAPI.GET("/client/:track_id/pages", service.GetConversationPages)

	adapterAPI.PUT("/visit/:visitor_id", service.UpdateVisitorName)
	adapterAPI.GET("/visit/:visit_id/pages", service.GetVisitorPages)
	adapterAPI.GET("/visit/search", service.OnlineVisits)
	adapterAPI.POST("/visit/:track_id/invite", service.InviteVisitor)
	adapterAPI.GET("/client/:track_id/attrs", service.GetVisitorInfo)
	adapterAPI.POST("/client/:track_id/attrs", service.UpdateVisitorInfo)

	adapterAPI.POST("/users/search", service.SearchVisitors)
	adapterAPI.GET("/users/segments", service.UserSegments)
	adapterAPI.PUT("/users/segments/:seg_id", service.UpdateUserSegments)
	adapterAPI.DELETE("/users/segments/:seg_id", service.DeleteSegments)
	adapterAPI.POST("/users/segments", service.CreateUserSegments)
	adapterAPI.GET("/users/config", service.UsersConfig)

	adapterAPI.POST("/analytics/overview", service.OverView)
	adapterAPI.POST("/analytics/agent_online_detail", service.AgentOnlineDetail)
	adapterAPI.POST("/analytics/agent/workload", service.AgentWorkLoad)
	adapterAPI.POST("/stats/conv", service.ConversationStats)
	adapterAPI.POST("/stats/evaluation", service.EvaluationStats)
	adapterAPI.POST("/stats/traffic_overview", service.TrafficOverview)

	adapterAgentGroup := adapterAPI.Group("/agent")
	adapterAgentGroup.GET("/online", service.SendOnlineEvent)
	adapterAgentGroup.GET("/online_agents", service.OnlineAgentsInfo)
	adapterAgentGroup.GET("/queue", service.QueueVisitors)
	adapterAgentGroup.GET("/get_visit_filter", service.GetVisitFilter)
	adapterAgentGroup.POST("/conversation/:conversation_id/group_redirect", service.TransferConversation)
	adapterAgentGroup.POST("/conversation/:conversation_id/redirect", service.TransferToAgent)

	adapterAgentGroup.POST("/kick_person_offline", service.KickOffline)

	adapterAgentGroup.GET("/enterprise_config", service.GetEntConfigsV1)
	adapterAgentGroup.POST("/enterprise_config", service.CreateOrUpdateEntConfigsV1)
	adapterAgentGroup.POST("/enterprise_config/change_msg_status", service.ChangeMsgStatus)

	adapterAgentGroup.GET("/info", service.GetCurrentAgentInfo)
	adapterAgentGroup.GET("/agent_invitations", service.GetInvitations)
	adapterAgentGroup.POST("/agent_invitations", service.SendInvitation)
	adapterAgentGroup.PUT("/agent_invitations/:invitation_id", service.CancelInvitation)
	adapterAgentGroup.POST("/resend_invitations/:invitation_id", service.ResendInvitation)

	adapterAgentGroup.GET("/agents", service.GetEntAgents)
	adapterAgentGroup.PUT("/agents/:agent_id", service.UpdateAgentInfoV1)
	adapterAgentGroup.DELETE("/agents/:agent_id", service.DeleteAgent)
	adapterAgentGroup.PUT("/agents/:agent_id/status", service.UpdateAgentStatus)
	adapterAgentGroup.GET("/agents/:agent_id/perms_ranges", service.GetAgentPerms)
	adapterAgentGroup.GET("/agents/:agent_id/mentioned_convs", service.GetMentionedConvs)

	adapterAgentGroup.POST("/client/:track_id/tags", service.AddTagToVisitor)
	adapterAgentGroup.DELETE("/client/:track_id/tags/:tag_id", service.DeleteTagFromVisitor)
	adapterAgentGroup.GET("/client_tags", service.GetEntVisitorTags)
	adapterAgentGroup.POST("/client_tags", service.CreateVisitorTag)
	adapterAgentGroup.PUT("/client_tags/:tag_id", service.UpdateVisitorTag)
	adapterAgentGroup.DELETE("/client_tags/:tag_id", service.DeleteVisitorTag)

	adapterAgentGroup.GET("/agent_groups", service.GetEntAgentGroups)
	adapterAgentGroup.POST("/agent_groups", service.CreateAgentGroup)
	adapterAgentGroup.DELETE("/agent_groups/:group_id", service.DeleteAgentGroup)
	adapterAgentGroup.PUT("/agent_groups/:group_id", service.UpdateAgentGroup)

	adapterAgentGroup.GET("/selecting_rules", service.GetSelectingRules)

	adapterAgentGroup.GET("/quick_replies", service.GetQuickReplies)
	adapterAgentGroup.POST("/quick_replies", service.CreateQuickReplyItem)
	adapterAgentGroup.POST("/quick_replies/import", service.ImportQuickReply)
	adapterAgentGroup.PUT("/quick_replies/:reply_id", service.UpdateQuickReplyItem)
	adapterAgentGroup.DELETE("/quick_replies/:reply_id", service.DeleteQuickReplyItem)
	adapterAgentGroup.POST("/quick_reply_groups", service.CreateQuickReplyGroup)
	adapterAgentGroup.PUT("/quick_reply_groups/:group_id", service.UpdateQuickReplyGroup)
	adapterAgentGroup.DELETE("/quick_reply_groups/:group_id", service.DeleteQuickReplyGroup)

	adapterAgentGroup.GET("/sales_config", service.GetEnterpriseSalesConfig)
	adapterAgentGroup.GET("/online_agents_info", service.OnlineAgentsInfo)
	adapterAgentGroup.GET("/personal_config", service.GetAgentPersonalConfigs)
	adapterAgentGroup.POST("/personal_config", service.UpdateAgentPersonalConfigs)

	adapterAgentGroup.GET("/visit/blacklist", service.GetVisitBlackList)
	adapterAgentGroup.PUT("/visit/blacklist", service.AddVisitorToBlacklistV1)
	adapterAgentGroup.DELETE("/visit/blacklist", service.DeleteVisitBlacklist)
	adapterAgentGroup.POST("/send_msg", service.AgentSendMsgV1)
	adapterAgentGroup.POST("/send_internal_msg", service.SendInternalMsg)
	adapterAgentGroup.POST("/end_conversation", service.EndConversationV1)

	adapterAgentGroup.GET("/forms", service.GetEntForms)
	adapterAgentGroup.PUT("/forms/:form_id", service.UpdateForm)
}
