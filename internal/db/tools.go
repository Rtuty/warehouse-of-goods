package db

import (
	"context"
	"errors"
	"fmt"
	"modules/pkg/dbclient"
)

// checkDbDublicate проверяет на дубликаты в БД, при создании нового товара или скалада
func checkDbDublicate(arg string, dbName string, ctx context.Context, c dbclient.Client) (bool, error) {
	var query string = "select exists (select 1 from "
	var exist bool

	if arg == "" {
		return false, errors.New("argument is empty")
	}

	switch dbName {
	case "goods":
		query = query + "goods where code::text"
	case "stocks":
		query = query + "stocks where name::text"
	default:
		return false, errors.New("database was not specified in the arguments of the function")
	}

	t, err := c.Begin(ctx) // Открываем транзакцию
	if err != nil {
		return false, fmt.Errorf("error when creating a transaction: %s", err)
	}

	query = query + "=$1::text)"

	if err := t.QueryRow(ctx, query, arg).Scan(&exist); err != nil {
		return false, err
	}

	if err := t.Commit(ctx); err != nil {
		return false, err
	}

	return exist, nil
}
