package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/middleware"
	"github.com/k1ender/task-master-go/internal/response"
	"github.com/k1ender/task-master-go/internal/storage"
)

type UserHandler struct {
	store    *storage.Storage
	validate *validator.Validate
	config   *config.Config
	log      *slog.Logger
}

func NewUserHandler(store *storage.Storage, validator *validator.Validate, config *config.Config, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		store:    store,
		validate: validator,
		config:   config,
		log:      logger,
	}
}

// @Summary Get user details
// @Description Get user details
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} response.Response
// @Router /user [get]
// @Security ApiKeyAuth
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetAuthUserFromContext(r.Context())
	response.OK(w, user)
}
