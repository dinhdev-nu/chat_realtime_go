-- +goose Up
-- +goose StatementBegin
CREATE TABLE `go_db_chat_message_status` (
  `message_id` BIGINT UNSIGNED NOT NULL COMMENT 'Message ID (from direct)',
  `message_user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID - usually receiver',
  `message_is_read` BOOLEAN DEFAULT FALSE COMMENT 'Read status',
  `message_read_at` timestamp NULL DEFAULT NULL COMMENT 'Read timestamp',
  PRIMARY KEY (`message_id`, `message_user_id`),
  FOREIGN KEY (`message_id`) REFERENCES `go_db_chat_messages_direct`(`message_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_user_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Read status for direct messages';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX `idx_msg_status_user_id` ON `go_db_chat_message_status`(`message_user_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `go_db_chat_message_status`;
-- +goose StatementEnd

