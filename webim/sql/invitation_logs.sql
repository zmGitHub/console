CREATE TABLE `invitation_logs` (
	`id` CHAR(20) PRIMARY KEY,
	`ent_id` CHAR(20) NOT NULL COMMENT '企业ID',
	`trace_id` CHAR(20) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL,
	`look_id` CHAR(20) NOT NULL COMMENT '外观ID',
	`look_config_index` INT NOT NULL COMMENT '所属规则中外观配置中的索引',
	`mech_id` CHAR(20) NOT NULL COMMENT '机制ID',
	`mech_config_index` INT NOT NULL COMMENT '所属规则中外观配置中的索引',
	`is_accepted` BOOL NOT NULL DEFAULT FALSE,
	`conversation_id` CHAR(20),
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '更新时间',
	INDEX `idx_ent_id_created_at` (`ent_id`, `created_at`),
	INDEX `idx_ent_id_look_id_created_at` (`ent_id`, `look_id`, `created_at`),
	INDEX `idx_ent_id_mech_id_created_at` (`ent_id`, `mech_id`, `created_at`),
	INDEX `idx_ent_id_trace_id_created_at` (`ent_id`, `trace_id`, `created_at`),
	INDEX `idx_ent_id_look_id_config_index_created_at` (`ent_id`, `look_id`, `look_config_index`, `created_at`),
	INDEX `idx_ent_id_mech_id_config_index_created_at` (`ent_id`, `mech_id`, `mech_config_index`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
