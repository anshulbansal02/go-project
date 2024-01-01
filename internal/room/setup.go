package room

import (
	"anshulbansal02/scribbly/internal/repository"
	aggJoinRequest "anshulbansal02/scribbly/internal/room/aggregates/user_join_request"
	aggUserRoom "anshulbansal02/scribbly/internal/room/aggregates/user_room_relation"
)

func SetupConcreteService(repository repository.Repository) *RoomService {
	roomService := NewService(NewRepository(repository),
		aggUserRoom.NewUserRoomRelation(repository),
		aggJoinRequest.NewUserJoinRequestRepository(repository),
	)

	return roomService
}
