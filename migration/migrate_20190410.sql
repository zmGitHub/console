ALTER TABLE `enterprise` ADD COLUMN `admin_id` CHAR(20) NOT NULL AFTER `name`;
ALTER TABLE `enterprise` ADD COLUMN `allocation_rule` VARCHAR(100) NOT NULL DEFAULT '' AFTER `admin_id`;
ALTER TABLE `enterprise` ADD UNIQUE INDEX (`admin_id`);