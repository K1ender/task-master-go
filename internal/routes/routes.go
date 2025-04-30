package routes

import (
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/docs"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/handlers"
	"github.com/k1ender/task-master-go/internal/middleware"
	"github.com/k1ender/task-master-go/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"gorm.io/gorm"
)

func New(db *gorm.DB, config *config.Config, store *storage.Storage, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", config.HttpServer.Port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "Task Master API"
	docs.SwaggerInfo.Version = "1.0"

	validator := validator.New(validator.WithRequiredStructEnabled())

	userHandlers := handlers.NewUserHandler(store, validator, config, logger)
	authHandlers := handlers.NewAuthHandler(store, validator, config, logger)
	taskHandlers := handlers.NewTaskHandler(store, validator, config, logger)

	authMiddleware := middleware.Auth(db, config.JWT.Secret)
	taskMiddleware := middleware.TaskMiddleware(db)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf(
				"http://localhost:%s/swagger/doc.json",
				config.HttpServer.Port,
			),
		),
	))

	r.Post("/register",
		authHandlers.RegisterUser,
	)
	r.Post("/login",
		authHandlers.LoginUser,
	)
	r.Route("/user", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", userHandlers.GetUser)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", taskHandlers.GetTasks)
		r.Post("/", taskHandlers.CreateTask)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(taskMiddleware)
			r.Get("/", taskHandlers.GetTask)
			r.Delete("/", taskHandlers.DeleteTask)
			r.Patch("/", taskHandlers.UpdateTask)
		})
	})

	return r
}
