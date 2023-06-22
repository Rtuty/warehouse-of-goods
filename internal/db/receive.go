package db

import (
	"context"
	"fmt"
	"modules/internal/entities"
	"modules/utils"
)

// GetAllGoods получает все товары со всех складов
func (d *db) GetAllGoods(ctx context.Context) ([]entities.Good, error) {
	q := `select code, name, size, value from goods`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	rows, err := d.client.Query(ctx, q)
	if err != nil {
		d.logger.Fatalf("request execution error: %s", err)
		return nil, err
	}

	goods := make([]entities.Good, 0)

	for rows.Next() {
		var g entities.Good

		if err := rows.Scan(&g.Code, &g.Name, &g.Size, &g.Value); err != nil {
			d.logger.Fatalf("func getAllGoods scan rows error: %s", err)
			return nil, err
		}

		goods = append(goods, g)
	}

	if err := rows.Err(); err != nil {
		d.logger.Fatalf("error checking failed: %s", err)
		return nil, err
	}

	return goods, nil
}

// GetGoodByCode получает товар по его коду
func (d *db) GetGoodByCode(ctx context.Context, code string) (entities.Good, error) {
	q := `select code, name, size, value from goods where code::text = $1::text`
	d.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	var g entities.Good

	if err := d.client.QueryRow(ctx, q, code).Scan(&g.Code, &g.Name, &g.Size, &g.Value); err != nil {
		d.logger.Fatalf("request execution error: %s query: %s", err, q)
		return entities.Good{}, nil
	}
	return g, nil
}

// GetGoodsCountByStockId возвращает количество всех товаров на складе. Если указан code товара => выводим общее количество конкретного товара
func (d *db) GetGoodsCountByStockId(ctx context.Context, stockId string, code string) (int64, error) {
	var count int64
	q := `select sum(value) from goods g
			inner join stocks s on s.id = g.stock_id
		  where g.stock_id = $1 and s.available`

	if code != "" {
		q = q + ` and code::text = $2`

		if err := d.client.QueryRow(ctx, q, stockId, code).Scan(&count); err != nil {
			d.logger.Fatalf("request execution error: %s query: %s", err, q)
			return -1, nil
		}

		return count, nil
	}

	if err := d.client.QueryRow(ctx, q, stockId).Scan(&count); err != nil {
		d.logger.Fatalf("request execution error: %s query: %s", err, q)
		return -1, nil
	}

	return count, nil
}
