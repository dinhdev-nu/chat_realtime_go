-- name: GetRoomById :one
SELECT * FROM `go_db_chat_rooms` 
WHERE `room_id` = ? LIMIT 1;

-- name: GetRoomByName :one
SELECT `room_id`, `room_is_group` 
FROM `go_db_chat_rooms` 
WHERE `room_name` IN (?, ?) 
AND `room_is_group` = 0
LIMIT 1;

-- name: GetRoomGroupByName :one
SELECT `room_id`, `room_is_group`
FROM `go_db_chat_rooms` 
WHERE `room_name` = ? AND `room_is_group` = 1
LIMIT 1;

-- name: CreateRoom :execresult
INSERT INTO `go_db_chat_rooms` (
    `room_name` , `room_description`, `room_avatar`, 
    `room_is_group`, `room_created_by`
) VALUES (?, ?, ?, ? , ?);

-- name: GetGroupRoomsByUserId :many
SELECT 
  r.room_id, 
  r.room_name, 
  r.room_description, 
  r.room_avatar,
  r.room_created_by,
  r.room_created_at,
  mm.member_user_id, 
  mm.member_nickname, 
  mm.member_role,
  mm.member_last_seen,
  g.message_id, 
  g.message_content, 
  g.message_sender_id,
  g.message_type, 
  g.message_sent_at
FROM go_db_chat_rooms r
JOIN go_db_chat_room_members m 
  ON r.room_id = m.room_id AND m.member_user_id = ?
LEFT JOIN (
  SELECT gm.*
  FROM go_db_chat_messages_group gm
  INNER JOIN (
    SELECT message_room_id, MAX(message_sent_at) AS max_sent_at
    FROM go_db_chat_messages_group
    GROUP BY message_room_id
  ) latest
    ON gm.message_room_id = latest.message_room_id
    AND gm.message_sent_at = latest.max_sent_at
) g
  ON g.message_room_id = r.room_id
RIGHT jOIN go_db_chat_room_members mm 
  ON mm.room_id = r.room_id
WHERE r.room_is_group = 1;


-- name: GetPrivateRoomsByUserId :many
SELECT 
  r.room_id, 
  
  m.member_user_id as current_user_id,
  m.member_last_seen as current_user_last_seen,

  mm.member_user_id,
  mm.member_last_seen,
  g.message_id, 
  g.message_content, 
  g.message_receiver_id,
  g.message_type, 
  g.message_sent_at
FROM go_db_chat_rooms r 
JOIN go_db_chat_room_members m
  ON r.room_id = m.room_id AND m.member_user_id = ? AND r.room_is_group = 0
LEFT JOIN (
  SELECT dm.*
  FROM go_db_chat_messages_direct dm
  INNER JOIN (
    SELECT message_room_id, MAX(message_sent_at) AS max_sent_at
    FROM go_db_chat_messages_direct
    GROUP BY message_room_id
  ) latest
    ON dm.message_room_id = latest.message_room_id
    AND dm.message_sent_at = latest.max_sent_at
) g 
  ON g.message_room_id = r.room_id
  LEFT JOIN go_db_chat_room_members mm
  ON mm.room_id = r.room_id AND mm.member_user_id != m.member_user_id