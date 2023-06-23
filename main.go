package main

import (
	"context"
	"log"
	"modules/internal/db"
	"modules/internal/server"
	"modules/pkg/dbclient"
	"modules/pkg/logger"

	"github.com/joho/godotenv"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	dbclient.GetConnection()    // Получаем данные для соединением с БД, записывая в сущность PstgCon
	ctx := context.Background() // Создаем контекст
	log := logger.GetLogger()   // Создаем логгер

	cl, err := dbclient.NewClient(ctx, 5, dbclient.PstgCon, log) // Получаем соединение с PostgreSQL
	if err != nil {
		log.Fatalf("failed to connect PostgreSQL error: %s", err)
	}

	server.RunJRPC(ctx, ":8080", db.NewRepository(cl, log), log) // Запуск JRPC сервера
}
