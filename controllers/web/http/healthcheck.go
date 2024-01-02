package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type healthcheckHttpControllers struct {
	BaseHandler
	router chi.Router
}

func SetupHealthcheckHttpControllers() *healthcheckHttpControllers {
	return &healthcheckHttpControllers{
		router: chi.NewRouter(),
	}
}

func (h *healthcheckHttpControllers) Routes() chi.Router {

	h.router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		h.JSON(w, http.StatusOK, "Working")
	})

	return h.router
}
