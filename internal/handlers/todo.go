package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"github.com/k1ender/task-master-go/internal/utils"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db       *gorm.DB
	validate *validator.Validate
	config   *config.Config
}

func NewTaskHandler(db *gorm.DB, validator *validator.Validate, config *config.Config) *TaskHandler {
	return &TaskHandler{
		db:       db,
		validate: validator,
		config:   config,
	}
}

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var payload CreateTaskRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	task := models.Task{
		Title: payload.Title,
		Body:  payload.Body,
	}

	res := h.db.Create(&task)

	if res.Error != nil {
		response.InternalServerError(w)
		return
	}

	response.Created(w, task)
}
