package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	GetMessagesFromRoom(roomId string, page string, offset string) (dto.GetMessagesFromRoomOutput, error)
	// rooms
	GetRoomsByUserID(userID int64) ([]dto.GetGroupRoomResponse, error)
	GetGroupRoomsByUserID2(userInfo *model.GoDbUserInfo) ([]dto.GetGroupRoomResponse, error)

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
		data.MessageID = msgId

	case "group":
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
		// Get user info
		user, err := s.Arepo.GetUserInfoByID(room.MemberUserID.Int64)
		if err != nil {
			return dto.InitChatOutPut{}, err
		}
		// get status member
		status, _ := s.Urepo.GetStatusByUserId(room.MemberUserID.Int64)

		res.Rooms = append(res.Rooms, dto.RoomInitChat{
			RoomInfo: room,
			Users: dto.InfoUserPrivateChat{
				UserID:     user.UserID,
				UserName:   user.UserNickname,
				UserAvatar: user.UserAvatar,
				UserStatus: status,
			},
		})
		res.Followers = append(res.Followers, room.MemberUserID.Int64)
	}

	res.SocketUrl = fmt.Sprintf("ws://%s:%s/v1/api/chat/ws?user_id=%d&token=",
		global.Config.Server.Host, global.Config.Server.Port, userId)
	return res, nil
}

// messages
func (s *chatService) GetMessagesFromRoom(roomId string, page string, offset string) (dto.GetMessagesFromRoomOutput, error) {
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	if offset == "" {
		offset = "0" // Default to offset 0 if not provided
	}
	roomIdUint := utils.StringToUint64(roomId)

	room, err := s.Crepo.GetRoomById(roomIdUint)
	if err != nil {
		return dto.GetMessagesFromRoomOutput{}, err
	}
	var result dto.GetMessagesFromRoomOutput
	result.RoomIsGroup = room.RoomIsGroup
	if !room.RoomIsGroup {
		messegers, err := s.Crepo.GetMessagesFromRoom(roomIdUint, utils.StringToInt64(page), utils.StringToInt64(offset))
		if err != nil {
			return dto.GetMessagesFromRoomOutput{}, err
		}
		result.MessagesDriect = messegers
	} else {
		messegers, err := s.Crepo.GetMessagesGruopFromRoom(roomIdUint, utils.StringToInt64(page), utils.StringToInt64(offset))
		if err != nil {
			return dto.GetMessagesFromRoomOutput{}, err
		}
		result.MessagesGroup = messegers
	}

	return result, nil // Placeholder for message retrieval logic
}

// rooms
func (s *chatService) GetRoomsByUserID(userID int64) ([]dto.GetGroupRoomResponse, error) {
	rooms, err := s.Crepo.GetGroupRoomsByUserId(userID)
	if err != nil {
		fmt.Println("GetRoomsByUserID: error fetching group rooms:", err)
		return nil, err
	}
	if len(rooms) == 0 {
		return []dto.GetGroupRoomResponse{}, nil
	}
	now := time.Now()
	// Get map rooms -> members and get 1 array memberids
	roomMembersMap := make(map[uint64][]uint64)
	menberIDs := make([]uint64, 0)
	for _, r := range rooms {
		roomMembersMap[r.RoomID] = append(roomMembersMap[r.RoomID], r.MemberUserID)
		if !utils.Contains(menberIDs, r.MemberUserID) || r.MemberUserID != uint64(userID) {
			menberIDs = append(menberIDs, r.MemberUserID)
		}
	}

	// get user info by member ids
	users, err := s.Urepo.GetUserInfoByIDs(menberIDs)
	if err != nil {
		return nil, err
	}

	// handle data response
	res := make([]dto.GetGroupRoomResponse, 0)
	for _, r := range rooms {
		if ids, ok := roomMembersMap[r.RoomID]; ok {
			dtoRoom := dto.GetGroupRoomResponse{
				RoomID:          r.RoomID,
				RoomName:        r.RoomName.String,
				RoomDescription: r.RoomDescription.String,
				RoomAvatar:      r.RoomAvatar.String,
				RoomIsGroup:     true,
				RoomMembers:     []dto.GroupMember{},
				RoomLastMessage: dto.GroupMessage{
					MessageID:       r.MessageID,
					MessageContent:  r.MessageContent,
					MessageSenderID: r.MessageSenderID,
					MessageType:     r.MessageType,
					MessageSentAt:   r.MessageSentAt,
				},
			}
			for _, id := range ids {
				for _, us := range users {
					if us.UserID == int64(id) {
						// Get user status
						status, _ := s.Urepo.GetStatusByUserId(us.UserID)
						dtoRoom.RoomMembers = append(dtoRoom.RoomMembers, dto.GroupMember{
							UserID:       us.UserID,
							UserName:     us.UserNickname,
							UserAvatar:   us.UserAvatar,
							UserStatus:   status, // Default status, will be updated later
							UserNickname: r.MemberNickname.String,
							Role:         r.MemberRole,
						})
						break // Found the user, no need to continue
					}
				}
			}
			delete(roomMembersMap, r.RoomID) // Remove processed room
			res = append(res, dtoRoom)
		}
	}

	fmt.Println("GetRoomsByUserID: rooms found:", time.Since(now).Milliseconds(), "ms")

	return res, nil
}

