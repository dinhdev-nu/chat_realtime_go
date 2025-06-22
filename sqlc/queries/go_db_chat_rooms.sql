-- name: GetRoomById :one
SELECT * FROM `go_db_chat_rooms` 
WHERE `room_id` = ? LIMIT 1;

-- name: GetRoomByName :one
SELECT `room_id`, `room_is_group` 
FROM `go_db_chat_rooms` 
WHERE `room_name` IN (?, ?) 
LIMIT 1;

-- name: CreateRoom :execresult
INSERT INTO `go_db_chat_rooms` (`room_name`, `room_is_group`, `room_created_by`) VALUES (?, ?, ?);


