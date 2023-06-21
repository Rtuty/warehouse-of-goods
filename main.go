package main

import (
	"fmt"
	"log"
	"modules/internal/goods"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/joho/godotenv"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Создание нового сервера JSON-RPC
	server := rpc.NewServer()

	// Регистрация сервиса склада
	err := server.Register(&goods.WarehouseService{})
	if err != nil {
		log.Fatal("Error registering service: ", err)
	}

	// Слушаем порт и обслуживаем запросы JSON-RPC
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err)
		}
		fmt.Println("check")
		go jsonrpc.ServeConn(conn)
	}
}
