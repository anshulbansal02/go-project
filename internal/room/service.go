package room

import (
	roomuser "anshulbansal02/scribbly/internal/room/aggregates/room_user"
	"context"
)

type RoomService struct {
	roomRepo         *RoomRepository
	userRoomRelation *roomuser.UserRoomRelationRepository
	joinRequests     *roomuser.UserJoinRequestRepository

	UserEventsChannel chan UserEvent
}

func NewService(roomRepo *RoomRepository, userRoomRelation *roomuser.UserRoomRelationRepository, joinRequestsRepo *roomuser.UserJoinRequestRepository) *RoomService {
	return &RoomService{
		roomRepo:         roomRepo,
		userRoomRelation: userRoomRelation,
		joinRequests:     joinRequestsRepo,

		UserEventsChannel: make(chan UserEvent),
	}
}

// // User Service Methods
func (s *RoomService) CreatePrivateRoom(ctx context.Context, adminId string) (*Room, error) {

	room := s.roomRepo.NewRoom(&adminId, "private")

	if err := s.roomRepo.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) CreatePublicRoom(ctx context.Context, adminId string) (*Room, error) {
	room := s.roomRepo.NewRoom(nil, "public")

	if err := s.roomRepo.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) GetRoom(ctx context.Context, roomId string) (*Room, error) {
	room, err := s.roomRepo.GetRoom(ctx, roomId)

	if err != nil {
		return nil, err
	}

	users, err := s.userRoomRelation.GetUsersByRoomId(ctx, roomId)
	if err != nil {
		return nil, err
	}

	room.Participants = users

	return room, nil
}

func (s *RoomService) GetRoomUsers(ctx context.Context, roomId string) ([]string, error) {
	users, err := s.userRoomRelation.GetUsersByRoomId(ctx, roomId)
	return users, err
}

func (s *RoomService) DeleteRoom(ctx context.Context, roomId string) error {
	defer s.roomRepo.LockKey(roomId)()

	return s.roomRepo.DeleteRoom(ctx, roomId)
}

func (s *RoomService) GetRoomAdmin(ctx context.Context, roomId string) *string {
	room, _ := s.roomRepo.GetRoom(ctx, roomId)

	return room.Admin
}

func (s *RoomService) CreateRoomJoiningRequest(ctx context.Context, roomId string, userId string) error {

	defer s.roomRepo.LockKey(roomId)()

	// Check if room exists
	roomExists, err := s.roomRepo.RoomExists(ctx, roomId)
	if err != nil {
		return err
	}
	if !roomExists {
		return ErrRoomNotFound
	}

	// Check if user is not in any other room
	rId, err := s.userRoomRelation.GetRoomIdByUserId(ctx, userId)
	if err != nil {
		return err
	}
	if rId != "" {
		return ErrUserAlreadyInRoom
	}

	// Create new join request
	err = s.joinRequests.CreateJoinRequest(ctx, userId, roomId)
	if err != nil {
		return err
	}

	// Send request to room user events channel
	s.UserEventsChannel <- UserEvent{
		Type:   JoinRequestEvent,
		UserId: userId,
		RoomId: roomId,
	}

	return nil

}

func (s *RoomService) CancelJoinRequest(ctx context.Context, userId string) {

}
