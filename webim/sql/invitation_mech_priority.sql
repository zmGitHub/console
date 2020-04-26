CREATE TABLE `invitation_mechs_priority` (
	`ent_id` CHAR(20) NOT NULL COMMENT '企业ID',
	`priority` TEXT NOT NULL COMMENT '优先级配置',
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '更新时间',
	PRIMARY KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
