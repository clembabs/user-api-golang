package models

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
