package web

import (
	"anshulbansal02/scribbly/internal/user"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type userHttpControllers struct {
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
	h.router.Post("/", func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Username string `json:"username"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		user, err := h.userService.CreateUser(r.Context(), body.Username)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	h.router.Get("/{userId}", func(w http.ResponseWriter, r *http.Request) {

		userId := chi.URLParam(r, "userId")

		if userId == "" {
			c.JSON(http.StatusBadRequest, nil)
		}

		user, err := h.userService.GetUser(r.Context(), userId)

		if err != nil {
			c.JSON(http.StatusNotFound, err)
		}

		c.JSON(http.StatusOK, user)
	})

	return h.router
}
