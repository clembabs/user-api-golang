package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clembabs/user-api/models"
	"github.com/clembabs/user-api/repositories"
	"github.com/clembabs/user-api/response"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type UserHandler struct {
	Repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.Repo.GetAll()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ApiResponseWrapper{
			Message: "Failed to fetch users",
			Error:   true,
		})
		return
	}
	response.WriteJSON(w, http.StatusOK, response.ApiResponseWrapper{
		Message: "Users retrieved successfully",
		Error:   false,
		Data:    users,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := (mux.Vars(r)["id"])
	user, err := h.Repo.GetByID(id)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: "Failed to fetch user",
			Error:   true,
		})

		return
	}
	response.WriteJSON(w, http.StatusOK, response.ApiResponseWrapper{
		Message: "User retrieved successfully",
		Error:   false,
		Data:    user,
	})

}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	err := validate.Struct(user)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}

	if err := h.Repo.Create(&user); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}
	response.WriteJSON(w, http.StatusOK, response.ApiResponseWrapper{
		Message: "User Created successfully",
		Error:   false,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	err := validate.Struct(user)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}

	// Call the repository to update the user
	if err := h.Repo.Update(&user); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ApiResponseWrapper{
		Message: "User Updated successfully",
		Error:   false,
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL params
	id := mux.Vars(r)["id"]

	// Call the repository to delete the user
	if err := h.Repo.Delete(id); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}

	// Send no content response
	response.WriteJSON(w, http.StatusOK, response.ApiResponseWrapper{
		Message: "User Deleted successfully",
		Error:   false,
	})
}
