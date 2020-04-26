CREATE TABLE `visitor`(
    `id` CHAR(20) PRIMARY KEY,
    `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
    `trace_id` CHAR(20) NOT NULL COMMENT 'trace id',
    `name` VARCHAR(128) NOT NULL COMMENT '名称',
    `age` INT NOT NULL DEFAULT 0,
    `gender` VARCHAR(10) NOT NULL DEFAULT '',
    `avatar` TEXT NOT NULL COMMENT '头像',
    `mobile` VARCHAR(50) NOT NULL DEFAULT '',
    `weibo` VARCHAR(50) NOT NULL DEFAULT '',
    `wechat` VARCHAR(50) NOT NULL  DEFAULT '',
    `email`  VARCHAR(50) NOT NULL  DEFAULT '',
    `qq_num` VARCHAR(20) NOT NULL DEFAULT '',
    `address` VARCHAR(200) NOT NULL DEFAULT '',
    `remark` VARCHAR(200)NOT NULL  DEFAULT '' COMMENT '备注',
    `visit_cnt` INT NOT NULL DEFAULT 1 COMMENT '累计访问次数',
    `visit_page_cnt` INT NOT NULL DEFAULT 1 COMMENT '累计访问页数',
    `residence_time_sec` INT NOT NULL DEFAULT 0 COMMENT '累计停留时长',
    `last_visit_id` CHAR(20) NOT NULL COMMENT '最近访问id',
    `visited_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '最近访问时间戳',
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
    `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
    UNIQUE KEY `udx_ent_id_trace_id` (`ent_id`, `trace_id`),
    INDEX `idx_ent_id_created_at` (`ent_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `visitor_attrs` (
    `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
    `trace_id` CHAR(20) NOT NULL COMMENT 'trace id',
    `attrs` TEXT NOT NULL COMMENT 'attrs',
    UNIQUE KEY `udx_ent_id_trace_id` (`ent_id`, `trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `visitor_queue` (
    `ent_id` CHAR(20) NOT NULL,
    `track_id` CHAR(20) NOT NULL UNIQUE,
    `visit_id` CHAR(20) NOT NULL,
    `enqueue_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    INDEX `idx_ent` (`ent_id`)
)