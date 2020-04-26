CREATE TABLE `invitation_mechs_rules` (
	`id` CHAR(20) PRIMARY KEY  COMMENT '规则ID',
	`ent_id` CHAR(20) NOT NULL COMMENT '企业ID',
	`rule` TEXT NOT NULL COMMENT '规则配置(条件)',
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '更新时间',
	INDEX `idx_enterprise_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
