package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"foodDelivery/domain"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
	"time"
)

type AuthHandler struct {
	userUseCase usecase.UserUseCase
}

var (
	ErrNotFound = errors.New("not found")
)

func NewAuthHandler(userUseCase usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
	}
}

// Add this function to check if the provided password is valid.
func isPasswordValid(inputPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "email is not valid", http.StatusForbidden)
		return
	}

	_, err = ah.userUseCase.GetUserByEmail(user.Email)
	if err != nil {
		if !IsNotFoundError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash the user's password before storing it in the database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Create the user in the database (assuming UserRepository has the necessary method).
	err = ah.userUseCase.RegisterUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message or JWT token, if required.
	response := []byte(`{"message": "User registered successfully"}`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (uh *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.userUseCase.GetUserByEmail(loginRequest.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !isPasswordValid(loginRequest.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}{
		AccessToken:  "Bearer " + accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(time.Minute * 15).Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateTokens(userID int64) (string, string, error) {
	var secretKey = []byte(os.Getenv("SECRETKEY"))
	var accessTokenExpirationTime = time.Now().Add(time.Minute * 15)
	var refreshTokenExpirationTime = time.Now().Add(time.Hour * 24)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    accessTokenExpirationTime.Unix(),
	})

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    refreshTokenExpirationTime.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (int64, error) {
	// Define the JWT secret key used to sign and verify the tokens.
	// Replace "your-secret-key" with a strong secret key.
	secretKey := []byte(os.Getenv("SECRETKEY"))

	// Parse the JWT token using the secret key.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used in the token.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Check if the token is valid.
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract the user ID from the token's claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	// Get the user ID claim from the token as a float64.
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	// Convert the user ID to int64.
	userID := int64(userIDFloat)

	return userID, nil
}

func isValidEmail(email string) bool {
	// Regular expression for validating email format
	// This regular expression is a simplified version and may not cover all cases.
	// You can find more comprehensive email validation regex patterns online.
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	regexpEmail := regexp.MustCompile(emailRegex)

	// Use the MatchString function to check if the email matches the pattern
	return regexpEmail.MatchString(email)
}

func IsNotFoundError(err error) bool {
	return err == ErrNotFound || err == repository.ErrUserNotFound
}
