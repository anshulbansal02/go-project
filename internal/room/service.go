package room

import "context"

type RoomService struct {
	roomRepo *RoomRepository
}

func NewService(roomRepo *RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

/**
- Create Room
- Get Room
- Request User join Room
- Accept User join Room
- Reject User join Room
- Cancel User join Roomcvfv
-

*/

// // User Service Methods
func (s *RoomService) CreatePrivateRoom(ctx context.Context, adminId string) (*Room, error) {

	room := s.roomRepo.NewRoom(&adminId, []string{adminId}, "private")

	if err := s.roomRepo.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) CreatePublicRoom(ctx context.Context, adminId string) (*Room, error) {
	room := s.roomRepo.NewRoom(nil, []string{}, "public")

	if err := s.roomRepo.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
