ALTER TABLE `enterprise` ADD COLUMN `nick_name` VARCHAR(100) NOT NULL DEFAULT '' AFTER `full_name`;
ALTER TABLE `enterprise` ADD COLUMN `contact_name` VARCHAR(255) NOT NULL DEFAULT '' AFTER `contact_signature`;