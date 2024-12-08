package middleware

import (
	"day_03/internal/config"
	"day_03/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func AuthMiddleware(config *config.Config, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.UserName == "" {
			http.Error(w, "Invalid token claims: missing username", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
