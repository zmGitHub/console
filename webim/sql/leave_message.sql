CREATE TABLE `leave_message` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `track_id` CHAR(20) NOT NULL,
  `last_option_agent` CHAR(20) NULL,
  `mobile` VARCHAR(32) NOT NULL DEFAULT '',
  `content` TEXT,
  `status` VARCHAR(20) NOT NULL DEFAULT 'open',
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
  KEY (`ent_id`),
  KEY (`track_id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
