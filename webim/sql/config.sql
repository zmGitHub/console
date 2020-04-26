/* 排队配置 */
CREATE TABLE `queue_config` (
  `ent_id` CHAR(20) NOT NULL,
  `queue_size` INT NOT NULL DEFAULT 0 COMMENT '排队访客数量',
  `description`  VARCHAR(255) NOT NULL COMMENT '排队提示文案',
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 什么条件下结束对话 */
CREATE TABLE `ending_conversation` (
  `ent_id` CHAR(20) NOT NULL,
  `no_message_duration` INT NOT NULL DEFAULT -1, /*对话无消息多长时间后结束*/
  `offline_duration` INT NOT NULL DEFAULT -1,    /*访客离线多长时间后结束*/
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 对话结束消息 */
CREATE TABLE `ending_message` (
  `ent_id` CHAR(20) NOT NULL,
  `platform` VARCHAR(10) NOT NULL COMMENT 'web, sdk, wechat, weibo',
  `agent` VARCHAR(500) NOT NULL DEFAULT '',
  `system` VARCHAR(500) NOT NULL DEFAULT '',
  `prompt` TINYINT(1) NOT NULL DEFAULT 0,
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  UNIQUE KEY (`ent_id`, `platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*  访客首条消息无应答，对话转接 */
CREATE TABLE `conversation_transfer` (
  `ent_id` CHAR(20) NOT NULL,
  `duration` INT NOT NULL DEFAULT 30, /*30秒无应答，转接*/
  `transfer_target` CHAR(20) NOT NULL,
  `target_type` VARCHAR(10) NOT NULL, /*group or auto*/
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 对话质量标准 */
CREATE TABLE `conversation_quality` (
  `ent_id` CHAR(20) NOT NULL,
  `grade` VARCHAR(50) NOT NULL,
  `visitor_msg_count` INT NOT NULL DEFAULT 0,
  `agent_msg_count` INT NOT NULL DEFAULT 0,
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 登录限制 */
CREATE TABLE `login_limit` (
  `ent_id` CHAR(20) NOT NULL UNIQUE,
  `group_ids` TEXT,
  `city_list` TEXT,
  `allowed_ip_list` TEXT,
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 发送文件开启 */
CREATE TABLE `send_file` (
  `ent_id` CHAR(20) NOT NULL UNIQUE,
  `status` BOOLEAN NOT NULL DEFAULT FALSE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/* 留言配置
  "qq": "close",
    "wechat": "open",
    "name": "open",
    "intro": "您好，现在为下班时间，请您留下个人信息，我们会在第一时间跟您取得联系。",
    "permission": "open",
    "captcha": "close",
    "tel": "open",
    "defaultTemplate": "open",
    "contactRule": "single",
    "email": "open",
    "category": "open",
    "defaultTemplateContent": ""
*/
CREATE TABLE `leave_message_config` (
  `ent_id` CHAR(20) NOT NULL,
  `introduction` VARCHAR(500) NOT NULL DEFAULT '',
  `show_visitor_name` BOOLEAN NOT NULL DEFAULT FALSE,
  `show_telephone` BOOLEAN NOT NULL DEFAULT FALSE,
  `show_email` BOOLEAN NOT NULL DEFAULT FALSE,
  `show_wechat` BOOLEAN NOT NULL DEFAULT FALSE,
  `show_qq` BOOLEAN NOT NULL DEFAULT FALSE,
  `auto_create_category` BOOLEAN NOT NULL DEFAULT FALSE,
  `fill_contact` VARCHAR(20) NOT NULL DEFAULT 'single', /* single, multi */
  `use_default_content` BOOLEAN NOT NULL DEFAULT FALSE,
  `default_content` VARCHAR(200) NOT NULL DEFAULT '',
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ent_all_configs` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL UNIQUE,
  `config_content` TEXT,
  `create_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `update_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `personal_config` (
  `agent_id` CHAR(20) NOT NULL UNIQUE,
  `config_content` TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


create TABLE `promotion_msgs` (
  `id` CHAR(20) PRIMARY KEY,
  `enterprise_id` CHAR(20) NOT NULL,
  `source` VARCHAR(50) NOT NULL,
  `content` TEXT,
  `content_sdk` VARCHAR(50) NOT NULL,
  `countdown` INT NOT NULL DEFAULT 0,
  `enabled` BOOLEAN NOT NULL DEFAULT FALSE,
  `summary` VARCHAR(300) NOT NULL DEFAULT '',
  `thumbnail` VARCHAR(250) NOT NULL DEFAULT '',
  `created_on` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_on` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
  KEY `idx_enterprise` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
