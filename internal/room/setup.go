package room

import (
	"anshulbansal02/scribbly/internal/repository"
	aggJoinRequest "anshulbansal02/scribbly/internal/room/aggregates/user_join_request"
	aggUserRoom "anshulbansal02/scribbly/internal/room/aggregates/user_room_relation"
)

func NewRepository(repo repository.Repository) *RoomRepository {
	return &RoomRepository{
		Repository: repo,
	}
}

func NewService(
	roomRepo *RoomRepository,
	userRoomRelationRepo *aggUserRoom.UserRoomRelationRepository,
	joinRequestsRepo *aggJoinRequest.UserJoinRequestRepository,

	roomCodeIdMap *RoomCodeIdMapRepository,
) *RoomService {
	return &RoomService{
		roomRepo:             roomRepo,
		userRoomRelationRepo: userRoomRelationRepo,
		joinRequestsRepo:     joinRequestsRepo,

		UserEventsChannel: make(chan UserEvent),

		roomCodeIdMap: roomCodeIdMap,
	}
}

func (r *RoomService) SetDependencies(services DependingServices) {
	r.DependingServices = services
}

func SetupConcreteService(repository repository.Repository) *RoomService {
	roomService := NewService(NewRepository(repository),
		aggUserRoom.NewUserRoomRelation(repository),
		aggJoinRequest.NewUserJoinRequestRepository(repository),
		&RoomCodeIdMapRepository{Repository: repository},
	)

	return roomService
}
