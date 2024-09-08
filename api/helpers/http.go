package helpers

import (
	"encoding/json"
	"net/http"
)

func Reply(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if body == nil {
		return
	}

	var response interface{}

	switch v := body.(type) {
	case string:
		response = struct {
			Message string `json:"message"`
		}{
			Message: v,
		}
	case error:
		response = struct {
			Error string `json:"error"`
		}{
			Error: v.Error(),
		}
	default:
		response = body
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
