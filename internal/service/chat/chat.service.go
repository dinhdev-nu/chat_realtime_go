package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
)

type IChatService interface {
	InitChat(userInfo *model.GoDbUserInfo) (dto.InitChatOutPut, error)
	// messages
	HandleSendMesage(msg dto.OnMessage) (dto.SaveMessageDTO, error)
	GetMessagesFromRoom(roomId string, page string) (dto.GetMessagesFromRoomOutput, error)
	// rooms
	GetRoomChatById(id uint64) (map[string]interface{}, error)
	CreateRoomChat(data *dto.CreateRoomDTO) (dto.CreateRoomDTO, error)
	UpdateStatusMessages(data *dto.UpdateStatusInput) error
}

type chatService struct {
	Arepo repo.IAuthRepo
	Crepo repo.IChatRepo
	Urepo repo.IUserRepo
}

func NewChatService(chatRepo repo.IChatRepo, authRepo repo.IAuthRepo, userRepo repo.IUserRepo) IChatService {
	return &chatService{
		Arepo: authRepo,
		Crepo: chatRepo,
		Urepo: userRepo,
	}
}

func (s *chatService) HandleSendMesage(msg dto.OnMessage) (dto.SaveMessageDTO, error) {
	// check valid message

	var data dto.SaveMessageDTO
	data.MessageRoomID = uint64(msg.Message.RoomID)
	data.MessageContent = msg.Message.Content
	data.MessageSentAt = msg.Message.SendAt

	// checck valid type
	switch msg.Type {
	case "single":

		// Inset message to database
		data.MessageReceiverID = uint64(msg.ReceiverID)

		msgId, err := s.Crepo.SaveMessegeDirect(data)

		if err != nil {
			return dto.SaveMessageDTO{}, err
		}

		// Insert message type
		err = s.Crepo.SaveMessageStatus(uint64(msgId), msg.ReceiverID)

		if err != nil {
			return dto.SaveMessageDTO{}, err
		}
		data.MessageID = msgId

		fmt.Println("HandleSendMesage:10 ", msg)

	case "multi":
		// Insert message to database
		data.MessageSenderID = uint64(msg.SendID)
		msgId, err := s.Crepo.SaveMessegeGroup(data)
		if err != nil {
			return dto.SaveMessageDTO{}, err
		}
		data.MessageID = msgId
	}

	return data, nil
}

// init chat
func (s *chatService) InitChat(userInfo *model.GoDbUserInfo) (dto.InitChatOutPut, error) {
	var res dto.InitChatOutPut
	res.CurrentUser = userInfo

	userId := userInfo.UserID

	// Get user rooms
	data, err := s.Crepo.GetRoomsByUserId(userId)
	if err != nil {
		return dto.InitChatOutPut{}, err
	}
	// Get last message
	for _, room := range data {
		if room.MessageReceiverID == uint64(userId) {
			// Get other user member info
			anotherUser, err := s.Crepo.GetAnotherUserID(room.RoomID, userId)
			if err != nil {
				return dto.InitChatOutPut{}, err
			}
			// Get user info
			user, err := s.Arepo.GetUserInfoByID(anotherUser)
			if err != nil {
				return dto.InitChatOutPut{}, err
			}
			res.Rooms = append(res.Rooms, dto.RoomInitChat{
				RoomInfo: room,
				Users: dto.InfoUserPrivateChat{
					UserID:     user.UserID,
					UserName:   user.UserNickname,
					UserAvatar: user.UserAvatar,
				},
			})
		} else {
			// Get user info
			user, err := s.Arepo.GetUserInfoByID(int64(room.MessageReceiverID))
			if err != nil {
				return dto.InitChatOutPut{}, err
			}
			res.Rooms = append(res.Rooms, dto.RoomInitChat{
				RoomInfo: room,
				Users: dto.InfoUserPrivateChat{
					UserID:     user.UserID,
					UserName:   user.UserNickname,
					UserAvatar: user.UserAvatar,
				},
			})
		}
	}
	// get user status
	for i, room := range res.Rooms {
		if !room.RoomInfo.RoomIsGroup {
			status, _ := s.Urepo.GetStatusByUserId(room.Users.UserID)
			res.Rooms[i].Users.UserStatus = status

			res.Followers = append(res.Followers, room.Users.UserID)
		} else {
			res.Rooms[i].Users.UserStatus = "offline"
		}
	}

	res.SocketUrl = fmt.Sprintf("ws://%s:%s/v1/api/chat/ws?user_id=%d&token=",
		global.Config.Server.Host, global.Config.Server.Port, userId)
	return res, nil
}

// messages
func (s *chatService) GetMessagesFromRoom(roomId string, page string) (dto.GetMessagesFromRoomOutput, error) {
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	roomIdUint := utils.StringToUint64(roomId)

	room, err := s.Crepo.GetRoomById(roomIdUint)
	if err != nil {
		return dto.GetMessagesFromRoomOutput{}, err
	}
	var result dto.GetMessagesFromRoomOutput
	result.RoomIsGroup = room.RoomIsGroup
	if !room.RoomIsGroup {
		messegers, err := s.Crepo.GetMessagesFromRoom(roomIdUint, utils.StringToInt64(page))
		if err != nil {
			return dto.GetMessagesFromRoomOutput{}, err
		}
		result.MessagesDriect = messegers
	} else {
		messegers, err := s.Crepo.GetMessagesGruopFromRoom(roomIdUint, utils.StringToInt64(page))
		if err != nil {
			return dto.GetMessagesFromRoomOutput{}, err
		}
		result.MessagesGroup = messegers
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

func (s *chatService) CreateRoomChat(data *dto.CreateRoomDTO) (dto.CreateRoomDTO, error) {
	// Check if room already exists
	room, err := s.Crepo.GetRoomByName(data.RoomName)
	if !errors.Is(err, sql.ErrNoRows) {
		return dto.CreateRoomDTO{}, err
	}

	if room.RoomID > 0 {
		return dto.CreateRoomDTO{
			RoomID:      room.RoomID,
			RoomName:    data.RoomName,
			RoomIsGroup: data.RoomIsGroup,
			RoomMembers: data.RoomMembers,
		}, nil
	}

	err = s.Crepo.CreateRoom(data)
	if err != nil {
		return dto.CreateRoomDTO{}, err
	}
	// Create room members
	err = s.Crepo.AddMembersToRoom(data.RoomID, data.RoomMembers)
	if err != nil {
		return dto.CreateRoomDTO{}, err
	}

	return *data, nil
}

func (s *chatService) UpdateStatusMessages(data *dto.UpdateStatusInput) error {
	err := s.Crepo.UpdateMessageStatus(data)
	if err != nil {
		return err
	}
	return nil
}
