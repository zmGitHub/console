/*
{
  "created_on": "2016-07-28T02:05:28.298909",
  "enterprise_id": 5869,
  "form_def": {
    "captcha": false,
    "inputs": {
      "fields": [
        {
          "display_name": "\u59d3\u540d",
          "field_name": "name",
          "ignore_returned_customer": false,
          "optional": true,
          "type": "text"
        },
        {
          "display_name": "1.\u5145\u503c\u95ee\u9898",
          "field_name": "1.\u5145\u503c\u95ee\u9898",
          "ignore_returned_customer": false,
          "optional": true,
          "type": "text"
        },
        {
          "display_name": "bbbbbbbbbbb",
          "field_name": "bbbbbbbbbbb",
          "ignore_returned_customer": false,
          "optional": true,
          "type": "text"
        },
        {
          "display_name": "\u7535\u8bdd",
          "field_name": "tel",
          "ignore_returned_customer": false,
          "optional": true,
          "type": "text"
        }
      ],
      "imclient": "open",
      "title": "\u60a8\u60f3\u54a8\u8be2\u54ea\u65b9\u9762"
    },
    "menus": {
      "assignments": [
        {
          "description": "ggggggg",
          "target": null,
          "target_kind": null
        },
        {
          "description": "\u81f3\u5c0a\u5361",
          "target": null,
          "target_kind": null
        },
        {
          "description": "\u552e\u524d",
          "target": null,
          "target_kind": null
        },
        {
          "description": "\u9648\u8001\u677f",
          "target": null,
          "target_kind": null
        }
      ],
      "imclient": "close",
      "title": "<\u6d4b\u8bd5>\u60a8\u662f\u4f55\u79cd\u5361\u9879"
    },
    "version": 1
  },
  "id": 22,
  "last_updated": "2018-11-09T00:25:22.519722",
  "title": "\u8be2\u524d\u8868\u5355"
}

*/

CREATE TABLE `prechat_form` (
  `id` CHAR(20) PRIMARY KEY,
  `ent_id` CHAR(20) NOT NULL,
  `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '启用/禁用',
  `title` VARCHAR(200) NOT NULL,
  `form_fields` TEXT COLLATE utf8mb4_unicode_ci,
  `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
  `updated_at` DATETIME(6) NOT NULL DEFAULT NOW(6) ON UPDATE NOW(6),
  UNIQUE KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;

