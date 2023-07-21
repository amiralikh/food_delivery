package http

import (
	"encoding/json"
	"foodDelivery/domain"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

func (ch *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryIDStr := vars["id"]
	categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 64)

	category, err := ch.categoryUseCase.GetCategoryByID(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(category)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (ch *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	err := json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = ch.categoryUseCase.CreateCategory(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := []byte(`{"message": "Create category successfully"}`)
	_, _ = w.Write(response)
}

func (ch *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catIDStr := vars["id"]
	catID, _ := strconv.ParseInt(catIDStr, 10, 64)

	var category domain.Category

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category.ID = catID

	err = ch.categoryUseCase.UpdateCategory(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Category updated successfully"}`)
	_, _ = w.Write(response)
}

func (ch *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catIDStr := vars["id"]
	catID, _ := strconv.ParseInt(catIDStr, 10, 64)

	err := ch.categoryUseCase.DeleteCategory(catID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []byte(`{"message": "Category deleted successfully"}`)
	_, _ = w.Write(response)

}

func (ch *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := ch.categoryUseCase.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the categories into JSON.
	response, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
