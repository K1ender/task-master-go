package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/k1ender/task-master-go/docs"
	"github.com/k1ender/task-master-go/internal/config"
	"github.com/k1ender/task-master-go/internal/handlers"
	"github.com/k1ender/task-master-go/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"gorm.io/gorm"
)

func New(db *gorm.DB, config *config.Config) *chi.Mux {
	r := chi.NewRouter()

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", config.HttpServer.Port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "Task Master API"
	docs.SwaggerInfo.Version = "1.0"

	validator := validator.New(validator.WithRequiredStructEnabled())
	authHandlers := handlers.NewAuthHandler(db, validator, config)
	taskHandlers := handlers.NewTaskHandler(db, validator, config)

	authMiddleware := middleware.Auth(db, config.JWT.Secret)
	taskMiddleware := middleware.TaskMiddleware(db)

	r.Post("/register",
		authHandlers.RegisterUser,
	)
	r.Post("/login",
		authHandlers.LoginUser,
	)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf(
				"http://localhost:%s/swagger/doc.json",
				config.HttpServer.Port,
			),
		),
	))

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
