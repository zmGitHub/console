CREATE TABLE `plan` (
  `id` CHAR(20) PRIMARY KEY,
  `plan_type` TINYINT NOT NULL DEFAULT 0,
  `agent_serve_limit` INT NOT NULL DEFAULT 0,
  UNIQUE (`plan_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ent_plan` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL UNIQUE,
  `plan_type` TINYINT NOT NULL DEFAULT 0,
  `trial_status` INT NOT NULL DEFAULT 0,
  `agent_serve_limit` INT NOT NULL DEFAULT 0,
  `login_agent_limit` INT NOT NULL DEFAULT 0,
  `agent_num` INT NOT NULL DEFAULT 0,
  `pay_amount` INT NOT NULL DEFAULT 0,
  `expiration_time` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `create_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  `update_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
