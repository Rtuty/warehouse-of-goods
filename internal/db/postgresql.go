package db

import (
	"context"
	"modules/internal/entities"
	"modules/pkg/dbclient"
	"modules/pkg/logger"

	"github.com/jackc/pgconn"
)

type Storage interface {
	CreateNewGood(ctx context.Context, g entities.Good) error
	CreateNewStock(ctx context.Context, s entities.Stock) error

	GetGoodsCountByStockId(ctx context.Context, stockId string, code string) (int64, error)

	GetAllGoods(ctx context.Context) ([]entities.Good, error)
	GetGoodByCode(ctx context.Context, code string) (entities.Good, error)
	AddGood(ctx context.Context, code string, stockId string, value int64, dynamic bool) error

	ReservationGood(ctx context.Context, code string, stockId string, value int64) error
	CancelGoodReservation(ctx context.Context, resId string) error
}

type db struct {
	client dbclient.Client
	logger *logger.Logger
}

func NewRepository(client dbclient.Client, logger *logger.Logger) Storage {
	return &db{
		client: client,
		logger: logger,
	}
}

var pgErr *pgconn.PgError
var errQ error
