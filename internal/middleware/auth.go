package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k1ender/task-master-go/internal/models"
	"github.com/k1ender/task-master-go/internal/response"
	"gorm.io/gorm"
)

type AuthKey string

const AuthKeyUser AuthKey = "user"

func GetAuthUserFromContext(ctx context.Context) *models.User {
	return ctx.Value(AuthKeyUser).(*models.User)
}

func Auth(db *gorm.DB, jwtSecret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				response.Unauthorized(w, "Unauthorized")
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			validatedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

			if err != nil {
				response.Unauthorized(w, "Unauthorized")
				return
			}

			var userID int

			if claims, ok := validatedToken.Claims.(jwt.MapClaims); ok {
				userID = int(claims["user_id"].(uint))
			} else {
				response.Unauthorized(w, "Unauthorized")
				return
			}

			var user models.User

			if res := db.First(&user, userID); res.Error != nil {
				if res.Error == gorm.ErrRecordNotFound {
					response.Unauthorized(w, "Unauthorized")
					return
				}
				response.InternalServerError(w)
				return
			}

			ctx := context.WithValue(r.Context(), AuthKeyUser, &user)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
