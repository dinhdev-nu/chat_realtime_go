-- DROP TABLES (đảm bảo đúng thứ tự)
DROP TABLE IF EXISTS `go_db_chat_message_status`;
DROP TABLE IF EXISTS `go_db_chat_messages_direct`;
DROP TABLE IF EXISTS `go_db_chat_messages_group`;
DROP TABLE IF EXISTS `go_db_chat_room_members`;
DROP TABLE IF EXISTS `go_db_chat_rooms`;

-- ROOMS
CREATE TABLE `go_db_chat_rooms` (
  `room_id` CHAR(36) NOT NULL COMMENT 'Room ID',
  `room_name` VARCHAR(100) DEFAULT NULL COMMENT 'Room name (null nếu 1-1)',
  `room_is_group` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Is group chat?',
  `room_created_by` BIGINT UNSIGNED DEFAULT NULL COMMENT 'Created by user ID',
  `room_created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
  PRIMARY KEY (`room_id`),
  FOREIGN KEY (`room_created_by`) REFERENCES `go_db_user_info`(`user_id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Chat rooms';


-- ROOM MEMBERS
CREATE TABLE `go_db_chat_room_members` (
  `room_id` CHAR(36) NOT NULL COMMENT 'Room ID',
  `member_user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID',
  `member_joined_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'Joined at',
  PRIMARY KEY (`room_id`, `member_user_id`),
  FOREIGN KEY (`room_id`) REFERENCES `go_db_chat_rooms`(`room_id`) ON DELETE CASCADE,
  FOREIGN KEY (`member_user_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Room members';


-- MESSAGES - GROUP
CREATE TABLE `go_db_chat_messages_group` (
  `message_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `message_room_id` CHAR(36) NOT NULL COMMENT 'Room ID',
  `message_sender_id` BIGINT UNSIGNED NOT NULL COMMENT 'Sender user ID',
  `message_content` TEXT NOT NULL,
  `message_type` ENUM('text', 'image', 'file') DEFAULT 'text',
  `message_sent_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`message_room_id`) REFERENCES `go_db_chat_rooms`(`room_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_sender_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Group messages';


-- MESSAGES - DIRECT
CREATE TABLE `go_db_chat_messages_direct` (
  `message_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `message_sender_id` BIGINT UNSIGNED NOT NULL COMMENT 'Sender ID',
  `message_receiver_id` BIGINT UNSIGNED NOT NULL COMMENT 'Receiver ID',
  `message_content` TEXT NOT NULL,
  `message_type` ENUM('text', 'image', 'file') DEFAULT 'text',
  `message_sent_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`message_sender_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_receiver_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='1-1 direct messages';


-- MESSAGE STATUS (áp dụng cho 1-1)
CREATE TABLE `go_db_chat_message_status` (
  `message_id` BIGINT UNSIGNED NOT NULL COMMENT 'Message ID (from direct)',
  `message_user_id` BIGINT UNSIGNED NOT NULL COMMENT 'User ID - usually receiver',
  `message_is_read` BOOLEAN DEFAULT FALSE COMMENT 'Read status',
  `message_read_at` DATETIME DEFAULT NULL COMMENT 'Read time',
  PRIMARY KEY (`message_id`, `message_user_id`),
  FOREIGN KEY (`message_id`) REFERENCES `go_db_chat_messages_direct`(`message_id`) ON DELETE CASCADE,
  FOREIGN KEY (`message_user_id`) REFERENCES `go_db_user_info`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Read status for direct messages';



-- INDEXES
-- Room members
CREATE INDEX `idx_room_members_user_id` ON `go_db_chat_room_members`(`member_user_id`);

-- Group messages
CREATE INDEX `idx_group_room_id` ON `go_db_chat_messages_group`(`message_room_id`);
CREATE INDEX `idx_group_sender_id` ON `go_db_chat_messages_group`(`message_sender_id`);
CREATE INDEX `idx_group_sent_at` ON `go_db_chat_messages_group`(`message_sent_at`);

-- Direct messages
CREATE INDEX `idx_direct_sender_id` ON `go_db_chat_messages_direct`(`message_sender_id`);
CREATE INDEX `idx_direct_receiver_id` ON `go_db_chat_messages_direct`(`message_receiver_id`);
CREATE INDEX `idx_direct_sent_at` ON `go_db_chat_messages_direct`(`message_sent_at`);

-- Message status
CREATE INDEX `idx_msg_status_user_id` ON `go_db_chat_message_status`(`message_user_id`);
