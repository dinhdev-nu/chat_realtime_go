-- +goose Up
-- +goose StatementBegin
CREATE TABLE `go_db_chat_messages_direct` (
  `message_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `message_room_id` bigint unsigned NOT NULL COMMENT 'Room ID',
  `message_receiver_id` BIGINT UNSIGNED NOT NULL COMMENT 'Receiver ID',
  `message_content` TEXT NOT NULL,
  `message_type` ENUM('text', 'image', 'file') DEFAULT 'text',
  `message_sent_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Sent at',
  FOREIGN KEY (`message_room_id`) REFERENCES `go_db_chat_rooms`(`room_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_receiver_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='1-1 direct messages';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX `idx_direct_sent_at` ON `go_db_chat_messages_direct`(`message_sent_at`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `go_db_chat_messages_direct`;
-- +goose StatementEnd
