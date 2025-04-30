package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/middleware"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"github.com/k1ender/task-master-go/internal/storage"
	"github.com/k1ender/task-master-go/internal/utils"
)

type TaskHandler struct {
	store    *storage.Storage
	validate *validator.Validate
	config   *config.Config
	log      *slog.Logger
}

func NewTaskHandler(store *storage.Storage, validator *validator.Validate, config *config.Config, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		store:    store,
		validate: validator,
		config:   config,
		log:      logger,
	}
}

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

// @Summary Create a new task
// @Description Create a new task
// @Tags Task
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /tasks [post]
// @Security ApiKeyAuth
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetAuthUserFromContext(r.Context())
	var payload CreateTaskRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		h.log.Error("failed to read request body", slog.Any("error", err))
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		h.log.Error("failed to validate request body", slog.Any("error", err))
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	task := models.Task{
		Title:  payload.Title,
		Body:   payload.Body,
		UserID: user.ID,
	}

	_, err := h.store.Tasks.CreateTask(&task)

	if err != nil {
		h.log.Error("failed to create task", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	response.Created(w, task)
}

// @Summary Get all tasks for a user
// @Description Get all tasks for a user
// @Tags Task
// @Accept json
// @Produce json
// @Success 200 {object} []models.Task
// @Failure 500 {object} response.Response
// @Router /tasks [get]
// @Security ApiKeyAuth
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetAuthUserFromContext(r.Context())

	tasks, err := h.store.Tasks.GetTasks(user.ID)

	if err != nil {
		h.log.Error("failed to get tasks", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	response.OK(w, tasks)
}

// @Summary Get a task by ID
// @Description Get a task by ID
// @Tags Task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /tasks/{id} [get]
// @Security ApiKeyAuth
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	task := middleware.GetTaskFromContext(r.Context())

	response.OK(w, task)
}

// @Summary Delete a task by ID
// @Description Delete a task by ID
// @Tags Task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204
// @Failure 500 {object} response.Response
// @Router /tasks/{id} [delete]
// @Security ApiKeyAuth
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	task := middleware.GetTaskFromContext(r.Context())

	err := h.store.Tasks.DeleteTask(task.ID)

	if err != nil {
		h.log.Error("failed to delete task", slog.Any("error", err))
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

// @Summary Update a task by ID
// @Description Update a task by ID
// @Tags Task
// @Accept json
// @Produce json
// @Param task body UpdateTaskRequest true "Task details"
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /tasks/{id} [patch]
// @Security ApiKeyAuth
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// user := middleware.GetAuthUserFromContext(r.Context())
	task := middleware.GetTaskFromContext(r.Context())
	var payload UpdateTaskRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		h.log.Error("failed to read request body", slog.Any("error", err))
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		h.log.Error("failed to validate request body", slog.Any("error", err))
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	updates := map[string]any{}

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

	err := h.store.Tasks.UpdateTask(task, updates)

	if err != nil {
		h.log.Error("failed to update task", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	response.OK(w, task)
}
