package dto

var APPPerms = map[string][]string{
	"visitor_and_conv": {
		"change_agent_online_status",    // : '允许修改客服在线状态',
		"invite_visitor_to_chat",        // : '允许邀请访客进行对话',
		"redirect_others_conv",          // : '允许转接他人的对话',
		"allow_others_redirect_my_conv", // : '允许他人转接自己的对话',
		"end_others_conv",               // : '允许手动结束他人的对话',
		"add_del_client_tag",            // : '允许添加或删除顾客标签',
		"add_del_client_card",           // : '允许添加或修改顾客名片',
		"add_del_client_blacklist",      // : '允许拉黑/洗白顾客',
	},

	"history_conv": {
		"see_others_history_conv", // : '允许查看他人的历史对话',
		"export_history_conv",     // : '允许导出历史对话',
	},

	"engage": {
		"see_engage", // : '允许使用顾客名片系统',
	},

	"data_report": {
		"see_data_report",      // : '允许使用数据报表系统',
		"see_conv_visit_conv",  // 允许查看 「对话-访问量和对话量报表」
		"see_conv_client_tag",  // 对话-顾客标签报表
		"see_agent_service",    // 允许查看 「客服-服务量报表」
		"see_agent_evaluation", // 允许查看 「客服-客服评价报表」
		"see_agent_status",     // 允许查看 「客服-客服状态报表」
		"manage_visit_card",    // : '允许查看和修改顾客名片管理',
		"manage_visit_list",    // : '允许查看和修改顾客列表管理',
	},

	"ent_info": {
		"config_ent_team",               // : '允许查看和修改企业信息',
		"config_ent_payment",            // : '允许查看和修改企业付费信息',
		"see_groups",                    // : '允许查看企业客服分组信息',
		"create_change_delete_group",    // : '允许创建,修改和删除客服分组',
		"create_change_delete_agent",    // : '允许创建,修改和删除客服账号（和名片）',
		"config_role",                   // : '允许查看和修改企业角色设置',
		"config_security",               // : '允许查看和修改企业安全设置',
		"config_visible_visitor_region", // : '允许查看和修改 「访客地区隔离」 设置',
	},

	"access_config": {
		"config_web_widget", // ": '允许查看和修改 「网站插件」 设置',
		"config_chat_link",  // : '允许查看和修改 「聊天链接」 设置',
	},

	"online_agent_config": {
		"config_tag",              // "允许查看和修改 「标签」 设置",
		"config_ent_quick_reply",  // "允许查看和修改企业的 「快捷回复」 设置",
		"config_ent_auto_message", // "允许查看和修改企业的 「自动消息」 设置",
		"config_agent_allocation", // "允许查看和修改 「客服分配」 设置",
		"config_invitation",       // "允许查看和修改 「对话邀请」 设置",
		"config_chat_rule",        // "允许查看和修改 「对话规则」 设置",
		"config_evaluation",       // "允许查看和修改 「客服评价」 设置",
		"config_blacklist",        // "允许查看和管理 「黑名单」",
		"config_queuing",          // "允许查看和修改 「顾客排队」 设置",
		"config_prechat_survey",   // "允许查看和修改 「询前表单」 设置",
	},
}
