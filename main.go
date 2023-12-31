package main

import (
	"database/sql"
	intPkg "foodDelivery/delivery/http"
	"foodDelivery/migrations"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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
	err = migrations.CreateGalleryTable(db)
	err = migrations.CreateAddressesTable(db)
	err = migrations.CreateOrdersTable(db)
	err = migrations.CreateOrderItemsTable(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create an instance of the repository.
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	supplierRepository := repository.NewSupplierRepository(db)
	foodRepository := repository.NewFoodRepository(db)
	galleryRepository := repository.NewGalleryRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	addressRepository := repository.NewAddressRepository(db)

	// Create an instance of the use case, passing in the UserRepository interface.
	userUseCase := usecase.NewUserUseCase(userRepository)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository)
	supplierUseCase := usecase.NewSupplierUseCase(supplierRepository)
	foodUseCase := usecase.NewFoodUseCase(foodRepository, categoryRepository, supplierRepository, galleryRepository)
	orderUseCase := usecase.NewOrderUseCase(orderRepository)
	addressUseCase := usecase.NewAddressUseCase(addressRepository)

	// Create an instance of the user handler, passing in the UserUseCase interface.
	userHandler := intPkg.NewUserHandler(userUseCase)
	categoryHandler := intPkg.NewCategoryHandler(categoryUseCase)
	supplierHandler := intPkg.NewSupplierHandler(supplierUseCase, categoryUseCase, foodUseCase)
	foodHandler := intPkg.NewFoodHandler(foodUseCase)
	authHandler := intPkg.NewAuthHandler(userUseCase)
	orderHandler := intPkg.NewOrderHandler(orderUseCase)
	addressHandler := intPkg.NewAddressHandler(addressUseCase)

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
	router.HandleFunc("/api/suppliers", supplierHandler.GetAllSuppliers).Methods("GET")
	router.HandleFunc("/api/suppliers", supplierHandler.CreateSupplier).Methods("POST")
	router.HandleFunc("/api/suppliers/{id}", supplierHandler.UpdateSupplier).Methods("PUT")
	router.HandleFunc("/api/suppliers/{id}", supplierHandler.DeleteSupplier).Methods("DELETE")
	router.HandleFunc("/api/suppliers/{id}/categories", supplierHandler.GetSupplierCategories).Methods("GET")
	router.HandleFunc("/api/supplier/{cat_id}/food-list/{supplier_id}", supplierHandler.GetFoodsByCategoryAndSupplier).Methods("GET")

	// foods API
	router.HandleFunc("/api/foods", foodHandler.GetAllFoodsWithImages).Methods("GET")
	router.HandleFunc("/api/foods", foodHandler.CreateFood).Methods("POST")
	router.HandleFunc("/api/foods/{id}", foodHandler.GetFoodByID).Methods("GET")
	router.HandleFunc("/api/foods/{id}", foodHandler.UpdateFood).Methods("PUT")
	router.HandleFunc("/api/foods/{id}", foodHandler.DeleteFood).Methods("DELETE")

	// auth
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")

	// orders
	router.HandleFunc("/api/orders", orderHandler.SubmitOrder).Methods("POST")
	router.HandleFunc("/api/orders", orderHandler.GetUserOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", orderHandler.GetOrderWithItems).Methods("GET")

	// addresses
	router.HandleFunc("/api/addresses", addressHandler.GetUsersAddresses).Methods("GET")
	router.HandleFunc("/api/addresses/{id}", addressHandler.GetAddressByID).Methods("GET")
	router.HandleFunc("/api/addresses/{id}", addressHandler.DeleteAddress).Methods("DELETE")
	router.HandleFunc("/api/addresses/{id}", addressHandler.UpdateAddress).Methods("PUT")
	router.HandleFunc("/api/addresses", addressHandler.CreateAddress).Methods("POST")

	// fix cross error
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Start the HTTP server.
	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))

	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
