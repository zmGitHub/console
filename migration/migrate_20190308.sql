alter table `quickreply_group` add column `creator_type` VARCHAR(50) NOT NULL after `created_by`;
alter table `quickreply_group` add column `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) after `created_at`;