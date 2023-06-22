package db

import (
	"context"
	"fmt"
	"modules/utils"
)

// faq: Таблица резервации. Вычетаем количество из таблицы goods и добавляем его в res_cen (reservation center)
func (d *db) ReserveGood(ctx context.Context, code string, stockId string, value int) error {
	if code == "" || stockId == "" || value == 0 {
		return fmt.Errorf("result: code = %s, stock_id = %s, value = %d. must not be equal to code == '' or stock id == '' or value == 0", code, stockId, value)
	}
	// Если value для резервации > количество товара на складе => ошибка
	var q string

	d.GetGoodsCountByStockId(ctx, stockId, code)
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	return nil
}

func (d *db) CancelGoodReserve(ctx context.Context, code string, stockId string, value int) error {
	return nil
}
