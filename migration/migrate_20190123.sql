ALTER TABLE `conversation` ADD COLUMN `agent_effective_msg_count` INT NOT NULL DEFAULT 0;
ALTER TABLE `conversation` ADD COLUMN  `client_last_send_time` DATETIME(6);
ALTER TABLE `conversation` ADD COLUMN  `first_msg_create_time` DATETIME(6);
ALTER TABLE `conversation` ADD COLUMN  `eval_content` VARCHAR(500) NOT NULL DEFAULT '';
ALTER TABLE `conversation` ADD COLUMN  `eval_level` INT NOT NULL DEFAULT 0;
ALTER TABLE `conversation` ADD COLUMN  `has_summary` BOOLEAN NOT NULL DEFAULT FALSE;