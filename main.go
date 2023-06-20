package main

import (
	"log"
	"modules/internal/db"

	"github.com/joho/godotenv"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	db.GetConnection()
	db, err := db.NewDbClient()
	if err != nil {
		panic(err)
	}
}
