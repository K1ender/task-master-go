package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"gorm.io/gorm"
)

type TaskKeyType string

const TaskKey TaskKeyType = "task"

func TaskMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetAuthUserFromContext(r.Context())
			taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
			if err != nil {
				response.BadRequest(w, "Bad Request")
				return
			}

			if taskID < 0 {
				response.BadRequest(w, "Bad Request")
				return
			}

			var task models.Task
			res := db.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task)

			if res.Error != nil {
				if res.Error == gorm.ErrRecordNotFound {
					response.NotFound(w, "Task not found")
					return
				}
				response.InternalServerError(w)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, TaskKey, &task)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTaskFromContext(ctx context.Context) *models.Task {
	return ctx.Value(TaskKey).(*models.Task)
}
