CREATE TABLE `ent_all_configs` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL UNIQUE,
  `config_content` TEXT,
  `create_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `update_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
