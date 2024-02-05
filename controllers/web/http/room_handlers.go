package web

import (
	"anshulbansal02/scribbly/controllers/middlewares"
	"anshulbansal02/scribbly/internal/room"
	"anshulbansal02/scribbly/internal/user"
	"errors"
	"fmt"
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
		var body struct {
			AdminId string `json:"adminId"`
		}
		err := h.DecodeBodyTo(r, &body)
		if err != nil {
			h.JSON(w, http.StatusBadRequest, err.Error())
			return
		}

		room, err := h.roomService.CreatePrivateRoom(r.Context(), body.AdminId)
		fmt.Println(err)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				h.JSON(w, http.StatusBadRequest, err.Error())
			}
			h.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		h.JSON(w, http.StatusCreated, room)
	})

	h.router.Get("/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "roomId")

		room, err := h.roomService.GetRoom(r.Context(), roomId)
		if err != nil {
			h.JSON(w, http.StatusInternalServerError, err)
			return
		}

		h.JSON(w, http.StatusOK, room)
	})

	h.router.Post("/join/{roomId}", middlewares.WithAuthorization(func(w http.ResponseWriter, r *http.Request) {
		roomId := chi.URLParam(r, "roomId")

		user := r.Context().Value(middlewares.UserCtxKey).(middlewares.UserAuthContext)

		// Create request to join room
		err := h.roomService.CreateJoinRequest(r.Context(), roomId, user.UserId)

		if err != nil {
			h.JSON(w, http.StatusBadRequest, err.Error())
			return
		}

		h.JSON(w, http.StatusProcessing, "Requested")
	}))

	h.router.Delete("/join", func(w http.ResponseWriter, r *http.Request) {
		// Cancel join request
	})

	h.router.Post("/leave", func(w http.ResponseWriter, r *http.Request) {
		// Leave room
	})

	return h.router
}
