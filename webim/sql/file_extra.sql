CREATE TABLE `file_extra` (
  `ent_id` CHAR(20) NOT NULL,
  `name` VARCHAR(300) NOT NULL,
  `type` VARCHAR(50) NOT NULL DEFAULT 'file',
  `size` INT NOT NULL DEFAULT 0,
  `upload_time` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `expire_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  KEY `idx_ent` (`ent_id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;