package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"

	intPkg "foodDelivery/delivery/http"
	"foodDelivery/migrations"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new database connection.
	// Create a new database connection.
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/food_delivery?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create the users table.
	err = migrations.CreateUsersTable(db)
	err = migrations.CreateCategoriesTable(db)
	err = migrations.CreateSuppliersTable(db)
	err = migrations.CreateFoodsTable(db)
	err = migrations.CreateGalleryTable(
		db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create an instance of the repository.
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	supplierRepository := repository.NewSupplierRepository(db)
	foodRepository := repository.NewFoodRepository(db)
	galleryRepository := repository.NewGalleryRepository(db)

	// Create an instance of the use case, passing in the UserRepository interface.
	userUseCase := usecase.NewUserUseCase(userRepository)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	supplierUseCase := usecase.NewSupplierUseCase(supplierRepository)
	foodUseCase := usecase.NewFoodUseCase(foodRepository, categoryRepository, supplierRepository, galleryRepository)

	// Create an instance of the user handler, passing in the UserUseCase interface.
	userHandler := intPkg.NewUserHandler(userUseCase)
	categoryHandler := intPkg.NewCategoryHandler(categoryUseCase)
	supplierHandler := intPkg.NewSupplierHandler(supplierUseCase)
	foodHandler := intPkg.NewFoodHandler(foodUseCase)

	// Create a new router.
	router := mux.NewRouter()

	// users API
	router.HandleFunc("/api/users/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// categories Api
	router.HandleFunc("/api/categories/{id}", categoryHandler.GetCategoryByID).Methods("GET")
	router.HandleFunc("/api/categories", categoryHandler.GetAllCategories).Methods("GET")
	router.HandleFunc("/api/categories", categoryHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/api/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/api/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// suppliers API
	router.HandleFunc("/api/suppliers/{id}", supplierHandler.GetSupplierByID).Methods("GET")
	router.HandleFunc("/api/suppliers", supplierHandler.GetAllSuppliers).Methods("")
	router.HandleFunc("/api/suppliers", supplierHandler.CreateSupplier).Methods("POST")
	router.HandleFunc("/api/suppliers/{id}", supplierHandler.UpdateSupplier).Methods("PUT")
	router.HandleFunc("/api/suppliers/{id}", supplierHandler.DeleteSupplier).Methods("DELETE")

	// foods API
	router.HandleFunc("/api/foods", foodHandler.GetAllFoodsWithImages).Methods("GET")
	router.HandleFunc("/api/foods/{id}", foodHandler.GetFoodByID).Methods("GET")
	router.HandleFunc("/api/foods", foodHandler.CreateFood).Methods("POST")
	router.HandleFunc("/api/foods/{id}", foodHandler.UpdateFood).Methods("PUT")
	router.HandleFunc("/api/foods/{id}", foodHandler.DeleteFood).Methods("DELETE")

	// Start the HTTP server.
	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
