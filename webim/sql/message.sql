CREATE TABLE `message` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `trace_id` CHAR(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_id` CHAR(20) NOT NULL,
  `conversation_id` CHAR(20) NOT NULL,
  `from_type` VARCHAR(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` DATETIME(6) DEFAULT NULL,
  `read_time` DATETIME(6) DEFAULT NULL,
  `content` TEXT COLLATE utf8mb4_unicode_ci,
  `content_type` VARCHAR(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `msg_type` VARCHAR(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `extra` TEXT COLLATE utf8mb4_unicode_ci,
  KEY `idx_conversation_id` (`conversation_id`),
  KEY `msg_find_by_ent_and_trace` (`ent_id`,`trace_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
