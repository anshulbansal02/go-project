package web

import (
	"anshulbansal02/scribbly/internal/user"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type userHttpControllers struct {
	BaseHandler
	userService user.UserService
	router      chi.Router
}

func SetupUserHttpControllers(userService user.UserService) *userHttpControllers {
	return &userHttpControllers{
		userService: userService,
		router:      chi.NewRouter(),
	}
}

func (h *userHttpControllers) Routes() chi.Router {
	// @POST / - Create a user
	h.router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
		}
		err := h.DecodeBodyTo(r, &body)
		if err != nil {
			h.JSON(w, http.StatusBadRequest, err)
			return
		}

		user, err := h.userService.CreateUser(r.Context(), body.Username)
		if err != nil {
			h.JSON(w, http.StatusInternalServerError, err)
			return
		}

		h.JSON(w, http.StatusCreated, user)
	})

	// @GET /<userId> - Get user info by id
	h.router.Get("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")

		user, err := h.userService.GetUser(r.Context(), userId)

		if err != nil {
			fmt.Println(err)
			h.JSON(w, http.StatusNotFound, err)
			return
		}

		h.JSON(w, http.StatusOK, user)
	})

	// @PATCH /<userId> - Update user info
	h.router.Patch("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
		}
		err := h.DecodeBodyTo(r, &body)
		if err != nil {
			h.JSON(w, http.StatusBadRequest, err)
			return
		}

		userId := chi.URLParam(r, "userId")

		err = h.userService.UpdateUserName(r.Context(), userId, body.Username)
		if err != nil {
			h.JSON(w, http.StatusNotFound, err)
			return
		}

		h.JSON(w, http.StatusOK, nil)
	})

	return h.router
}
