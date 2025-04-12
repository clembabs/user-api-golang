package handlers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/clembabs/user-api/models"
)

const dataFile = "data/users.json"

func LoadUsers() []models.User {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return []models.User{}
	}
	defer file.Close()

	bytes, _ := os.ReadFile(dataFile)

	var users []models.User
	json.Unmarshal(bytes, &users)
	return users
}

func SaveUsers(users []models.User) {
	bytes, err := json.MarshalIndent(users, "", "  ") // convert Go users into JSON
	if err != nil {
		log.Println("Error marshaling:", err)
		return
	}

	err = os.WriteFile(dataFile, bytes, 0644) // write the JSON to file
	if err != nil {
		log.Println("Error writing file:", err)
	}
}
