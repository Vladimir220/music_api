package caching

import "time"

type DaoCaching interface {
	Get(key string) (res string, err error)
	Set(key, value string, ttl time.Duration) (err error)
	HSet(key string, data interface{}) (err error)
	HGet(key string) (res map[string]string, err error)
}
