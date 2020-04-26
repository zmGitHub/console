CREATE TABLE `visit_page` (
    `id` CHAR(20) PRIMARY KEY,
    `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
    `visit_id` CHAR(20) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL COMMENT '访问id',
    `ip` VARCHAR(50) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL DEFAULT '',
    `source` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '首次来路页source',
    `source_keyword` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '首次来路页keyword',
    `source_domain` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '首次来路页domain',
    `source_url` TEXT NOT NULL COMMENT '首次来路页url',
    `title` TEXT NOT NULL COMMENT '首次着陆页title',
    `domain` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '首次着陆页domain',
    `url` TEXT NOT NULL COMMENT '首次着陆页url',
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
    `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
    INDEX `idx_ent_id_visit_id_create_at` (`ent_id`, `visit_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


