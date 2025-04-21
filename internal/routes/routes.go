package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/handlers"
	"gorm.io/gorm"
)

func New(db *gorm.DB, config *config.Config) *chi.Mux {
	router := chi.NewRouter()

	validator := validator.New(validator.WithRequiredStructEnabled())

	router.Post("/register",
		handlers.NewAuthHandler(db, validator, config).RegisterUser,
	)

	return router
}
