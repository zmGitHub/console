CREATE TABLE `invitation_looks` (
	`rule_id` CHAR(20) PRIMARY KEY,
	`ent_id` CHAR(20) NOT NULL COMMENT '企业ID',
	`name` VARCHAR(30) NOT NULL COMMENT '名称',
	`target` VARCHAR(30) NOT NULL COMMENT '平台类型',
	`begin_on` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '开始时间',
	`expire_on` DATETIME(6) COMMENT '过期时间,如果为null,则为永久',
	`enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用',
	`style_config` TEXT NOT NULL COMMENT '样式配置',
	`is_any` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否任意匹配',
	`version` VARCHAR(30) NOT NULL DEFAULT  '' COMMENT '版本号',
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '更新时间',
	INDEX `idx_enterprise_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
