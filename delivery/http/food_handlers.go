package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"foodDelivery/domain"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type FoodHandler struct {
	foodUseCase usecase.FoodUseCase
}

func NewFoodHandler(foodUseCase usecase.FoodUseCase) *FoodHandler {
	return &FoodHandler{
		foodUseCase: foodUseCase,
	}
}

func (fh *FoodHandler) GetFoodByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foodIDStr := vars["id"]
	foodID, err := strconv.ParseInt(foodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid food ID", http.StatusBadRequest)
		return
	}

	food, err := fh.foodUseCase.GetFoodByID(foodID)
	if err != nil {
		if errors.Is(err, usecase.ErrFoodNotFound) {
			http.Error(w, "Food not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get food: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(food)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to serialize food: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (fh *FoodHandler) CreateFood(w http.ResponseWriter, r *http.Request) {
	var food domain.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	err = fh.foodUseCase.CreateFood(&food)
	if err != nil {
		if errors.Is(err, usecase.ErrFoodNameRequired) || errors.Is(err, usecase.ErrCategoryRequired) || errors.Is(err, usecase.ErrSupplierRequired) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, usecase.ErrCategoryNotFound) || errors.Is(err, usecase.ErrSupplierNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create food: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := []byte(`{"message": "Food created successfully"}`)
	_, _ = w.Write(response)
}

func (fh *FoodHandler) UpdateFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foodIDStr := vars["id"]
	foodID, err := strconv.ParseInt(foodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid food ID", http.StatusBadRequest)
		return
	}

	var food domain.Food
	err = json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	food.ID = foodID

	err = fh.foodUseCase.UpdateFood(&food)
	if err != nil {
		if errors.Is(err, usecase.ErrFoodNotFound) {
			http.Error(w, "Food not found", http.StatusNotFound)
			return
		} else if errors.Is(err, usecase.ErrSupplierNotFound) {
			http.Error(w, "Supplier not found", http.StatusBadRequest)
			return
		} else if errors.Is(err, usecase.ErrCategoryNotFound) {
			http.Error(w, "Category not found", http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle gallery update.
	err = fh.foodUseCase.SyncGallery(food.ID, food.Gallery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Food and Gallery updated successfully"}`)
	_, _ = w.Write(response)
}

func (fh *FoodHandler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foodIDStr := vars["id"]
	foodID, err := strconv.ParseInt(foodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid food ID", http.StatusBadRequest)
		return
	}

	err = fh.foodUseCase.DeleteFood(foodID)
	if err != nil {
		if errors.Is(err, usecase.ErrFoodNotFound) {
			http.Error(w, "Food not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Food deleted successfully"}`)
	_, _ = w.Write(response)
}
