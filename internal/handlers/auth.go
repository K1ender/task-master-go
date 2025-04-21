package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"github.com/k1ender/task-master-go/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db       *gorm.DB
	validate *validator.Validate
	config   *config.Config
}

func NewAuthHandler(db *gorm.DB, validator *validator.Validate, config *config.Config) *AuthHandler {
	return &AuthHandler{
		db:       db,
		validate: validator,
		config:   config,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	user := models.User{
		Username: payload.Username,
		Password: string(hashed_password),
	}

	res := h.db.Create(&user)

	if res.Error != nil {
		if err, ok := res.Error.(*pgconn.PgError); ok {
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
		response.InternalServerError(w)
		return
	}

	response.Created(w, map[string]string{"token": ss})
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserRequest
	if err := utils.ReadJSON(r, &payload); err != nil {
		response.BadRequest(w, "Bad Request")
		return
	}

	if err := h.validate.Struct(payload); err != nil {
		response.ValidationError(w, err.(validator.ValidationErrors))
		return
	}

	var user models.User

	res := h.db.Where("username = ?", payload.Username).First(&user)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			response.NotFound(w, "User not found")
			return
		}
		response.InternalServerError(w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		response.Unauthorized(w, "Unauthorized")
		return
	}

	ss, err := utils.SignToken(user.ID, h.config.JWT.Secret)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.OK(w, map[string]string{"token": ss})
}
