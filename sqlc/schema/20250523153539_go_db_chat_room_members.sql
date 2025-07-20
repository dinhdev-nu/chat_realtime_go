-- +goose Up
-- +goose StatementBegin
CREATE TABLE `go_db_chat_room_members` (
  `room_id` bigint unsigned NOT NULL COMMENT 'Room ID',
  `member_user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
  `member_nickname` VARCHAR(25) DEFAULT NULL COMMENT 'Nickname in the room',
  `member_joined_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Joined at',
  `member_role` VARCHAR(15) NOT NULL DEFAULT 'member' COMMENT 'Role in the room',
  `member_last_seen` BIGINT UNSIGNED DEFAULT NULL COMMENT "Message ID last Seen",
  FOREIGN KEY (`room_id`) REFERENCES `go_db_chat_rooms`(`room_id`) ON DELETE CASCADE,
  FOREIGN KEY (`member_user_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Room members';
-- +goose StatementEnd

-- +goose StatementBegin  
CREATE INDEX `idx_room_members_user_id` ON `go_db_chat_room_members`(`member_user_id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `go_db_chat_room_members`;
-- +goose StatementEnd
