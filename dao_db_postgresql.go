package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type daoPostgreSQL[T any] struct {
	db        *sql.DB
	tableName string
}

func (dao daoPostgreSQL[T]) Create(values T) (rowsAffected int64, err error) {
	cols := make([]string, 0, 5)
	vals := make([]string, 0, 5)
	v := reflect.ValueOf(values)
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i).String()
		vals = append(vals, fieldVal)

		fieldName := v.Type().Field(i).Name
		fieldName = fmt.Sprintf(`"%s"`, fieldName)
		cols = append(cols, fieldName)
	}

	valQuery := make([]string, 0)
	for i := 1; i <= len(vals); i++ {
		valQuery = append(valQuery, fmt.Sprintf("$%d", i))
	}

	colsStr := strings.Join(cols, ", ")
	tableStr := dao.tableName
	valQueryStr := strings.Join(valQuery, ", ")
	queryStr := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s);`, tableStr, colsStr, valQueryStr)

	allQueryParams := make([]any, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		allQueryParams = append(allQueryParams, vals[i])
	}

	var res sql.Result
	res, err = dao.db.Exec(queryStr, allQueryParams...)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения команды insert: %v", err)
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("ошибка получения количества изменений: %v", err)
		return
	}

	return
}

func (dao daoPostgreSQL[T]) Read(columns []string, whereEquals *T, recordStart, recordEnd int) (res []T, err error) {
	whereQuery := make([]string, 0, 5)
	vals := make([]string, 0, 5)
	if whereEquals != nil {
		v := reflect.ValueOf(whereEquals).Elem()
		j := 1
		for i := 0; i < v.NumField(); i++ {
			fieldVal := v.Field(i).String()
			if fieldVal == "" {
				continue
			}
			vals = append(vals, fieldVal)

			fieldName := v.Type().Field(i).Name
			fieldName = fmt.Sprintf(`"%s"`, fieldName)
			whereQuery = append(whereQuery, fmt.Sprintf(`%s = $%d`, fieldName, j))
			j++
		}
	}

	columnsStr := ""
	if len(columns) == 0 {
		columnsStr = "*"
	} else {
		for i := range columns {
			if columns[i] == "*" {
				columnsStr = "*"
				break
			}
			columns[i] = fmt.Sprintf(`"%s"`, columns[i])
		}
	}
	if columnsStr != "*" {
		columnsStr = strings.Join(columns, ", ")
	}

	whereStr := ""
	if whereEquals != nil && len(whereQuery) != 0 {
		whereStr = strings.Join(whereQuery, " AND ")
		whereStr = fmt.Sprintf("WHERE %s", whereStr)
	}

	tableStr := dao.tableName

	queryStr := fmt.Sprintf(`SELECT %s FROM %s %s OFFSET %d LIMIT %d;`, columnsStr, tableStr, whereStr, recordStart, recordEnd)

	allQueryParams := make([]any, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		allQueryParams = append(allQueryParams, vals[i])
	}

	var rows *sql.Rows
	rows, err = dao.db.Query(queryStr, allQueryParams...)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения команды select: %v", err)
		return
	}
	defer rows.Close()

	var colNames []string
	colNames, err = rows.Columns()
	if err != nil {
		err = fmt.Errorf("ошибка чтения схемы таблицы: %v", err)
		return
	}

	for rows.Next() {
		buf := make([]any, len(colNames))
		for i := range buf {
			buf[i] = new(string)
		}
		err = rows.Scan(buf...)
		if err != nil {
			err = fmt.Errorf("ошибка получения результатов: %v", err)
			return
		}

		var tt T
		var t *T = &tt
		v := reflect.ValueOf(t).Elem()
		for i, c := range colNames {
			field := v.FieldByName(c)
			if !field.IsValid() {
				err = fmt.Errorf("несоответствие схемы БД и указанной структуры: отсутствует поле %s", c)
				return
			}

			str := buf[i].(*string)
			field.SetString(*str)

		}
		res = append(res, *t)

	}
	return
}

func (dao daoPostgreSQL[T]) Update(whereEquals T, newValues T) (rowsAffected int64, err error) {
	j := 1
	set := make([]string, 0, 5)
	setVals := make([]string, 0, 5)
	v := reflect.ValueOf(newValues)
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i).String()
		if fieldVal == "" {
			continue
		}
		setVals = append(setVals, fieldVal)

		fieldName := v.Type().Field(i).Name
		fieldName = fmt.Sprintf(`"%s"`, fieldName)
		set = append(set, fmt.Sprintf(`%s = $%d`, fieldName, j))
		j++
	}
	if len(set) == 0 {
		err = fmt.Errorf("отсутствуют новые значения для update")
		return
	}

	where := make([]string, 0, 5)
	whereVals := make([]string, 0, 5)
	v = reflect.ValueOf(whereEquals)
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i).String()
		if fieldVal == "" {
			continue
		}
		whereVals = append(whereVals, fieldVal)

		fieldName := v.Type().Field(i).Name
		fieldName = fmt.Sprintf(`"%s"`, fieldName)
		where = append(where, fmt.Sprintf(`%s = $%d`, fieldName, j))
		j++
	}
	if len(where) == 0 {
		err = fmt.Errorf("отсутствует условия для update")
		return
	}

	tableStr := dao.tableName
	setStr := strings.Join(set, ", ")
	whereStr := strings.Join(where, " AND ")
	queryStr := fmt.Sprintf(`UPDATE %s SET %s WHERE %s;`, tableStr, setStr, whereStr)

	paramLen := len(whereVals) + len(setVals)
	allQueryParams := make([]any, 0, paramLen)
	for i := 0; i < len(setVals); i++ {
		allQueryParams = append(allQueryParams, setVals[i])
	}
	for i := 0; i < len(whereVals); i++ {
		allQueryParams = append(allQueryParams, whereVals[i])
	}

	var res sql.Result
	res, err = dao.db.Exec(queryStr, allQueryParams...)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения update: %v", err)
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("ошибка получения количества изменений: %v", err)
		return
	}
	return
}

func (dao daoPostgreSQL[T]) Delete(whereEquals T) (rowsAffected int64, err error) {
	j := 1
	where := make([]string, 0, 5)
	vals := make([]string, 0, 5)
	v := reflect.ValueOf(whereEquals)
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i).String()
		if fieldVal == "" {
			continue
		}
		vals = append(vals, fieldVal)

		fieldName := v.Type().Field(i).Name
		fieldName = fmt.Sprintf(`"%s"`, fieldName)
		where = append(where, fmt.Sprintf(`%s = $%d`, fieldName, j))
		j++
	}
	if len(where) == 0 {
		err = fmt.Errorf("отсутствует условия для delete")
		return
	}

	tableStr := dao.tableName
	whereStr := strings.Join(where, " AND ")
	queryStr := fmt.Sprintf(`DELETE FROM %s WHERE %s;`, tableStr, whereStr)

	allQueryParams := make([]any, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		allQueryParams = append(allQueryParams, vals[i])
	}

	var res sql.Result
	res, err = dao.db.Exec(queryStr, allQueryParams...)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения команды update: %v", err)
		return
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("ошибка получения количества изменений: %v", err)
		return
	}
	return
}

func (dao daoPostgreSQL[T]) CheckMigrations() (err error) {
	var driver database.Driver
	driver, err = postgres.WithInstance(dao.db, &postgres.Config{})
	if err != nil {
		err = fmt.Errorf("ошибка создания драйвера: %v", err)
		return
	}

	dbName := os.Getenv("DB_NAME")
	var m *migrate.Migrate
	m, err = migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		err = fmt.Errorf("ошибка создания мигратора: %v", err)
		return
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		err = fmt.Errorf("ошибка применения миграций: %v", err)
		return
	} else {
		err = nil
	}

	return
}

func (dao *daoPostgreSQL[T]) Init() (err error) {
	dao.tableName = os.Getenv("DB_TABLE_NAME")
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
	)

	loginInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	dao.db, err = sql.Open("postgres", loginInfo)
	if err != nil {
		err = fmt.Errorf("ошибка подключения к БД: %v", err)
		return
	}

	return
}

func (dao *daoPostgreSQL[T]) Close() (err error) {
	err = dao.db.Close()
	if err != nil {
		err = fmt.Errorf("ошибка закрытия БД: %v", err)
	}
	return
}

func createDaoPostgreSQL[T any]() (dao DaoDB[T], err error) {
	daoPSQL := &daoPostgreSQL[T]{}
	err = daoPSQL.Init()
	if err != nil {
		return
	}

	err = daoPSQL.CheckMigrations()
	if err != nil {
		return
	}

	var test T
	err = checkStructForDAO(test, daoPSQL.db)
	if err != nil {
		return
	}

	return daoPSQL, err
}
