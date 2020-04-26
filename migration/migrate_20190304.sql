alter table agent add column `public_email` VARCHAR(255) NOT NULL DEFAULT '' after `mobile`;
alter table agent add column `public_telephone` VARCHAR(32)  NOT NULL DEFAULT '' after `public_email`;
