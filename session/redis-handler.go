package session

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisHandler struct {
	pool     *redis.Pool // redis connection pool
	prefix   string
	lifetime time.Duration
}

// NewRedisHandler Create a redis session handler
func NewRedisHandler(pool *redis.Pool, prefix string, lifetime time.Duration) *RedisHandler {
	return &RedisHandler{pool, prefix, lifetime}
}

func (rh *RedisHandler) Read(id string) string {
	c := rh.pool.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", rh.prefix+":"+id))
	if err != nil {
		return ""
	}

	var value string

	json.Unmarshal(b, value)

	return value
}

func (rh *RedisHandler) Write(id string, data string) {
	c := rh.pool.Get()
	defer c.Close()

	c.Do("SETEX", rh.prefix+":"+id, int64(rh.lifetime), data)
}
