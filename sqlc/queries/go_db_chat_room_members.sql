-- name: GetAnotherPrivateMenberByRoomId :one
SELECT member_user_id, member_last_seen FROM go_db_chat_room_members 
WHERE room_id = ? AND member_user_id != ?;

-- name: InsetMemberToRoom :exec
INSERT INTO go_db_chat_room_members (
  room_id, 
  member_user_id,
  member_role
) 
VALUES (?, ?, ?);

-- name: UpdateMemberLastSeen :exec
UPDATE go_db_chat_room_members
SET member_last_seen = ?
WHERE room_id = ? AND member_user_id = ?;