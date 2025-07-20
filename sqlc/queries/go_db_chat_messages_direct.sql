-- name: GetMessagesDirectByRoomId :many
SELECT message_id, 
       message_room_id, 
       message_receiver_id, 
       message_content, 
       message_type,
       message_sent_at
FROM go_db_chat_messages_direct
WHERE message_room_id = ?
ORDER BY message_sent_at DESC
LIMIT ? OFFSET ?;

-- name: SaveMessageDirect :execresult
INSERT INTO go_db_chat_messages_direct (
    message_room_id, 
    message_receiver_id, 
    message_content, 
    message_type, 
    message_sent_at
) VALUES (?, ?, ?, ?, ?);