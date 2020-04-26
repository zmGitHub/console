ALTER TABLE `online_agents` DROP INDEX `idx_agent`;

ALTER TABLE `online_agents` ADD COLUMN `ent_id` CHAR(20) NOT NULL FIRST;

ALTER TABLE `online_agents` ADD UNIQUE INDEX `idx_ent_agent` (`ent_id`, `agent_id`);
