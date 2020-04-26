CREATE TABLE `evaluation` (
    `ent_id` CHAR(20) NOT NULL,
    `agent_id` CHAR(20) NOT NULL,
    `conv_id` CHAR(20) NOT NULL COMMENT '对话id',
    `level` TINYINT NOT NULL,
    `content` VARCHAR(255) NOT NULL DEFAULT '',
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    KEY `idx_ent` (`ent_id`),
    KEY `idx_agent` (`agent_id`),
    UNIQUE KEY (`conv_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
