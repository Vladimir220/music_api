package db

import (
	"database/sql"
	"fmt"
	"music_api/models"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DaoPostgresMusicApi struct {
	db *sql.DB
}

func (dao DaoPostgresMusicApi) Create(values models.Track) (rowsAffected int64, err error) {
	if values.Song == "" || values.Group_name == "" {
		return int64(0), fmt.Errorf("создания песни с пустым названием песни или группы")
	}

	var id string
	queryStr := `WITH ins AS (
				INSERT INTO Groups ("Group_name")
				VALUES ($1)
				ON CONFLICT ("Group_name") 
				DO NOTHING
				RETURNING "Id"
			)
			SELECT COALESCE((SELECT "Id" FROM ins), (SELECT "Id" FROM Groups WHERE "Group_name" = $1));`
	rows, err := dao.db.Query(queryStr, values.Group_name)
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		err = fmt.Errorf("ошибка получения результатов: %v", err)
		return rowsAffected, err
	}
	rows.Close()

	queryStr = fmt.Sprintf(`INSERT INTO Tracks ("Group_id","Song","Release_date","Song_lyrics","Link") VALUES ((%s),$1,$2,$3,$4);`, id)

	res, err := dao.db.Exec(queryStr, values.Song, values.Release_date, values.Song_lyrics, values.Link)
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

func (dao DaoPostgresMusicApi) Read(columns []string, whereEquals *models.Track, recordStart, recordEnd int) (res []models.Track, err error) {
	var id string
	where := []string{}
	vals := []any{}
	i := 1

	if whereEquals != nil {
		if whereEquals.Group_name != "" {
			queryStr := `SELECT "Id" FROM Groups WHERE "Group_name" = $1`
			rows, err := dao.db.Query(queryStr, whereEquals.Group_name)
			if err != nil {
				err = fmt.Errorf("ошибка выполнения команды select: %v", err)
				return res, err
			}

			rows.Next()
			err = rows.Scan(&id)
			if err != nil {
				err = fmt.Errorf("ошибка получения результатов: %v", err)
				return res, err
			}
			rows.Close()
		}

		if whereEquals.Group_name != "" {
			where = append(where, fmt.Sprintf(`"Group_id" = $%d`, i))
			vals = append(vals, id)
			i++
		}
		if whereEquals.Song != "" {
			where = append(where, fmt.Sprintf(`"Song" = $%d`, i))
			vals = append(vals, whereEquals.Song)
			i++
		}
		if whereEquals.Release_date != "" {
			where = append(where, fmt.Sprintf(`"Release_date" = $%d`, i))
			vals = append(vals, whereEquals.Release_date)
			i++
		}
		if whereEquals.Song_lyrics != "" {
			where = append(where, fmt.Sprintf(`"Song_lyrics" = $%d`, i))
			vals = append(vals, whereEquals.Song_lyrics)
			i++
		}
		if whereEquals.Link != "" {
			where = append(where, fmt.Sprintf(`"Link" = $%d`, i))
			vals = append(vals, whereEquals.Link)
			i++
		}
	}
	whereStr := ""
	if whereEquals != nil && len(where) != 0 {
		whereStr = strings.Join(where, " AND ")
		whereStr = fmt.Sprintf("WHERE %s", whereStr)
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
	} else {
		columnsStr = `"Group_name", "Song", "Release_date", "Song_lyrics", "Link"`
	}

	queryStr := fmt.Sprintf(`SELECT %s FROM Tracks INNER JOIN "groups" ON Tracks."Group_id" = Groups."Id" %s OFFSET $%d LIMIT $%d;`, columnsStr, whereStr, i, i+1)
	vals = append(vals, recordStart)
	vals = append(vals, recordEnd)

	rows, err := dao.db.Query(queryStr, vals...)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения команды select: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		track := models.Track{}
		err = rows.Scan(&track.Group_name, &track.Song, &track.Release_date, &track.Song_lyrics, &track.Link)
		if err != nil {
			err = fmt.Errorf("ошибка получения результатов: %v", err)
			return
		}

		res = append(res, track)
	}
	return
}

