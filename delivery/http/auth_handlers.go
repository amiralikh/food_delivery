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
	"strconv"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	err = ah.userUseCase.RegisterUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

const (
	accessTokenDuration  = time.Minute * 15
	refreshTokenDuration = time.Hour * 24 * 7
)

func generateTokens(userID int64) (string, string, error) {
	secretKey := []byte(os.Getenv("M@ggie&&mIK@"))

	accessClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
		Issuer:    "your-app",
		Subject:   strconv.FormatInt(userID, 10),
	}

	refreshClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
		Issuer:    "foodDelivery",
		Subject:   strconv.FormatInt(userID, 10),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func ValidateToken(tokenStr string) (bool, error) {
	secretKey := []byte(os.Getenv("M@ggie&&mIK@"))

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		fmt.Println("User ID:", claims.Subject)
		return true, nil
	}

	return false, nil
}

func isValidEmail(email string) bool {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regexpEmail := regexp.MustCompile(emailRegex)

	return regexpEmail.MatchString(email)
}

func IsNotFoundError(err error) bool {
	return err == ErrNotFound || err == repository.ErrUserNotFound
}