func (s *chatService) CreateRoomChat(data *dto.CreateRoomDTO) (dto.CreateRoomDTO, error) {

	// Check if room already exists
	roomID := uint64(0)
	if data.RoomIsGroup {
		room, err := s.Crepo.GetRoomGroupByName(data.RoomName)
		roomID = room.RoomID
		if !errors.Is(err, sql.ErrNoRows) {
			return dto.CreateRoomDTO{}, err
		}
	} else {
		room, err := s.Crepo.GetRoomByName(data.RoomName)
		roomID = room.RoomID
		if !errors.Is(err, sql.ErrNoRows) {
			return dto.CreateRoomDTO{}, err
		}
	}
	if roomID != 0 {
		return dto.CreateRoomDTO{
			RoomID:          roomID,
			RoomName:        data.RoomName,
			RoomIsGroup:     data.RoomIsGroup,
			RoomMembers:     data.RoomMembers,
			RoomAvatar:      data.RoomAvatar,
			RoomDescription: data.RoomDescription,
		}, errors.New("room already exists")
	}

	err := s.Crepo.CreateRoom(data)
	if err != nil {
		return dto.CreateRoomDTO{}, err
	}
	// add first message for group room
	if data.RoomIsGroup {
		newMessageGroup := dto.SaveMessageDTO{
			MessageID:         0,
			MessageRoomID:     data.RoomID,
			MessageSenderID:   uint64(data.RoomCreateBy),
			MessageReceiverID: 0, // No receiver in group chat
			MessageContent:    "[system] Welcome to the group chat! ",
			MessageType:       "text", // Default type for group room creation
			MessageSentAt:     time.Now(),
		}
		msgID, err := s.Crepo.SaveMessegeGroup(newMessageGroup)
		if err != nil {
			fmt.Println("CreateRoomChat 1: error saving group message", err)
			return dto.CreateRoomDTO{}, err
		}

		data.RoomMessageID = msgID
	}

	// Create room members
	err = s.Crepo.AddMembersToRoom(data.RoomID, uint64(data.RoomCreateBy), data.RoomMembers)
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

// optimize to O(n^2)
func (s *chatService) GetGroupRoomsByUserID2(userInfo *model.GoDbUserInfo) ([]dto.GetGroupRoomResponse, error) {
	userID := userInfo.UserID
	rooms, err := s.Crepo.GetGroupRoomsByUserId(userID)
	if err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return []dto.GetGroupRoomResponse{}, nil
	}

	// Get map roomID -> RoomResponse and Get 1 array Set memberIDs
	roomMenbersMap := make(map[uint64]*dto.GetGroupRoomResponse)
	memberIDs := make([]uint64, 0)
	memberMapSet := make(map[uint64]*dto.GroupMember)
	memberIDsMap := make(map[uint64]struct{})

	for _, r := range rooms {
		if _, ok := roomMenbersMap[r.RoomID]; !ok {
			// create new room response
			roomMenbersMap[r.RoomID] = &dto.GetGroupRoomResponse{
				RoomID:          r.RoomID,
				RoomName:        r.RoomName.String,
				RoomDescription: r.RoomDescription.String,
				RoomAvatar:      r.RoomAvatar.String,
				RoomCreatedBy:   r.RoomCreatedBy.Int64,
				RoomCreatedAt:   r.RoomCreatedAt.Time,
				RoomIsGroup:     true,
				RoomMembers:     []dto.GroupMember{},
				RoomLastMessage: dto.GroupMessage{
					MessageID:       r.MessageID,
					MessageContent:  r.MessageContent,
					MessageSenderID: r.MessageSenderID,
					MessageType:     r.MessageType,
					MessageSentAt:   r.MessageSentAt,
				},
			}
		}

		// Add member to room
		room := roomMenbersMap[r.RoomID]
		newMember := dto.GroupMember{
			UserID:         int64(r.MemberUserID),
			UserNickname:   r.MemberNickname.String,
			Role:           r.MemberRole,
			UserAvatar:     "",
			UserName:       "",
			UserStatus:     "offline", // Default status, will be updated later
			MemberLastSeen: r.MemberLastSeen.Int64,
		}

		if r.MemberUserID == uint64(userID) {
			newMember.UserName = userInfo.UserNickname
			newMember.UserAvatar = userInfo.UserAvatar
			newMember.UserStatus = "online"

			room.CurrentLastSeen = r.MemberLastSeen.Int64
		}

		room.RoomMembers = append(room.RoomMembers, newMember)

		if r.MemberUserID != uint64(userID) {
			if _, exits := memberIDsMap[r.MemberUserID]; !exits {
				memberIDsMap[r.MemberUserID] = struct{}{}
				memberIDs = append(memberIDs, r.MemberUserID)
				memberMapSet[r.MemberUserID] = &dto.GroupMember{
					UserID:       int64(r.MemberUserID),
					UserNickname: r.MemberNickname.String,
					Role:         r.MemberRole,
					UserAvatar:   "",
					UserName:     "",
					UserStatus:   "offline", // Default status, will be updated later
				}
			}
		}
	}

	// Get user info by member ids
	users, err := s.Urepo.GetUserInfoByIDs(memberIDs)
	if err != nil {
		return nil, err
	}

	// Update user info in room members
	for _, user := range users {
		if member, ok := memberMapSet[uint64(user.UserID)]; ok {
			member.UserName = user.UserNickname
			member.UserAvatar = user.UserAvatar

			// Get user status
			status, _ := s.Urepo.GetStatusByUserId(user.UserID)
			member.UserStatus = status
		}
	}

	// handle data response
	res := make([]dto.GetGroupRoomResponse, 0)
	for _, room := range roomMenbersMap {
		for i, member := range room.RoomMembers {
			if memberInfo, ok := memberMapSet[uint64(member.UserID)]; ok {
				room.RoomMembers[i] = *memberInfo
			}
		}
		res = append(res, *room)
	}

	return res, nil
}
