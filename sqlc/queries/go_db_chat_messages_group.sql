-- name: GetMessagesGroupByRoomId :many
SELECT * FROM go_db_chat_messages_group   
WHERE message_room_id = ?
ORDER BY message_send_at DESC
LIMIT ? OFFSET ?;

-- name: SaveMessageGroup :execresult
INSERT INTO go_db_chat_messages_group (
    message_room_id, 
    message_sender_id, 
    message_content, 
    message_type, 
    message_sent_at
) VALUES (?, ?, ?, ?, ?);
