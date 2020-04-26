create TABLE `promotion_msgs` (
  `id` CHAR(20) PRIMARY KEY,
  `enterprise_id` CHAR(20) NOT NULL,
  `source` VARCHAR(50) NOT NULL,
  `content` TEXT,
  `content_sdk` VARCHAR(50) NOT NULL,
  `countdown` INT NOT NULL DEFAULT 0,
  `enabled` BOOLEAN NOT NULL DEFAULT FALSE,
  `summary` VARCHAR(300) NOT NULL DEFAULT '',
  `thumbnail` VARCHAR(250) NOT NULL DEFAULT '',
  `created_on` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_on` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
