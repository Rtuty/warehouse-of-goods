package db

import (
	"context"
	"errors"
	"fmt"
	"modules/utils"

	"github.com/jackc/pgconn"
)

// faq: Таблица резервации. Вычетаем количество из таблицы goods и добавляем его в res_cen (reservation center)
// Ммм, переписывать буду все функции.
// Открываем транзацию, чтобы обеспечить более безопасную работу с постгрей
func (d *db) ReservationGood(ctx context.Context, code string, stockId string, value int64) error {
	if code == "" || stockId == "" || value == 0 {
		return fmt.Errorf("result: code = %s, stock_id = %s, value = %d. must not be equal to code == '' or stock id == '' or value == 0", code, stockId, value)
	}

	var q string

	c, err := d.GetGoodsCountByStockId(ctx, stockId, code)
	if err != nil {
		d.logger.Fatalf("reservation good error. couldn't get value of goods: %s", err)
		return err
	}

	if c < value { // Проверяем, можно ли зарезервировать запрашиваемое количество товара
		d.logger.Fatal("is not possible to reserve a good because it is not in stock")
		return errors.New("cannot reserve 0 goods")
	}

	// Открываем транзакцию. Обновляем значения в goods и res_cen
	t, err := d.client.Begin(ctx)
	if err != nil {
		return err
	}

	mu.Lock()

	errChan := make(chan error, 1)

	// После исполнения запроса проверяем на ошибки (выписать отдельной функцией и вынести)
	go func(errs <-chan error) {
		for err := range errs {
			if errors.Is(err, pgErr) {
				pgErr = errQ.(*pgconn.PgError)
				newErr := fmt.Errorf(
					fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
						pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
				)
				d.logger.Error(newErr)

				t.Rollback(ctx) // Отменяем транзакцию, если есть ошибка
			}

			mu.Unlock()
		}
	}(errChan)

	// Вычетаем количество резервируемого товара из таблицы goods
	q = `update goods set value = (value - $1) where code::text = $2 and stock_id = $3`
	_, errQ = t.Exec(ctx, q, value, code, stockId)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	errChan <- errQ

	// Создаем новую строку в таблице res_cen
	q = `insert into res_cen (good_code, stock_id, value) values ($1, $2, $3)`
	_, errQ = t.Exec(ctx, q, code, stockId, value)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	errChan <- errQ

	close(errChan)
	mu.Unlock()

	t.Commit(ctx) // Фиксируем транзакцию, если все окей
	return nil
}

func (d *db) CancelGoodReservation(ctx context.Context, code string, stockId string, value int64) error {
	return nil
}
