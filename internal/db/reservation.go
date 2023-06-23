package db

import (
	"context"
	"errors"
	"fmt"
	"modules/utils"

	"github.com/jackc/pgconn"
)

// faq: Таблица резервации. ReservationGood вычитает количество из таблицы goods и добавляем его в res_cen (reservation center)
func (d *db) ReservationGood(ctx context.Context, code string, stockId string, value int64) error {
	if code == "" || stockId == "" || value == 0 {
		return fmt.Errorf("result: code = %s, stock_id = %s, value = %d. must not be equal to code == '' or stock id == '' or value == 0", code, stockId, value)
	}

	var q string

	c, err := d.GetGoodsCountByStockId(ctx, stockId, code) // Находим количество товара доступное на складе
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
		d.logger.Fatal("couldn't open transaction")
		return err
	}

	chErr := make(chan error)

	// После исполнения запроса проверяем на ошибки
	go func(errs chan error) {
		for err := range errs { // faq: не нужно проверять на открытость канала, range делает это под капотом, поэтому без select
			if errors.Is(err, pgErr) {
				pgErr = errQ.(*pgconn.PgError)
				newErr := fmt.Errorf(
					fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
						pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
				)
				d.logger.Error(newErr)
			}
		}

	}(chErr)

	// Вычетаем количество резервируемого товара из таблицы goods
	q = `update goods set value = (value - $1) where code::text = $2 and stock_id = $3`
	_, errQ = t.Exec(ctx, q, value, code, stockId)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	chErr <- errQ

	// Создаем новую строку в таблице res_cen
	q = `insert into res_cen (good_code, stock_id, value) values ($1, $2, $3)`
	_, errQ = t.Exec(ctx, q, code, stockId, value)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	chErr <- errQ

	// Фиксируем транзакцию, если все окей
	if err := t.Commit(ctx); err != nil {
		d.logger.Errorf("failed to commit transaction: %s", err)
		return err
	}

	return nil
}

func (d *db) CancelGoodReservation(ctx context.Context, resId string) error {
	if resId == "" {
		return fmt.Errorf("resId is null")
	}

	t, err := d.client.Begin(ctx)
	if err != nil {
		d.logger.Fatal("couldn't open transaction")
		return err
	}

	var q string = `select good_code, value, stock_id::text from res_cen rc where rc.id::text = $1`

	var stock_id string
	var good_code, res_vl int64

	if errQ = t.QueryRow(ctx, q, resId).Scan(&good_code, &res_vl, &stock_id); errQ != nil { // Получаем данные с резервационного центра
		d.logger.Fatalf("func CancelGoodReservation query for search stock error %s", errQ)
		return errQ
	}

	chErr := make(chan error)

	go func(errs chan error) {
		for err := range errs {
			if errors.Is(err, pgErr) {
				pgErr = errQ.(*pgconn.PgError)
				newErr := fmt.Errorf(
					fmt.Sprintf("sql error: %s,  Detail: %s, Where: %s, Code: %s, SQLState: %s",
						pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()),
				)
				d.logger.Error(newErr)
			}
		}

	}(chErr)

	q = `delete from res_cen where id::text = $1` // Удаляем строку с резервационного центра
	_, errQ = t.Exec(ctx, q, resId)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	chErr <- errQ

	q = `update goods set value = (select value from goods where code = $2) + $1 where code = $2`
	_, errQ = t.Exec(ctx, q, res_vl, good_code)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	chErr <- errQ

	// Фиксируем транзакцию, если все окей
	if err := t.Commit(ctx); err != nil {
		d.logger.Errorf("failed to commit transaction: %s", err)
		return err
	}

	return nil
}
