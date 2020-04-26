CREATE TABLE `personal_config` (
  `agent_id` CHAR(20) NOT NULL UNIQUE,
  `config_content` TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;