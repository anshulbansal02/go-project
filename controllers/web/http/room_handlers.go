package web

import (
	"anshulbansal02/scribbly/internal/room"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type roomHttpControllers struct {
	BaseHandler
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

	h.router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		// Create new room only if there exists a user client
	})

	h.router.Get("/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		// Get room by id

		roomId := chi.URLParam(r, "roomId")

		room, err := h.roomService.GetRoom(r.Context(), roomId)

		if err != nil {
		}

		h.JSON(w, http.StatusOK, room)

	})

	h.router.Post("/join/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		// Create request to join room
		h.roomService.CreateJoinRequest(r.Context(), "", "")
	})

	h.router.Delete("/join", func(w http.ResponseWriter, r *http.Request) {
		// Cancel join request
	})

	h.router.Post("/leave", func(w http.ResponseWriter, r *http.Request) {
		// Leave room
	})

	return h.router
}
