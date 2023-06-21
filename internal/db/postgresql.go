package db

import (
	"context"
	"errors"
	"fmt"
	"modules/internal/entities"
	"modules/pkg/dbclient"
	"modules/pkg/logger"
	"modules/utils"
	"strconv"
	"sync"

	"github.com/jackc/pgconn"
)

type Storage interface {
	CreateNewGood(ctx context.Context, g entities.Good) error
	CreateNewStock(ctx context.Context, s entities.Stock) error

	AddGood(ctx context.Context, code string, stockId string, value int) error

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

var pgErr *pgconn.PgError
var errQ error
var mu sync.Mutex

// CreateNewGood cоздает новый типа товара в базе данных
func (d *db) CreateNewGood(ctx context.Context, g entities.Good) error {
	mu.Lock() // Для избежания "гонок данных" и коллизий используем mutex

	exist, err := checkDbDublicate(strconv.Itoa(int(g.Code)), "goods", ctx, d.client) // Проверяем код нового продукта на дубликаты
	if err != nil {
		return fmt.Errorf("good being created is already in stock, error: %s", err)
	}

	if exist { // если код товара найден в БД - логируем и возвращаем ошибку
		exErr := "specified new product code exists in stock"
		d.logger.Error(exErr)
		mu.Unlock()
		return errors.New(exErr)
	}

	q := `insert into goods (name, code, size, value) values ($1, $2, $3, $4)` // Исполняем запрос по добавлению нового товара + проверяем на всевозможные ошибки и логируем
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	_, errQ = d.client.Exec(ctx, q, g.Name, g.Code, g.Size, g.Value)

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(
			fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
		)
		d.logger.Error(newErr)
		mu.Unlock()
		return newErr
	}
	mu.Unlock()
	d.logger.Trace("new good has been successfully added to the database")

	return nil
}

// CreateNewStock cоздает новый склад в базе данных. Функция аналогична CreateNewGood
func (d *db) CreateNewStock(ctx context.Context, s entities.Stock) error {
	mu.Lock()

	exist, err := checkDbDublicate(s.Name, "stocks", ctx, d.client)
	if err != nil {
		return fmt.Errorf("good being created is already in stock, error: %s", err)
	}

	if exist {
		exErr := "specified name of the stock being created already exists in the database"
		d.logger.Error(exErr)
		mu.Unlock()
		return errors.New(exErr)
	}

	q := `insert into stocks (name, available) values ($1, $2)`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	_, errQ = d.client.Exec(ctx, q, s.Name, s.Available)

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(
			fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
		)
		d.logger.Error(newErr)
		mu.Unlock()
		return newErr
	}
	mu.Unlock()
	d.logger.Trace("new stock has been successfully added to the database")

	return nil
}

// AddGood получает данные по товару и добавляет их на склад
func (d *db) AddGood(ctx context.Context, code string, stockId string, value int) error {
	return nil
}

func (d *db) ReserveGood(ctx context.Context, code int, stockId int, value int) error {
	return nil
}

func (d *db) CancelGoodReserve(ctx context.Context, code int, stockId int, value int) error {
	return nil
}

// Дополнительный функционал, реализующийся внутри пакета

// checkDbDublicate проверяет на дубликаты в БД, при создании нового товара или скалада
func checkDbDublicate(arg string, dbName string, ctx context.Context, c dbclient.Client) (bool, error) {
	var query string
	var exist bool

	if arg == "" {
		return false, errors.New("argument is empty")
	}

	switch dbName {
	case "goods":
		query = "select case when (select * from goods where code::text = '$1'::text) is not null then true else false end"
	case "stocks":
		query = "select case when (select * from stocks where name::text = '$1'::text) is not null then true else false end"
	default:
		return false, errors.New("database was not specified in the arguments of the function")
	}

	if err := c.QueryRow(ctx, query, arg).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}
