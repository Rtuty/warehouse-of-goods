package server

import (
	"context"
	"modules/internal/db"
	"modules/pkg/logger"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func RunJRPC(ctx context.Context, port string, strg db.Storage, log *logger.Logger) {
	s := rpc.NewServer() // Создание нового сервера JSON-RPC

	if err := s.RegisterName("goods", &Handler{ctx: ctx, db: strg, log: log}); err != nil { // Регистрация сервиса склада
		log.Fatal("Error registering service: ", err)
	}

	// Слушаем порт и обслуживаем запросы JSON-RPC
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("error listening: ", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err)
		}

		go jsonrpc.ServeConn(conn)
	}
}
