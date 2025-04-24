package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/clembabs/user-api/models"
	"github.com/clembabs/user-api/repositories"
	"github.com/clembabs/user-api/response"
	"github.com/clembabs/user-api/utils"
)

type AuthHandler struct {
	Repo repositories.UserRepository
}

func NewAuthHandler(repo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

type AuthResponse struct {
	User         models.User `json:"user"`
	Token        string      `json:"token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: "Invalid request payload",
			Error:   true,
		})
		return
	}
	hashed, _ := utils.HashPassword(user.Password)
	user.Password = hashed

	err := h.Repo.Create(&user)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: err.Error(),
			Error:   true,
		})
		return
	}
	token, _ := utils.GenerateJWT(user.ID)
	response.WriteJSON(w, http.StatusCreated, response.ApiResponseWrapper{
		Message: "User created",
		Error:   false,
		Data: AuthResponse{
			User:  user,
			Token: token,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ApiResponseWrapper{
			Message: "Invalid request payload",
			Error:   true,
		})
		return
	}

	user, err := h.Repo.GetByEmail(creds.Email)

	if err != nil || user == nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.ApiResponseWrapper{
			Message: "User does not exist",
			Error:   true,
		})
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		response.WriteJSON(w, http.StatusUnauthorized, response.ApiResponseWrapper{
			Message: "Invalid credentials",
			Error:   true,
		})
		return
	}
	// Generate a JWT token
	token, _ := utils.GenerateJWT(user.ID)
	// Generate a refresh token
	// Note: In a real-world application, you should store the refresh token securely
	// and implement a mechanism to refresh the access token using the refresh token.
	// For simplicity, we're just generating a new refresh token here.
	//Your code works and matches your comment: you're just generating a refresh token and validating it.

	// For learning/dev, it's fine.

	// In production, you'd need storage, rotation, and revocation mechanisms to protect against abuse.

	refreshToken, _ := utils.GenerateRefreshToken(user.ID)
	response.WriteJSON(w, http.StatusCreated, response.ApiResponseWrapper{
		Message: "Login successful",
		Error:   false,
		Data: AuthResponse{
			User:         *user, // dereference the pointer
			Token:        token,
			RefreshToken: refreshToken,
		},
	})

}
