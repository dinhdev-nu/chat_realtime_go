-- name: GetLastMessagesByRoomId :many
SELECT d.message_id FROM go_db_chat_messages_direct d 
JOIN go_db_chat_message_status s ON d.message_id = s.message_id
WHERE d.message_room_id = ?;

-- name: GetMessagesDirectByRoomId :many
SELECT d.message_id, 
       d.message_room_id, 
       d.message_receiver_id, 
       d.message_content, 
       d.message_type,
       d.message_sent_at,
       s.message_is_read,
       s.message_read_at
FROM go_db_chat_messages_direct d 
JOIN go_db_chat_message_status s ON d.message_id = s.message_id
WHERE d.message_room_id = ?
ORDER BY d.message_sent_at DESC
LIMIT ? OFFSET ?;

-- name: SaveMessageDirect :execresult
INSERT INTO go_db_chat_messages_direct (
    message_room_id, 
    message_receiver_id, 
    message_content, 
    message_type, 
    message_sent_at
) VALUES (?, ?, ?, ?, ?);