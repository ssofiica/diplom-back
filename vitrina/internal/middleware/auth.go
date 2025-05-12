package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"back/vitrina/internal/delivery"
	"back/vitrina/internal/entity"
	"back/vitrina/utils/response"

	"github.com/golang-jwt/jwt/v5"
)

var userKey string = "user"

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	handler := "JWTMiddleware"
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WithError(w, 401, handler, delivery.ErrDefault401)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Неверный формат токена")
			response.WithError(w, 401, handler, delivery.ErrDefault401)
			return
		}
		tokenString := parts[1]

		claims := &entity.Claims{}
		t, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
			if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("недопустимый метод подписи")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			response.WithError(w, 401, handler, delivery.ErrDefault401)
			return
		}
		if !t.Valid {
			fmt.Println("Недействительный токен")
			response.WithError(w, 401, handler, delivery.ErrDefault401)
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			response.WithError(w, 401, handler, delivery.ErrDefault401)
			return
		}

		user := claims.User
		ctx := context.WithValue(r.Context(), userKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
