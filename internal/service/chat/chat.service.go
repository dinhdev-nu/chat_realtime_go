package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/input"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
)

type IChatService interface {
	InitChat(userInfo *model.GoDbUserInfo) (map[string]interface{}, error)
	// messages
	HandleSendMesage(msg input.Message) (map[string]interface{}, error)
	GetMessagesFromRoom(roomId string, page string) (map[string]interface{}, error)
	// rooms
	GetRoomChatById(id uint64) (map[string]interface{}, error)
	CreateRoomChat(data *input.CreateRoomInput) (map[string]interface{}, error)
	UpdateStatusMessages(data *input.UpdateStatusInput) error
}

type chatService struct {
	Arepo repo.IAuthRepo
	Crepo repo.IChatRepo
}

func NewChatService(chatRepo repo.IChatRepo, authRepo repo.IAuthRepo) IChatService {
	return &chatService{
		Arepo: authRepo,
		Crepo: chatRepo,
	}
}

func (s *chatService) HandleSendMesage(msg input.Message) (map[string]interface{}, error) {
	// check valid message

	data := make(map[string]interface{})
	data["message_room_id"] = msg.RoomID
	data["message_content"] = msg.Content
	data["message_sent_at"] = msg.SendAt

	// checck valid type
	switch msg.Type {
	case "single":

		// Inset message to database
		data["message_receiver_id"] = msg.ReceiverID

		msgId, err := s.Crepo.SaveMessegeDirect(data)

		if err != nil {
			return nil, err
		}

		// Insert message type
		err = s.Crepo.SaveMessageStatus(uint64(msgId), msg.ReceiverID)

		if err != nil {
			return nil, err
		}
		data["message_id"] = msgId

		fmt.Println("HandleSendMesage:10 ", msg)

	case "multi":
		// Insert message to database
		data["message_sender_id"] = msg.SendID
		msgId, err := s.Crepo.SaveMessegeGroup(data)
		if err != nil {
			return nil, err
		}
		data["message_id"] = msgId
	}

	return data, nil
}

// init chat
func (s *chatService) InitChat(userInfo *model.GoDbUserInfo) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	res["user"] = userInfo

	userId := userInfo.UserID

	// Get user rooms
	rooms := make([]map[string]interface{}, 0)
	roomsId, err := s.Crepo.GetRoomsByUserId(userId)
	if err != nil {
		return nil, err
	}
	// Get last message
	for _, room := range roomsId {
		if room.MessageReceiverID == uint64(userId) {
			// Get other user member info
			anotherUser, err := s.Crepo.GetAnotherUserID(room.RoomID, userId)
			if err != nil {
				return nil, err
			}
			// Get user info
			user, err := s.Arepo.GetUserInfoByID(anotherUser)
			if err != nil {
				return nil, err
			}
			rooms = append(rooms, map[string]interface{}{
				"room": room,
				"info": map[string]interface{}{
					"user_id":     user.UserID,
					"user_name":   user.UserNickname,
					"user_avatar": user.UserAvatar,
				},
			})
		} else {
			// Get user info
			user, err := s.Arepo.GetUserInfoByID(int64(room.MessageReceiverID))
			if err != nil {
				return nil, err
			}
			rooms = append(rooms, map[string]interface{}{
				"room": room,
				"info": map[string]interface{}{
					"user_id":     user.UserID,
					"user_name":   user.UserNickname,
					"user_avatar": user.UserAvatar,
				},
			})
		}
	}
	res["rooms"] = rooms
	res["socket_url"] = fmt.Sprintf("ws://%s:%s/v1/api/chat/ws?user_id=%d&token=",
		global.Config.Server.Host, global.Config.Server.Port, userId)
	return res, nil
}

// messages
func (s *chatService) GetMessagesFromRoom(roomId string, page string) (map[string]interface{}, error) {
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	roomIdUint := utils.StringToUint64(roomId)

	room, err := s.Crepo.GetRoomById(roomIdUint)
	if err != nil {
		return nil, err
	}
	var result = map[string]interface{}{
		"room_id": room.RoomID,
	}
	if !room.RoomIsGroup {
		messegers, err := s.Crepo.GetMessagesFromRoom(roomIdUint, utils.StringToInt64(page))
		if err != nil {
			return nil, err
		}
		result["messages"] = messegers
	} else {
		messegers, err := s.Crepo.GetMessagesGruopFromRoom(roomIdUint, utils.StringToInt64(page))
		if err != nil {
			return nil, err
		}
		result["messages"] = messegers
	}

	return result, nil // Placeholder for message retrieval logic
}

// rooms
func (s *chatService) GetRoomChatById(id uint64) (map[string]interface{}, error) {
	room, err := s.Crepo.GetRoomById(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"room_id":   room.RoomID,
		"room_name": room.RoomName,
	}, nil
}

func (s *chatService) CreateRoomChat(data *input.CreateRoomInput) (map[string]interface{}, error) {
	// Check if room already exists
	room, err := s.Crepo.GetRoomByName(data.RoomName)
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if room.RoomID > 0 {
		return map[string]interface{}{
			"room_id":  room.RoomID,
			"is_group": data.RoomIsGroup,
		}, nil
	}

	// Prepare payload for creating a new room
	payload := map[string]interface{}{
		"room_name":      data.RoomName,
		"room_create_by": data.RoomCreateBy,
		"room_is_group":  data.RoomIsGroup,
	}

	var id int64
	id, err = s.Crepo.CreateRoom(payload)
	if err != nil {
		return nil, err
	}
	// Create room members
	err = s.Crepo.AddMembersToRoom(id, data.RoomMembers)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"room_id":  id,
		"is_group": data.RoomIsGroup,
	}, nil
}

func (s *chatService) UpdateStatusMessages(data *input.UpdateStatusInput) error {
	err := s.Crepo.UpdateMessageStatus(
		data.RoomID,
		data.UserId,
	)
	if err != nil {
		return err
	}
	return nil
}
