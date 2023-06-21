package dbclient

import (
	"fmt"
	"log"
	"modules/utils"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

/*
Для подключения к postgresql, достаем переменные окружения из файла .env
Если все переменные существуют и иметют значение -> записываем в структуру подключения Pstgcon
*/
type dataSource struct {
	Host, Port, User, Passwd, Dbname, Sslmode string
}

var PstgCon dataSource

func GetConnection() {
	var res *dataSource = &PstgCon
	envVars := []string{"HOST", "PORT", "USER", "PASSWD", "DBNAME", "SSLMODE"}

	for _, v := range envVars {
		value := os.Getenv(v)
		if value == "" {
			panic(fmt.Sprintf("invalid environment variable %s", v))
		} else {
			field := reflect.ValueOf(res).Elem().FieldByNameFunc(
				func(fieldName string) bool {
					return strings.EqualFold(fieldName, v)
				})
			if field.IsValid() {
				field.SetString(value)
			}
		}
	}
	fmt.Println(PstgCon)
}

func NewClient(ctx context.Context, maxAttempts int, cn dataSource) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cn.User, cn.Passwd, cn.Host, cn.Port, cn.Dbname)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}
