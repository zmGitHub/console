CREATE TABLE `automatic_message` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `channel_type` VARCHAR(50) NOT NULL COMMENT '渠道类型(web, sdk, ...)', /**/
  `msg_type` VARCHAR(50) NOT NULL COMMENT '推广消息，企业欢迎消息，客服无应答, 顾客无响应, ...',
  `msg_content` TEXT,
  `after_seconds` INT NOT NULL DEFAULT 0 COMMENT '多长时间(秒)之后发送',
  `enabled` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否启用',
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
