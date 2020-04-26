CREATE TABLE `visit_blacklist` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
  `trace_id` CHAR(20) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL,
  `visit_id` CHAR(20) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL DEFAULT '',
  `agent_id` CHAR(20) NOT NULL COMMENT '坐席id',
  `conv_id` CHAR(20) NOT NULL COMMENT '对话id',
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  UNIQUE KEY (`ent_id`,`trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

