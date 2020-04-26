CREATE TABLE `visitor_tag` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `creator` CHAR(20) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `color` VARCHAR(50) NOT NULL,
  `use_count` INT NOT NULL DEFAULT 0,
  `rank` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `visitor_tag_relation` (
  `visitor_id` CHAR(20) NOT NULL,
  `tag_id` CHAR(20) NOT NULL,
  UNIQUE KEY (`visitor_id`, `tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
