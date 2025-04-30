package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"github.com/k1ender/task-master-go/internal/storage"
	"github.com/k1ender/task-master-go/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	store    *storage.Storage
	validate *validator.Validate
	config   *config.Config
	log      *slog.Logger
}

func NewAuthHandler(store *storage.Storage, validator *validator.Validate, config *config.Config, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		store:    store,
		validate: validator,
		config:   config,
		log:      logger,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register [post]
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserRequest
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

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("failed to hash password", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	user := models.User{
		Username: payload.Username,
		Password: string(hashed_password),
	}

	err = h.store.Users.CreateUser(&user)

	if err != nil {
		h.log.Error("failed to create user", slog.Any("error", err))
		if err, ok := err.(*pgconn.PgError); ok {
			if err.ConstraintName == "uni_users_username" {
				response.BadRequest(w, "Username already exists")
				return
			}
		}
		response.InternalServerError(w)
		return
	}

	ss, err := utils.SignToken(user.ID, h.config.JWT.Secret)
	if err != nil {
		h.log.Error("failed to sign token", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	response.Created(w, map[string]string{"token": ss})
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// @Summary Login a user
// @Description Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body LoginUserRequest true "User details"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [post]
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserRequest
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

	user, err := h.store.Users.GetUserByUsername(payload.Username)

	if err != nil {
		h.log.Error("failed to get user", slog.Any("error", err))
		if err == gorm.ErrRecordNotFound {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalServerError(w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		h.log.Error("failed to compare password", slog.Any("error", err))
		response.Unauthorized(w, "Unauthorized")
		return
	}

	ss, err := utils.SignToken(user.ID, h.config.JWT.Secret)
	if err != nil {
		h.log.Error("failed to sign token", slog.Any("error", err))
		response.InternalServerError(w)
		return
	}

	response.OK(w, map[string]string{"token": ss})
}
