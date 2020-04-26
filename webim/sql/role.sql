CREATE TABLE `role` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `perms_range_groups` (
  `agent_id` CHAR(20) NOT NULL,
  `group_id` CHAR(20) NOT NULL,
  UNIQUE KEY (`agent_id`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

