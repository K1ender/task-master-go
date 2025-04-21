package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/middleware"
	"github.com/k1ender/task-master-go/internal/response"
	"gorm.io/gorm"
)

type UserHandler struct {
	db       *gorm.DB
	validate *validator.Validate
	config   *config.Config
}

func NewUserHandler(db *gorm.DB, validator *validator.Validate, config *config.Config) *UserHandler {
	return &UserHandler{
		db:       db,
		validate: validator,
		config:   config,
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
