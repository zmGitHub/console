ALTER TABLE quickreply_item ADD COLUMN `content_type` VARCHAR(50) NOT NULL DEFAULT '' AFTER `content`;
ALTER TABLE quickreply_item ADD COLUMN `rich_content` TEXT AFTER `content_type`;
ALTER TABLE quickreply_item ADD COLUMN `rank` INT NOT NULL DEFAULT 0 AFTER `rich_content`;
ALTER TABLE quickreply_item ADD COLUMN `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) AFTER `created_at`;

ALTER TABLE quickreply_group ADD COLUMN `rank` INT NOT NULL DEFAULT 0 AFTER `title`;