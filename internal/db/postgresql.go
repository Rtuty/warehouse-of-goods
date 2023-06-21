package db

import (
	"context"
	"errors"
	"fmt"
	"modules/internal/goods"
	"modules/pkg/dbclient"
	"modules/pkg/logger"
	"strings"

	"github.com/jackc/pgconn"
)

type Storage interface {
	CreateNewGood(ctx context.Context, g goods.Product) error
	CreateWarehouse(ctx context.Context, warehouse goods.Stock) error

	AddGood(ctx context.Context, code string, stock string, value int) error

	ReserveGood(ctx context.Context, code int, stockId int, value int) error
	CancelGoodReserve(ctx context.Context, code int, stockId int, value int) error
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

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (d *db) CreateNewGood(ctx context.Context, g goods.Product) error {
	var pgErr *pgconn.PgError
	var errQ error

	// Для создания нового типа продукта, треубется проверить, существует ли данный продукт на складе
	var exist bool
	if err := d.client.QueryRow(ctx,
		`select case when (select * from goods where code = $1) is not null 
		then true
		else false end`, g.Code).Scan(&exist); err != nil {
		return err
	}

	// Если запись по коду найдена - логируем и возвращаем ошибку
	if exist {
		exErr := "error! Указанный код нового продукта существует на складе"
		d.logger.Error(exErr)
		return errors.New(exErr)
	}

	// Исполняем запрос по добавлению нового товара + проверяем на всевозможные ошибки и логируем
	q := `insert into goods (name, code, size, value) values ($1, $2, $3, $4)`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", FormatQuery(q)))

	_, errQ = d.client.Exec(ctx, q, g.Name, g.Code, g.Size, g.Value)

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		d.logger.Error(newErr)
		return newErr
	}

	d.logger.Trace("the new product has been successfully added to the database")
	return nil
}

func (d *db) CreateWarehouse(ctx context.Context, warehouse goods.Stock) error { return nil }

func (d *db) AddGood(ctx context.Context, code string, stock string, value int) error { return nil }

func (d *db) ReserveGood(ctx context.Context, code int, stockId int, value int) error { return nil }
func (d *db) CancelGoodReserve(ctx context.Context, code int, stockId int, value int) error {
	return nil
}
