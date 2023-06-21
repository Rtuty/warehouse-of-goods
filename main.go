package main

import (
	"fmt"
	"log"
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

type WarehouseService struct{} // todo json rpc server

func main() {
	// Создание нового сервера JSON-RPC
	server := rpc.NewServer()

	// Регистрация сервиса склада
	err := server.Register(&WarehouseService{})
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
