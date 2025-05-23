package response

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/utils"
)

type Response struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func WriteResponse(w http.ResponseWriter, status int, data any, message string, success bool) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return utils.WriteJSON(
		w,
		status,
		Response{Success: success, Status: status, Data: data, Message: message},
	)
}

func OK(w http.ResponseWriter, data any) error {
	return WriteResponse(w, http.StatusOK, data, "", true)
}

func Created(w http.ResponseWriter, data any) error {
	return WriteResponse(w, http.StatusCreated, data, "", true)
}

func BadRequest(w http.ResponseWriter, message string) error {
	return WriteResponse(w, http.StatusBadRequest, nil, message, false)
}

func InternalServerError(w http.ResponseWriter) error {
	return WriteResponse(w, http.StatusInternalServerError, nil, "Internal Server Error", false)
}

func NotFound(w http.ResponseWriter, message string) error {
	return WriteResponse(w, http.StatusNotFound, nil, message, false)
}

func Unauthorized(w http.ResponseWriter, message string) error {
	return WriteResponse(w, http.StatusUnauthorized, nil, message, false)
}

func NoContent(w http.ResponseWriter) error {
	return WriteResponse(w, http.StatusNoContent, nil, "", true)
}

type Error struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

func ValidationError(w http.ResponseWriter, err []validator.FieldError) error {
	var errs []Error = make([]Error, len(err))

	for i, e := range err {
		errs[i] = Error{
			Field: e.Field(), Errors: append(errs[i].Errors, e.Error()),
		}
	}

	return WriteResponse(w, http.StatusBadRequest, errs, "Unprocessable Entity", false)
}
