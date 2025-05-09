package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/clembabs/user-api/db"
	"github.com/clembabs/user-api/handlers"
	"github.com/clembabs/user-api/middlewares"
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
	authHandler := handlers.NewAuthHandler(repo)

	r := mux.NewRouter()

	//Auth
	r.HandleFunc("/auth/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	//User routes
	r.HandleFunc("/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	//TODO:
	r.Handle("/me", middlewares.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middlewares.UserIDKey).(string)
		json.NewEncoder(w).Encode(map[string]string{"id": userID})
	}))).Methods("GET")

	// wrap middlewares
	wrapped := middlewares.Logger(
		middlewares.Recover(
			middlewares.CORS(r),
		),
	)

	fmt.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", wrapped)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
