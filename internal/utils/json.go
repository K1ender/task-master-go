package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, data any) error {
	encoder := json.NewDecoder(r.Body)

	encoder.DisallowUnknownFields()

	return encoder.Decode(data)
}
