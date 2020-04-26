
CREATE TABLE `user_segments` (
    `id` CHAR(20) NOT NULL PRIMARY KEY,
    `ent_id` CHAR(20) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `rules` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
    KEY  `idx_ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
