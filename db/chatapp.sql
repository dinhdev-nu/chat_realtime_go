-- DROP TABLES (đảm bảo đúng thứ tự)
DROP TABLE IF EXISTS message_status;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS room_members;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

-- USERS
CREATE TABLE users (
  id CHAR(36) PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  full_name VARCHAR(100),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ROOMS
CREATE TABLE rooms (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(100),
  is_group BOOLEAN NOT NULL DEFAULT TRUE,
  created_by CHAR(36),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (created_by) REFERENCES users(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE
) ENGINE=InnoDB;

-- ROOM MEMBERS
CREATE TABLE room_members (
  room_id CHAR(36) NOT NULL,
  user_id CHAR(36) NOT NULL,
  joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (room_id, user_id),
  FOREIGN KEY (room_id) REFERENCES rooms(id)
    ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;

-- MESSAGES
CREATE TABLE messages (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  sender_id CHAR(36) NOT NULL,
  room_id CHAR(36),
  receiver_id CHAR(36),
  content TEXT NOT NULL,
  type ENUM('text', 'image', 'file') DEFAULT 'text',
  sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (sender_id) REFERENCES users(id)
    ON DELETE CASCADE,
  FOREIGN KEY (room_id) REFERENCES rooms(id)
    ON DELETE CASCADE,
  FOREIGN KEY (receiver_id) REFERENCES users(id)
    ON DELETE CASCADE,
  CHECK (
    (room_id IS NOT NULL AND receiver_id IS NULL) OR
    (room_id IS NULL AND receiver_id IS NOT NULL)
  )
) ENGINE=InnoDB;

-- MESSAGE STATUS (optional)
CREATE TABLE message_status (
  message_id BIGINT NOT NULL,
  user_id CHAR(36) NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  read_at DATETIME,
  PRIMARY KEY (message_id, user_id),
  FOREIGN KEY (message_id) REFERENCES messages(id)
    ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;


-- Index cho tìm kiếm tin nhắn theo room hoặc user
CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_sent_at ON messages(sent_at);

-- Index cho truy vấn thành viên room
CREATE INDEX idx_room_members_user_id ON room_members(user_id);

-- Index cho message_status để lấy trạng thái đã đọc
CREATE INDEX idx_message_status_user_id ON message_status(user_id);
