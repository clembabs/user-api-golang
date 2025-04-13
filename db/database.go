package db

import (
	"database/sql"
	"log"

	"github.com/clembabs/user-api/repositories"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	// Open SQLite database connection (this could be a memory database or file-based)
	DB, err = sql.Open("sqlite", "data/users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize the repository
	repo := repositories.NewSQLiteUserRepository(DB)
	// Ensure the database schema is initialized (create tables if not already)
	if err := repo.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	return nil
}
