-- +goose Up
-- +goose StatementBegin
CREATE TABLE `go_db_chat_rooms` (
  `room_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Room ID',
  `room_name` VARCHAR(100) DEFAULT NULL COMMENT 'Room name (null nếu 1-1)',
  `room_is_group` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Is group chat?',
  `room_created_by` BIGINT UNSIGNED DEFAULT NULL COMMENT 'Created by user ID',
  `room_created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
  PRIMARY KEY (`room_id`),
  FOREIGN KEY (`room_created_by`) REFERENCES `go_db_user_info`(`user_id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Chat rooms';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `go_db_chat_rooms`;
-- +goose StatementEnd


