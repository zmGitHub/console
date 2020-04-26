CREATE TABLE `login_records` (
  `id` CHAR(20) PRIMARY KEY,
  `agent_id` CHAR(20) NOT NULL,
  `ent_id` CHAR(20) NOT NULL,
  `login_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '登录时间',
  `login_client` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '登录客户端',
  `login_ip` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '登录ip',
  `device_info` VARCHAR(200) NOT NULL DEFAULT '' COMMENT '登录设备信息',
  KEY `idx_ent` (`ent_id`),
  KEY `idx_agent` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
