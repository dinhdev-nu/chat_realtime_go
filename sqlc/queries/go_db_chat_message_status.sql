-- name: SaveMessageStatus :exec
INSERT INTO go_db_chat_message_status (
    message_id,
    message_user_id
) VALUES (?, ?);

-- name: UpdateMessageStatus :exec
UPDATE go_db_chat_message_status s
JOIN go_db_chat_messages_direct d ON s.message_id = d.message_id
JOIN go_db_chat_rooms r ON d.message_room_id = r.room_id
SET s.message_is_read = TRUE,
    s.message_read_at = NOW()
WHERE r.room_id = ? AND d.message_receiver_id = ? AND s.message_is_read = FALSE;
