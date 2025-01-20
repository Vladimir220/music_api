package caching

import (
	"context"
	"fmt"
	"music_api/dao/db"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type DaoRedis struct {
	rdb *redis.Client
}

func (r *DaoRedis) Init() (err error) {
	var (
		host     = os.Getenv("REDIS_HOST")
		port     = os.Getenv("REDIS_PORT")
		password = os.Getenv("REDIS_PASSWORD")
	)
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return
	}

	r.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	_, err = r.rdb.Ping(context.Background()).Result()
	if err != nil {
		err = fmt.Errorf("не удалось подключиться к Redis: %v", err)
		return
	}

	return
}

func (r *DaoRedis) Get(key string) (res string, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err = r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = fmt.Errorf("ключ не найден")
			return
		} else {
			err = fmt.Errorf("ошибка получения значения: %v", err)
			return
		}
	}
	return
}

func (r *DaoRedis) Set(key, value string, ttl time.Duration) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = r.rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		err = fmt.Errorf("ошибка установки значения: %v", err)
		return
	}
	return
}

func (r *DaoRedis) HSet(key string, data interface{}) (err error) {
	err = db.CheckStructFormat(data)
	if err != nil {
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = r.rdb.HSet(ctx, key, data).Err()
	if err != nil {
		err = fmt.Errorf("ошибка установки значения: %v", err)
		return
	}
	return
}

func (r *DaoRedis) HGet(key string) (res map[string]string, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err = r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = fmt.Errorf("ключ не найден")
			return
		} else {
			err = fmt.Errorf("ошибка получения значения: %v", err)
			return
		}
	}
	return
}

func CreateDaoRedis() (dao DaoCaching, err error) {
	redis := &DaoRedis{}
	err = redis.Init()
	dao = redis
	return
}
