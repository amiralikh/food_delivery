package main

import (
	"log"
	stdhttp "net/http"

	"foodDelivery/config"
	"foodDelivery/delivery/http"
	"foodDelivery/repository"
	"foodDelivery/usecase"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new database connection.
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create an instance of the repository.
	userRepository := repository.NewUserRepository(db)

	// Create an instance of the use case, passing in the UserRepository interface.
	userUseCase := usecase.NewUserUseCase(userRepository)

	// Create an instance of the user handler, passing in the UserUseCase interface.
	userHandler := http.NewUserHandler(userUseCase)

	// Create a new router.
	router := mux.NewRouter()

	// Register the user routes.
	router.HandleFunc("/api/users/{id}", userHandler.GetUserByID).Methods("GET")
	router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Start the HTTP server.
	log.Println("Server started on port 8080")
	err = stdhttp.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
