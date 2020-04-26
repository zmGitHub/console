CREATE TABLE `ent_app`(
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `app_name` VARCHAR(50) NOT NULL,
  `create_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  `update_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
  UNIQUE KEY (`ent_id`, `app_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `perm`(
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL COMMENT '企业id',
  `app_name` VARCHAR(30) NOT NULL COMMENT '产品名称',
  `name` VARCHAR(50) NOT NULL COMMENT '权限名',
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT "修改时间",
  UNIQUE KEY (`ent_id`, `app_name`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `role_perm`(
  `role_id` CHAR(20) NOT NULL,
  `perm_id` CHAR(20) NOT NULL COMMENT '权限id',
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6) COMMENT '创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6) COMMENT '修改时间',
  UNIQUE KEY (`role_id`, `perm_id`),
  INDEX (`perm_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

