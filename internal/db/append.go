package db

import (
	"context"
	"errors"
	"fmt"
	"modules/internal/entities"
	"modules/utils"

	"github.com/jackc/pgconn"
)

// AddGood получает данные по товару и добавляет их на склад
func (d *db) AddGood(ctx context.Context, code string, stockId string, value int64, dynamic bool) error { // Параметр dynamic позволяет переносить товар с недоступного склада на доступный
	if code == "" || stockId == "" || value == 0 {
		d.logger.Fatal("code, id stock or value is empty")
		return errors.New("code, id stock or value cannot be empty")
	}

	var q string

	q = `select available from stocks where id::text = $1`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	t, err := d.client.Begin(ctx)
	if err != nil {
		d.logger.Fatal("couldn't open transaction")
		return err
	}

	var s entities.Stock

	if errQ = t.QueryRow(ctx, q, stockId).Scan(&s.Available); errQ != nil { // Получаем статус доступности склада
		d.logger.Fatalf("func AddGood query for search stock error %s", errQ)
		return errQ
	}

	if s.Available { // Если склад доступен, обновляем value в таблице goods
		q = `update goods set value = (select value from goods where code::text = $2::text) + $1 where code::text = $2::text and stock_id = $3`
		d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

		_, errQ = t.Exec(ctx, q, value, code, stockId)
	} else {
		switch dynamic {
		case true: // Если был указан параметр dynamic => находим и возвращаем любой доступный склад
			q = `select s.id from stocks s
							inner join goods g on g.code::text = $1 and g.stock_id = $2
						where s.id != $2 and s.available limit 1`

			if errQ = t.QueryRow(ctx, q, code, stockId).Scan(&s.ID); errQ != nil {
				d.logger.Fatalf("func AddGood query for get another stock (dynamic case) error %s", errQ)

				return errQ
			}

			q = `update goods set value = (select value from goods where code::text = $2::text) + $1 where stock_id = $3` // TODO: Пофиксить запрос

			_, errQ = t.Exec(ctx, q, value, code, s.ID) // Если запрос прошел все проверки и новый склад найден => Делаем обновление количества
		case false:
			return errors.New("failed to add goods to the stock")
		}
	}

	if errors.Is(errQ, pgErr) {
		pgErr = errQ.(*pgconn.PgError)
		newErr := fmt.Errorf(
			fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
		)
		d.logger.Error(newErr)

		return newErr
	}

	if err := t.Commit(ctx); err != nil {
		d.logger.Errorf("failed to commit transaction: %s", err)
		return err
	}
	return nil
}
