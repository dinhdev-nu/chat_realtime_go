package dto

type SearchUsersOutput struct {
	UserID       int64  `json:"user_id"`
	UserNickname string `json:"user_nickname"`
	UserAvatar   string `json:"user_avatar"`
	UserGender   int32  `json:"user_gender"`
	UserStatus   string `json:"user_status"`
}
