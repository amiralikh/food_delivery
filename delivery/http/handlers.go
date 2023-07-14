// app/delivery/http/handlers.go

package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"foodDelivery/domain"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
)

// UserHandler represents the HTTP handler for user operations.
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetUserByID handles the request to retrieve a user by ID.
func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	// Convert the userID to the appropriate data type, such as int64, and use it as needed.

	user, err := uh.userUseCase.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the user object into JSON.
	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

// CreateUser handles the request to create a new user.
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.userUseCase.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := []byte(`{"message": "User created successfully"}`)
	_, _ = w.Write(response)
}

// UpdateUser handles the request to update an existing user.
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL path parameters.
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	// Parse the user ID string to int64.
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decode the request body into a user object.
	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the user ID in the user object.
	user.ID = userID

	// Call the use case method to update the user.
	err = uh.userUseCase.UpdateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "User updated successfully"}`)
	_, _ = w.Write(response)
}

// DeleteUser handles the request to delete a user.
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	// Convert the userID to the appropriate data type, such as int64, and use it as needed.

	err := uh.userUseCase.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "User deleted successfully"}`)
	_, _ = w.Write(response)
}
