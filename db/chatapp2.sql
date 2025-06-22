-- SET TIME ZONE cho PostgreSQL (nếu dùng)
-- SET TIME ZONE 'UTC';

-- DROP bảng nếu tồn tại (xóa từ bảng phụ -> bảng chính)
DROP TABLE IF EXISTS 
  message_reads,
  messages,
  conversation_members,
  conversations,
  friend_requests,
  friends,
  users
CASCADE;

-- USERS
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE,
  password_hash TEXT,
  avatar_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CONVERSATIONS
CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  is_group BOOLEAN DEFAULT FALSE,
  name VARCHAR(100),
  created_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CONVERSATION MEMBERS
CREATE TABLE conversation_members (
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_convo_members_user ON conversation_members (user_id);

-- MESSAGES
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
  text TEXT,
  attachment_url TEXT,
  type VARCHAR(20) DEFAULT 'text', -- 'text', 'image', 'file'
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_messages_convo_created ON messages (conversation_id, created_at DESC);

-- MESSAGE READS
CREATE TABLE message_reads (
  message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (message_id, user_id)
);

CREATE INDEX idx_message_reads_user ON message_reads (user_id);

-- FRIEND REQUESTS
CREATE TABLE friend_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  from_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  to_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'declined'
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (from_user_id, to_user_id)
);

-- FRIENDS
CREATE TABLE friends (
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  friend_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, friend_id)
);

CREATE INDEX idx_friends_user ON friends (user_id, friend_id);

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  type VARCHAR(50),
  message TEXT,
  metadata JSONB,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);
