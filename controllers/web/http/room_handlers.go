package web

import (
	"anshulbansal02/scribbly/internal/room"

	"github.com/go-chi/chi/v5"
)

type roomHttpControllers struct {
	roomService room.RoomService
	router      chi.Router
}

func SetupRoomHttpControllers(roomService room.RoomService) *roomHttpControllers {
	return &roomHttpControllers{
		roomService: roomService,
		router:      chi.NewRouter(),
	}

}

func (h *roomHttpControllers) Routes() chi.Router {
	return h.router
}
