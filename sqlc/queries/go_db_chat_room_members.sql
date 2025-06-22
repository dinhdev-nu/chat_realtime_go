-- name: GetPrivateRoomsByUserId :many
SELECT 
  m.room_id, 
  r.room_is_group, 
  d.message_receiver_id, 
  d.message_content, 
  d.message_type,
  d.message_id,
  d.message_sent_at,
  s.message_is_read
FROM go_db_chat_room_members m
JOIN go_db_chat_rooms r 
  ON m.room_id = r.room_id
LEFT JOIN (
    SELECT * FROM go_db_chat_messages_direct d1
    WHERE (d1.message_sent_at) IN (
        SELECT MAX(message_sent_at)
        FROM go_db_chat_messages_direct
        GROUP BY message_room_id
    ) 
) d ON m.room_id = d.message_room_id
LEFT JOIN go_db_chat_message_status s ON d.message_id = s.message_id
WHERE 
  m.member_user_id = ?
  AND r.room_is_group = 0;


-- name: GetAnotherPrivateMenberByRoomId :one
SELECT member_user_id FROM go_db_chat_room_members 
WHERE room_id = ? AND member_user_id != ?;

-- name: InsetMemberToRoom :exec
INSERT INTO go_db_chat_room_members (room_id, member_user_id) 
VALUES (?, ?);

