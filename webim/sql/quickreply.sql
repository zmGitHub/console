CREATE TABLE `quickreply_group` (
	`id` CHAR(20) PRIMARY KEY,
	`ent_id` CHAR(20) NOT NULL,
	`title` VARCHAR(100) NOT NULL COMMENT '标题',
	`rank` INT NOT NULL DEFAULT 0,
	`created_by` CHAR(20) NOT NULL,
	`creator_type` VARCHAR(50) NOT NULL,
	`position` INT NOT NULL DEFAULT 0,
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
	KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `quickreply_item` (
	`id` CHAR(20) PRIMARY KEY,
	`quickreply_group_id` CHAR(20) NOT NULL,
	`title` VARCHAR(100) NOT NULL,
	`content` TEXT COLLATE utf8mb4_unicode_ci,
	`content_type` VARCHAR(50) NOT NULL DEFAULT '',
	`rich_content` TEXT,
	`rank` INT NOT NULL DEFAULT 0,
	`hot_key` VARCHAR(500) NOT NULL DEFAULT '[]',
	`created_by` CHAR(20) NOT NULL,
	`created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
	`updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
	KEY `idx_quickreply_group` (`quickreply_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
