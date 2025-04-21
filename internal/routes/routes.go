package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/handlers"
	"github.com/k1ender/task-master-go/internal/middleware"
	"gorm.io/gorm"
)

func New(db *gorm.DB, config *config.Config) *chi.Mux {
	router := chi.NewRouter()

	validator := validator.New(validator.WithRequiredStructEnabled())
	authHandlers := handlers.NewAuthHandler(db, validator, config)
	taskHandlers := handlers.NewTaskHandler(db, validator, config)

	authMiddleware := middleware.Auth(db, config.JWT.Secret)

	router.Post("/register",
		authHandlers.RegisterUser,
	)
	router.Post("/login",
		authHandlers.LoginUser,
	)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/tasks", taskHandlers.CreateTask)
	})

	return router
}
