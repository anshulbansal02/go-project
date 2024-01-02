package room

import (
	"anshulbansal02/scribbly/internal/repository"
	aggJoinRequest "anshulbansal02/scribbly/internal/room/aggregates/user_join_request"
	aggUserRoom "anshulbansal02/scribbly/internal/room/aggregates/user_room_relation"
	"anshulbansal02/scribbly/internal/user"
)

func NewRepository(repo repository.Repository) *RoomRepository {
	return &RoomRepository{
		Repository: repo,
	}
}

func NewService(
	roomRepo *RoomRepository,
	userRoomRelation *aggUserRoom.UserRoomRelationRepository,
	joinRequestsRepo *aggJoinRequest.UserJoinRequestRepository,
	userService *user.UserService,
) *RoomService {
	return &RoomService{
		roomRepo:         roomRepo,
		userRoomRelation: userRoomRelation,
		joinRequests:     joinRequestsRepo,

		userService: userService,

		UserEventsChannel: make(chan UserEvent),
	}
}

func SetupConcreteService(repository repository.Repository, userService *user.UserService) *RoomService {
	roomService := NewService(NewRepository(repository),
		aggUserRoom.NewUserRoomRelation(repository),
		aggJoinRequest.NewUserJoinRequestRepository(repository),
		userService,
	)

	return roomService
}
