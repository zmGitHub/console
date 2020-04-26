CREATE TABLE `online_agents` (
    `ent_id` CHAR(20) NOT NULL,
    `agent_id` CHAR(20) NOT NULL,
    `online_count` INT NOT NULL DEFAULT 0,
    `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    UNIQUE KEY `idx_ent_agent` (`ent_id`, `agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `online_visitors` (
    `ent_id` CHAR(20) NOT NULL,
    `trace_id` CHAR(20) NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `idx_ent_trace` (`ent_id`, `trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
