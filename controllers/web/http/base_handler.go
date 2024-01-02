package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseHandler struct {
}

func (h *BaseHandler) JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json")

	json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Server Error: %w", err)))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(json)
}

func (h *BaseHandler) DecodeBodyTo(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)

	return err
}
