CREATE TABLE `agent_statistic` (
  `ent_id` CHAR(20) NOT NULL,
  `agent_id` CHAR(20) NOT NULL,
  `conversation_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `good_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `medium_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `bad_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `message_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `duration` INT NOT NULL DEFAULT 0,
  `first_resp_duration` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  KEY (`ent_id`),
  UNIQUE (`agent_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `visitor_statistic` (
  `ent_id` CHAR(20) NOT NULL,
  `visitor_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `visit_num`  INT UNSIGNED NOT NULL DEFAULT 0,
  `page_views` INT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
  KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `conversation_statistic` (
  `ent_id` CHAR(20) NOT NULL,
  `conversation_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `effective_conversation_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `message_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `avg_resp_duration` FLOAT NOT NULL DEFAULT 0.0,
  `avg_conversation_duration` FLOAT NOT NULL DEFAULT 0.0,
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
  KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `conversation_stats` (
    `ent_id` CHAR(20) NOT NULL,
    `total_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `effective_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `message_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `good_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `medium_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `bad_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `gold_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `silver_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `bronze_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `no_grade_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `duration_in_sec` INT UNSIGNED NOT NULL DEFAULT 0,
    `remark_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `visit_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `visit_page_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `visitor_count` INT UNSIGNED NOT NULL DEFAULT 0,
    `wait_time_in_sec` INT UNSIGNED NOT NULL DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    UNIQUE `ent_create_Time` (`ent_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
