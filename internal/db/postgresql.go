package db

import "modules/internal/goods"

type Storage interface {
	CreateNewGood(g goods.Product) error
	CreateWarehouse(warehouse goods.Stock) error

	AddGood(code string, stock string, value int) error

	ReserveGood(code int, stockId int, value int) error
	CancelGoodReserve(code int, stockId int, value int) error
}
