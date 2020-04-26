/*
querySQL := `
	SELECT id, type, url_type, match_type,
			match_string, match_rules, targets, inverted
	FROM selecting_rule
	WHERE enterprise_id = ? ORDER BY rank ASC`
*/

CREATE TABLE `select_agent_rule` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `rank` TINYINT DEFAULT 0,
  `type` VARCHAR(50) NOT NULL,            /* url, geo*/
  `url_type` VARCHAR(50) DEFAULT '',     /* landing_page_url */
  `match_type` VARCHAR(50) DEFAULT '',   /* starts_with, ends_with, contains */
  `match_string` VARCHAR(50) DEFAULT '', /* "baidu" */
  `match_rules` TEXT,         /* []model.MatchRule{{"province", 110000}, }, */
  `source_rules` TEXT,        /* map[string][]string{"sdk": {"app1"},},*/
  `targets` TEXT,             /* target_ids*/
  `inverted` BOOL DEFAULT FALSE,         /* true, false */
  KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `allocation_rule` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `rule_type` VARCHAR(50) NOT NULL,  /*"roundrobin", "random", */
  UNIQUE KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;