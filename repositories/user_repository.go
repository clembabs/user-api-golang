package repositories

import (
	"database/sql"
	"fmt"

	"github.com/clembabs/user-api/models"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var validate = validator.New()

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
}

type SQLiteUserRepository struct {
	DB *sql.DB
}

func (r *SQLiteUserRepository) Init() error {
	_, err := r.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		// If there's an error, print it to help with debugging
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *SQLiteUserRepository) GetByID(id string) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) Create(user *models.User) error {
	// Validate required fields
	if err := validate.Struct(user); err != nil {
		return err // returns validation error, e.g., if Name is empty
	}
	user.ID = uuid.New().String() // generate a UUID string

	// Insert user into the table
	_, err := r.DB.Exec("INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)", user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteUserRepository) Update(user *models.User) error {
	_, err := r.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func (r *SQLiteUserRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
