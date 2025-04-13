package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clembabs/user-api/db"
	"github.com/clembabs/user-api/handlers"
	"github.com/clembabs/user-api/repositories"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}
	// Ensure DB connection is closed when the program exits
	defer db.DB.Close()

	// Initialize the user repository with the database connection
	repo := repositories.NewSQLiteUserRepository(db.DB)

	// Create the handler with the repository
	handler := handlers.NewUserHandler(repo)

	r := mux.NewRouter()
	r.HandleFunc("/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	fmt.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
