CREATE TABLE `visit` (
    `id` CHAR(20) PRIMARY KEY,
    `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
    `trace_id` CHAR(20) NOT NULL,
    `visit_page_cnt` INT NOT NULL DEFAULT 1 COMMENT '访问页数',
    `residence_time_sec` INT NOT NULL DEFAULT 0 COMMENT '停留秒数',

    `browser_family` VARCHAR(50) NOT NULL DEFAULT '',
    `browser_version_string` VARCHAR(50) NOT NULL DEFAULT '',
    `browser_version` VARCHAR(50) NOT NULL DEFAULT '',
    `os_category` VARCHAR(50) NOT NULL DEFAULT '',
    `os_family` VARCHAR(50) NOT NULL DEFAULT '',
    `os_version_string` VARCHAR(50) NOT NULL DEFAULT '',
    `os_version` VARCHAR(50) NOT NULL DEFAULT '',
    `platform` VARCHAR(50) NOT NULL DEFAULT '',
    `ua_string` TEXT NOT NULL,

    `ip` VARCHAR(50) CHARACTER SET ASCII COLLATE ASCII_BIN NOT NULL DEFAULT '',
    `country` VARCHAR(30) NOT NULL DEFAULT '',
    `province` VARCHAR(30) NOT NULL DEFAULT '',
    `city` VARCHAR(30) NOT NULL DEFAULT '',
    `isp` VARCHAR(20) NOT NULL DEFAULT '',

    `first_page_source` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '首次来路页source',
    `first_page_source_keyword` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '首次来路页keyword',
    `first_page_source_domain` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '首次来路页domain',
    `first_page_source_url` TEXT NOT NULL COMMENT '首次来路页url',
    `first_page_title` TEXT NOT NULL COMMENT '首次着陆页title',
    `first_page_domain` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '首次着陆页domain',
    `first_page_url` TEXT NOT NULL COMMENT '首次着陆页url',

    `latest_title` TEXT NOT NULL COMMENT '最新着陆页title',
    `latest_url` TEXT NOT NULL COMMENT '最新着陆页url',

    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
    `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
    INDEX `idx_enterprise_id_trace_id` (`ent_id`, `trace_id`),
    INDEX `idx_enterprise_id_created_on` (`ent_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

