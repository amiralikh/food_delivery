package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("M@ggie&&mIK@")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
