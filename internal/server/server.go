package server

import (
	"context"
	"modules/internal/db"
	"modules/pkg/logger"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func NewService(ctx context.Context, strg db.Storage, log *logger.Logger) *Service {
	return &Service{
		ctx: ctx,
		db:  strg,
		log: log,
	}
}

func RunJRPC(ctx context.Context, strg db.Storage, log *logger.Logger) {
	s := rpc.NewServer()
	log.Info("run server")

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(ctx, strg, log), "")
	log.Info("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":8080", nil)
}
