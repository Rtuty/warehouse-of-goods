package db

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type dataSource struct {
	Host, Port, User, Passwd, Dbname, Sslmode string
}

var PstgCon dataSource

/*
Для подключения к postgresql, достаем переменные окружения из файла .env
Если все переменные существуют и иметют значение -> записываем в структуру подключения Pstgcon
*/
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

var dbConStr string = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	PstgCon.Host, PstgCon.Port, PstgCon.User, PstgCon.Passwd, PstgCon.Dbname, PstgCon.Sslmode)

// Создаем нового клиента postgresql
func NewDbClient() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