func (dao DaoPostgresMusicApi) Update(whereEquals models.Track, newValues models.Track) (rowsAffected int64, err error) {
	set := []string{}

	// Получаю id группы
	var id string
	queryStr := `SELECT "Id" FROM Groups WHERE "Group_name" = $1;`
	rows, err := dao.db.Query(queryStr, whereEquals.Group_name)
	if err != nil {
		err = fmt.Errorf("ошибка выполнения команды select: %v", err)
		return rowsAffected, err
	}
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		err = fmt.Errorf("ошибка получения результатов: %v", err)
		return rowsAffected, err
	}
	rows.Close()

	// Если в новых значениях есть новое название группы, то либо ищем id новой группы, либо создаём группу
	if newValues.Group_name != "" {
		var newId string
		queryStr := `WITH ins AS (
					INSERT INTO Groups ("Group_name")
					VALUES ($1)
					ON CONFLICT ("Group_name") 
					DO NOTHING
					RETURNING "Id"
				)
				SELECT COALESCE((SELECT "Id" FROM ins), (SELECT "Id" FROM Groups WHERE "Group_name" = $1));`
		rows, err = dao.db.Query(queryStr, newValues.Group_name)
		if err != nil {
			return rowsAffected, err
		}
		rows.Next()
		err = rows.Scan(&newId)
		if err != nil {
			err = fmt.Errorf("ошибка получения результатов: %v", err)
			return rowsAffected, err
		}
		rows.Close()
		set = append(set, fmt.Sprintf(`"Group_id" = %s`, newId))
	}

	vals := []any{}
	i := 1
	if newValues.Song != "" {
		set = append(set, fmt.Sprintf(`"Song" = $%d`, i))
		vals = append(vals, newValues.Song)
		i++
	}
	if newValues.Release_date != "" {
		set = append(set, fmt.Sprintf(`"Release_date" = $%d`, i))
		vals = append(vals, newValues.Release_date)
		i++
	}
	if newValues.Song_lyrics != "" {
		set = append(set, fmt.Sprintf(`"Song_lyrics" = $%d`, i))
		vals = append(vals, newValues.Song_lyrics)
		i++
	}
	if newValues.Link != "" {
		set = append(set, fmt.Sprintf(`"Link" = $%d`, i))
		vals = append(vals, newValues.Link)
		i++
	}
	setStr := strings.Join(set, ", ")

	where := []string{}
	if whereEquals.Group_name != "" {
		where = append(where, fmt.Sprintf(`"Group_id" = %s`, id))
	}
	if whereEquals.Song != "" {
		where = append(where, fmt.Sprintf(`"Song" = $%d`, i))
		vals = append(vals, whereEquals.Song)
		i++
	}
	if whereEquals.Release_date != "" {
		where = append(where, fmt.Sprintf(`"Release_date" = $%d`, i))
		vals = append(vals, whereEquals.Release_date)
		i++
	}
	if whereEquals.Song_lyrics != "" {
		where = append(where, fmt.Sprintf(`"Song_lyrics" = $%d`, i))
		vals = append(vals, whereEquals.Song_lyrics)
		i++
	}
	if whereEquals.Link != "" {
		where = append(where, fmt.Sprintf(`"Link" = $%d`, i))
		vals = append(vals, whereEquals.Link)
		i++
	}
	whereStr := strings.Join(where, " AND ")

	queryStr = fmt.Sprintf(`UPDATE Tracks SET %s WHERE %s;`, setStr, whereStr)

	var res sql.Result
	res, err = dao.db.Exec(queryStr, vals...)
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

func (dao DaoPostgresMusicApi) Delete(whereEquals models.Track) (rowsAffected int64, err error) {
	var id string
	queryStr := `SELECT "Id" FROM Groups WHERE "Group_name" = $1;`
	rows, err := dao.db.Query(queryStr, whereEquals.Group_name)
	if err != nil {
		return rowsAffected, err
	}
	rows.Next()
	err = rows.Scan(&id)
	if err != nil {
		err = fmt.Errorf("ошибка получения результатов: %v", err)
		return rowsAffected, err
	}
	rows.Close()

	where := []string{}
	i := 1
	vals := []any{}
	if whereEquals.Group_name != "" {
		where = append(where, fmt.Sprintf(`"Group_id" = %s`, id))
	}
	if whereEquals.Song != "" {
		where = append(where, fmt.Sprintf(`"Song" = $%d`, i))
		vals = append(vals, whereEquals.Song)
		i++
	}
	if whereEquals.Release_date != "" {
		where = append(where, fmt.Sprintf(`"Release_date" = $%d`, i))
		vals = append(vals, whereEquals.Release_date)
		i++
	}
	if whereEquals.Song_lyrics != "" {
		where = append(where, fmt.Sprintf(`"Song_lyrics" = $%d`, i))
		vals = append(vals, whereEquals.Song_lyrics)
		i++
	}
	if whereEquals.Link != "" {
		where = append(where, fmt.Sprintf(`"Link" = $%d`, i))
		vals = append(vals, whereEquals.Link)
		i++
	}
	whereStr := strings.Join(where, " AND ")
	queryStr = fmt.Sprintf(`DELETE FROM Tracks WHERE %s;`, whereStr)

	res, err := dao.db.Exec(queryStr, vals...)
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

func (dao DaoPostgresMusicApi) CheckMigrations() (err error) {
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

func (dao *DaoPostgresMusicApi) Init() (err error) {
	var (
		user      = os.Getenv("DB_USER")
		password  = os.Getenv("DB_PASSWORD")
		dbName    = os.Getenv("DB_NAME")
		host      = os.Getenv("DB_HOST")
		port      = os.Getenv("DB_PORT")
		dbDNSName = os.Getenv("DB_NAME")
		dockerMod = os.Getenv("DOCKER_MOD")
		loginInfo string
	)

	if dockerMod != "1" {
		loginInfo = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, "postgres")
	} else {
		loginInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, dbDNSName, "postgres")
	}

	db, err := sql.Open("postgres", loginInfo)
	if err != nil {
		err = fmt.Errorf("ошибка подключения к БД: %v", err)
		return
	}

	query := "SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"
	var isDbExist bool
	err = db.QueryRow(query, dbName).Scan(&isDbExist)
	if err != nil {
		err = fmt.Errorf("ошибка при проверке существования базы данных: %v", err)
		return
	}

	if !isDbExist {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			err = fmt.Errorf("ошибка при создании базы данных: %v", err)
			return
		}
	}

	if dockerMod != "1" {
		loginInfo = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	} else {
		loginInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, dbDNSName, dbName)
	}

	dao.db, err = sql.Open("postgres", loginInfo)
	if err != nil {
		err = fmt.Errorf("ошибка подключения к БД: %v", err)
		return
	}

	return
}

func (dao *DaoPostgresMusicApi) Close() (err error) {
	err = dao.db.Close()
	if err != nil {
		err = fmt.Errorf("ошибка закрытия БД: %v", err)
	}
	return
}

func CreateDaoPostgresMusicApi() (dao DaoDB[models.Track], err error) {
	psql := &DaoPostgresMusicApi{}
	err = psql.Init()
	if err != nil {
		return
	}

	err = psql.CheckMigrations()
	if err != nil {
		return
	}

	return psql, err
}
