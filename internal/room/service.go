package room

import (
	aggJoinRequest "anshulbansal02/scribbly/internal/room/aggregates/user_join_request"
	aggUserRoom "anshulbansal02/scribbly/internal/room/aggregates/user_room_relation"
	"anshulbansal02/scribbly/internal/user"
	"anshulbansal02/scribbly/pkg/repository"
	"errors"

	"context"
)

type DependingServices struct {
	UserService *user.UserService
}

type RoomService struct {
	roomRepo             *RoomRepository
	userRoomRelationRepo *aggUserRoom.UserRoomRelationRepository
	joinRequestsRepo     *aggJoinRequest.UserJoinRequestRepository

	DependingServices

	UserEventsChannel chan UserEvent

	roomCodeIdMap *RoomCodeIdMapRepository
}

/********************** Service Methods **********************/

func (s *RoomService) CreatePrivateRoom(ctx context.Context, adminId string) (*Room, error) {
	adminExists, err := s.UserService.UserExists(ctx, adminId)
	if err != nil {
		return nil, err
	}
	if !adminExists {
		return nil, user.ErrUserNotFound
	}

	_, err = s.userRoomRelationRepo.GetRoomIdByUserId(ctx, adminId)
	if err == nil {
		return nil, ErrUserAlreadyInRoom
	} else {
		if !errors.Is(err, repository.ErrEntityNotFound) {
			return nil, err
		}
	}

	room := s.roomRepo.NewRoom(&adminId, "private")

	if err := s.roomCodeIdMap.Set(ctx, room.Code, room.ID); err != nil {
		return nil, err
	}

	if err := s.userRoomRelationRepo.AddUserToRoom(ctx, room.ID, adminId); err != nil {
		return nil, err
	}

	if err := s.roomRepo.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) CreatePublicRoom(ctx context.Context) (*Room, error) {
	room := s.roomRepo.NewRoom(nil, "public")

	if err := s.roomCodeIdMap.Set(ctx, room.Code, room.ID); err != nil {
		return nil, err
	}

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

	users, err := s.userRoomRelationRepo.GetUsersByRoomId(ctx, roomId)
	if err != nil {
		return nil, err
	}

	room.Participants = users

	return room, nil
}

// Decide on returning room does not exist error
func (s *RoomService) GetRoomUsers(ctx context.Context, roomId string) ([]string, error) {
	userIds, err := s.userRoomRelationRepo.GetUsersByRoomId(ctx, roomId)
	return userIds, err
}

func (s *RoomService) DeleteRoom(ctx context.Context, roomId string) error {
	defer s.roomRepo.KeyMutex.Lock(roomId)()

	return s.roomRepo.DeleteRoom(ctx, roomId)
}

func (s *RoomService) GetRoomAdmin(ctx context.Context, roomId string) *string {
	room, _ := s.roomRepo.GetRoom(ctx, roomId)
	if room == nil {
		return nil
	}
	return room.Admin
}

func (s *RoomService) CreateJoinRequest(ctx context.Context, roomCode string, userId string) error {
	roomId, err := s.roomCodeIdMap.GetRoomId(ctx, roomCode)
	if err != nil {
		return err
	}

	defer s.roomRepo.KeyMutex.Lock(roomId)()

	// Check if room exists
	roomExists, err := s.roomRepo.RoomExists(ctx, roomId)
	if err != nil {
		return err
	}
	if !roomExists {
		return ErrRoomNotFound
	}

	// Check if user is not in any other room
	_, err = s.userRoomRelationRepo.GetRoomIdByUserId(ctx, userId)

	if err == nil {
		return ErrUserAlreadyInRoom
	} else {
		if !errors.Is(err, repository.ErrEntityNotFound) {
			return err
		}
	}

	// Create new join request
	err = s.joinRequestsRepo.CreateJoinRequest(ctx, userId, roomId)
	if err != nil {
		return err
	}

	// Send request to room user events channel
	s.UserEventsChannel <- UserEvent{
		Type:   JoinRequestedEvent,
		UserId: userId,
		RoomId: roomId,
	}

	return nil

}

func (s *RoomService) CancelJoinRequest(ctx context.Context, userId string) error {
	//Assuming delete only happens if request exists
	roomId, err := s.joinRequestsRepo.GetUserJoinRequestedRoom(ctx, userId)
	if err != nil {
		return err
	}

	err = s.joinRequestsRepo.DeleteJoinRequest(ctx, userId)
	if err != nil {
		return err
	}

	s.UserEventsChannel <- UserEvent{
		Type:   JoinRequestCancelledEvent,
		UserId: userId,
		RoomId: roomId,
	}

	return nil
}

func (s *RoomService) RejectJoinRequest(ctx context.Context, userId string) error {
	// delete request entity and send channel event
	roomId, err := s.joinRequestsRepo.GetUserJoinRequestedRoom(ctx, userId)
	if err != nil {
		return err
	}

	err = s.joinRequestsRepo.DeleteJoinRequest(ctx, userId)
	if err != nil {
		return err
	}

	s.UserEventsChannel <- UserEvent{
		Type:   JoinRequestRejectedEvent,
		UserId: userId,
		RoomId: roomId,
	}

	return nil
}

func (s *RoomService) AcceptJoinRequest(ctx context.Context, userId string) error {
	// delete request entity, add user to room and send channel event
	roomId, err := s.joinRequestsRepo.GetUserJoinRequestedRoom(ctx, userId)
	if err != nil {
		return err
	}

	defer s.roomRepo.KeyMutex.Lock(roomId)()

	err = s.joinRequestsRepo.DeleteJoinRequest(ctx, userId)
	if err != nil {
		return err
	}

	err = s.userRoomRelationRepo.AddUserToRoom(ctx, roomId, userId)
	if err != nil {
		return err
	}

	s.UserEventsChannel <- UserEvent{
		Type:   UserJoinedEvent,
		UserId: userId,
		RoomId: roomId,
	}

	return nil
}

func (s *RoomService) GetUserRoomId(ctx context.Context, userId string) (string, error) {
	roomId, err := s.userRoomRelationRepo.GetRoomIdByUserId(ctx, userId)
	return roomId, err
}
