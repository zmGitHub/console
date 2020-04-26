CREATE TABLE `agent_group`(
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  `description` VARCHAR(200) NOT NULL DEFAULT '',
  UNIQUE KEY `ent_id_name` (`ent_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `agent_group_relation` (
  `group_id` CHAR(20) NOT NULL,
  `uid` CHAR(20) NOT NULL,
  UNIQUE KEY (`uid`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
