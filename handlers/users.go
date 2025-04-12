package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clembabs/user-api/models"
	"github.com/gorilla/mux"
)

var users = LoadUsers()

func GetUsers(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := (mux.Vars(r)["id"])
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.NotFound(w, r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	users = append(users, user)
	SaveUsers(users)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := (mux.Vars(r)["id"])
	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	for i, user := range users {
		if user.ID == id {
			users[i] = updatedUser
			SaveUsers(users)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	http.NotFound(w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := (mux.Vars(r)["id"])
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			SaveUsers(users)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
