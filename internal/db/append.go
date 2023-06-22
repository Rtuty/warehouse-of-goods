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
func (d *db) AddGood(ctx context.Context, code string, stockId string, value int, dynamic bool) error { // Параметр dynamic позволяет переносить товар с недоступного склада на доступный
	if code == "" || stockId == "" || value == 0 {
		d.logger.Fatal("code, id stock or value is empty")
		return errors.New("code, id stock or value cannot be empty")
	}

	var q string

	q = `select available from stocks where id::text = $1`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	var s entities.Stock

	mu.Lock()
	if errQ = d.client.QueryRow(ctx, q, stockId).Scan(&s.Available); errQ != nil { // Получаем статус доступности склада
		d.logger.Fatalf("func AddGood query for search stock error %s", errQ)
		mu.Unlock()
		return errQ
	}

	if s.Available { // Если склад доступен, обновляем value в таблице goods
		q = `update goods set value = $1 where code::text = $2 and stock_id = $3`
		d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

		_, errQ = d.client.Exec(ctx, q, value, code, stockId)
	} else {
		switch dynamic {
		case true: // Если был указан параметр dynamic => находим и возвращаем любой доступный склад (у которого есть товар с таким же кодом)
			q = `select s.id from stocks s
							inner join goods g on g.code::text = $1 and g.stock_id = $2
						where s.id != $3 and s.available limit 1`

			if errQ = d.client.QueryRow(ctx, q, code, stockId, stockId).Scan(&s.ID); errQ != nil {
				d.logger.Fatalf("func AddGood query for get another stock (dynamic case) error %s", errQ)
				mu.Unlock()
				return errQ
			}

			q = `update goods set value = $1 where code::text = $2 and stock_id = $3`

			_, errQ = d.client.Exec(ctx, q, value, code, s.ID) // Если запрос прошел все проверки и новый склад найден => Делаем обновление количества
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
		mu.Unlock()
		return newErr
	}

	mu.Unlock()
	return nil
}
