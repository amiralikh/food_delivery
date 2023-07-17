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
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create an instance of the repository.
	userRepository := repository.NewUserRepository(db)

	// Create an instance of the use case, passing in the UserRepository interface.
	userUseCase := usecase.NewUserUseCase(userRepository)

	// Create an instance of the user handler, passing in the UserUseCase interface.
	userHandler := intPkg.NewUserHandler(userUseCase)

	// Create a new router.
	router := mux.NewRouter()

	// users API
	router.HandleFunc("/api/users/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// categories Api

	// Start the HTTP server.
	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
