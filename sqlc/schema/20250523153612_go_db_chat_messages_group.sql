-- +goose Up
-- +goose StatementBegin
CREATE TABLE `go_db_chat_messages_group` (
  `message_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `message_room_id` bigint unsigned NOT NULL COMMENT 'Room ID',
  `message_sender_id` BIGINT UNSIGNED NOT NULL COMMENT 'Sender user ID',
  `message_content` TEXT NOT NULL,
  `message_type` ENUM('text', 'image', 'file') DEFAULT 'text',
  `message_sent_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Sent at',
  FOREIGN KEY (`message_room_id`) REFERENCES `go_db_chat_rooms`(`room_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_sender_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Group messages';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX `idx_group_sender_id` ON `go_db_chat_messages_group`(`message_sender_id`);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX `idx_group_sent_at` ON `go_db_chat_messages_group`(`message_sent_at`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `go_db_chat_messages_group`;
-- +goose StatementEnd
