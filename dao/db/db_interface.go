package db

import (
	"database/sql"
	"fmt"
	"reflect"
)

type DaoDB[T any] interface {
	Create(values T) (rowsAffected int64, err error)
	Read(columns []string, whereEquals *T, recordStart, recordEnd int) (res []T, err error)
	Update(whereEquals, newValues T) (rowsAffected int64, err error)
	Delete(whereEquals T) (rowsAffected int64, err error)
	Close() (err error)
}

// 1) all structure fields must be strings
// 2) the names of the structure fields must be the same as the names of the columns of the DB schema
func CheckStructForDAO(obj interface{}, db *sql.DB) (err error) {
	v := reflect.ValueOf(obj)
	err = CheckStructFormat(obj)

	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM tracks LIMIT 0")
	if err != nil {
		err = fmt.Errorf("ошибка чтения схемы таблицы: %v", err)
		return
	}
	defer rows.Close()

	var colNames []string
	colNames, err = rows.Columns()
	if err != nil {
		err = fmt.Errorf("ошибка чтения схемы таблицы: %v", err)
		return
	}

	for _, c := range colNames {
		field := v.FieldByName(c)
		if !field.IsValid() {
			err = fmt.Errorf("несоответствие схемы БД и указанной структуры: отсутствует поле %s", c)
			return
		}
	}

	return
}

// all structure fields must be strings
func CheckStructFormat(obj interface{}) (err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("ожидалась структура, а получено %s", v.Kind())
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Struct {
			err = fmt.Errorf("поле %d - структура, но вложенные структуры недопустимы", i)
			return
		}

		if field.Kind() != reflect.String {
			err = fmt.Errorf("поле %d должно быть string, а получено %s", i, field.Kind())
			return
		}
	}
	return
}
