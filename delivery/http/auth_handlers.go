package http

import (
	"encoding/json"
	"errors"
	"foodDelivery/domain"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
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

	// Validate the password
	if !isPasswordValid(loginRequest.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token as a response
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateToken(userID int64) (string, error) {
	var secretKey = os.Getenv("SECRETKEY")
	var expirationTime = os.Getenv("EXPIRATIONTIME")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expirationTime,
	})

	// Sign the token with the secret key.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
