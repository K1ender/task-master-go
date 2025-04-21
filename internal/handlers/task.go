package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/middleware"
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
	user := middleware.GetAuthUserFromContext(r.Context())
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
		Title:  payload.Title,
		Body:   payload.Body,
		UserID: user.ID,
	}

	res := h.db.Create(&task)

	if res.Error != nil {
		response.InternalServerError(w)
		return
	}

	response.Created(w, task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetAuthUserFromContext(r.Context())

	var tasks []models.Task

	res := h.db.Where("user_id = ?", user.ID).Find(&tasks)

	if res.Error != nil {
		response.InternalServerError(w)
		return
	}

	response.OK(w, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	task := middleware.GetTaskFromContext(r.Context())

	response.OK(w, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	task := middleware.GetTaskFromContext(r.Context())

	res := h.db.Delete(&task)

	if res.Error != nil {
		response.InternalServerError(w)
		return
	}

	response.NoContent(w)
}

type UpdateTaskRequest struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetAuthUserFromContext(r.Context())
	task := middleware.GetTaskFromContext(r.Context())
	var payload UpdateTaskRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	updates := map[string]interface{}{}

	if payload.Title != "" {
		updates["title"] = payload.Title
	}

	if payload.Body != "" {
		updates["body"] = payload.Body
	}

	if payload.Completed != task.Completed {
		updates["completed"] = payload.Completed
	}

	if len(updates) == 0 {
		response.OK(w, task)
		return
	}

	res := h.db.Model(&task).Where("user_id = ?", user.ID).Updates(updates)

	if res.Error != nil {
		response.InternalServerError(w)
		return
	}

	response.OK(w, task)
}
