CREATE TABLE `message_beep` (
  `agent_id`  CHAR(20) NOT NULL,
  `client_type` CHAR(20) NOT NULL, /* web, sdk, ...*/
  `beep_type` VARCHAR(20) NOT NULL, /* desktop or popup */
  `new_conversation` BOOLEAN NOT NULL DEFAULT FALSE,
  `new_message` BOOLEAN NOT NULL DEFAULT FALSE,
  `conversation_transfer_in` BOOLEAN NOT NULL DEFAULT FALSE,
  `conversation_transfer_out` BOOLEAN NOT NULL DEFAULT FALSE,
  `colleague_conversation` BOOLEAN NOT NULL DEFAULT FALSE,
  UNIQUE KEY (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
