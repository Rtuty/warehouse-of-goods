package db

import (
	"context"
	"errors"
	"fmt"
	"modules/internal/entities"
	"modules/utils"
	"strconv"

	"github.com/jackc/pgconn"
)

// CreateNewGood cоздает новый типа товара в базе данных
func (d *db) CreateNewGood(ctx context.Context, g entities.Good) error {
	t, err := d.client.Begin(ctx) // Открываем транзакцию
	if err != nil {
		return fmt.Errorf("error when creating a transaction: %s", err)
	}

	exist, err := checkDbDublicate(strconv.Itoa(int(g.Code)), "goods", ctx, t) // Проверяем код нового продукта на дубликаты
	if err != nil {
		t.Rollback(ctx)
		return fmt.Errorf("good being created is already in stock, error: %s", err)
	}

	if exist { // если код товара найден в БД - логируем и возвращаем ошибку
		exErr := "specified new product code exists in stock"
		t.Rollback(ctx)
		d.logger.Error(exErr)
		return errors.New(exErr)
	}

	q := `insert into goods (name, size, value) values ($1, $2, $3)` // Создаем запрос по добавлению нового товара + проверяем на всевозможные ошибки и логируем
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	_, errQ = t.Exec(ctx, q, g.Name, g.Size, g.Value) // FAQ Метод Exec использовать для исполнения запросов, которые не возвращают данных update|delete|insert. Метод Query использовать для исполнения и возврата (select)

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(
			fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
		)
		t.Rollback(ctx)
		d.logger.Error(newErr)
		return newErr
	}

	if err := t.Commit(ctx); err != nil {
		d.logger.Errorf("failed to commit transaction: %s", err)
		return err
	}

	d.logger.Trace("new good has been successfully added to the database")

	return nil
}

// CreateNewStock cоздает новый склад в базе данных. Функция аналогична CreateNewGood
func (d *db) CreateNewStock(ctx context.Context, s entities.Stock) error {
	t, err := d.client.Begin(ctx) // Открываем транзакцию
	if err != nil {
		return fmt.Errorf("error when creating a transaction: %s", err)
	}

	exist, err := checkDbDublicate(s.Name, "stocks", ctx, t)
	if err != nil {
		t.Rollback(ctx)
		return fmt.Errorf("good being created is already in stock, error: %s", err)
	}

	if exist {
		exErr := "specified name of the stock being created already exists in the database"
		t.Rollback(ctx)
		d.logger.Error(exErr)
		return errors.New(exErr)
	}

	q := `insert into stocks (name, available) values ($1, $2)`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	_, errQ = t.Exec(ctx, q, s.Name, s.Available)

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(
			fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
		)
		t.Rollback(ctx)
		d.logger.Error(newErr)
		return newErr
	}

	if err := t.Commit(ctx); err != nil {
		d.logger.Errorf("failed to commit transaction: %s", err)
		return err
	}

	d.logger.Trace("new stock has been successfully added to the database")

	return nil
}
