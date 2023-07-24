package middleware

import (
	"context"
	authH "foodDelivery/delivery/http"
	"foodDelivery/usecase"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the "Authorization" header value from the request.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the header starts with "Bearer " to validate the scheme.
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header.
		token := authHeader[len("Bearer "):]

		// Your token validation logic goes here.
		// For example, you can use the ValidateToken function to validate the token and extract the user ID.
		userID, err := authH.ValidateToken(token)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set the user ID as a request context value to be used by the handler.
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	}
}

func AuthenticateMiddleware(userUseCase usecase.UserUseCase) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the JWT token from the request header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			// Check if the token is in the correct format.
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// Get the token string.
			tokenString := splitToken[1]

			// Validate the token and get the user ID from it.
			userID, err := authH.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Get the user from the database based on the user ID.
			user, err := userUseCase.GetUserByID(userID)
			if err != nil {
				if authH.IsNotFoundError(err) {
					http.Error(w, "User not found", http.StatusUnauthorized)
				} else {
					http.Error(w, "Error retrieving user", http.StatusInternalServerError)
				}
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			// Call the next handler in the chain with the updated context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
